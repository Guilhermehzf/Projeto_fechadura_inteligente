package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

// Estrutura para a resposta do endpoint /status
type StatusResposta struct {
	TrancaAberta    bool   `json:"tranca_aberta"`
	UltimaAtualizacao string `json:"ultima_atualizacao"`
}

// --- Variáveis Globais de Configuração ---
// Estas variáveis serão preenchidas no início da função main a partir do arquivo .env
var (
	mqttBroker      string
	mqttClientID    string
	mqttUser        string
	mqttPass        string
	mqttTopicCommands string
	mqttTopicState    string
)

// --- Variáveis Globais de Estado ---
var stateMutex sync.Mutex
var trancaEstaAberta bool
var mqttClient MQTT.Client
var initialMessageReceived = make(chan bool, 1)

// messageHandler é chamado quando a API recebe uma mensagem do ESP32 sobre seu estado
var messageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Printf(">> MENSAGEM DE ESTADO RECEBIDA | Tópico: %s | Mensagem: %s\n", msg.Topic(), msg.Payload())

	// Decodifica o payload da mensagem JSON
	type StatePayload struct {
		TrancaAberta bool `json:"tranca_aberta"`
	}
	var payload StatePayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Printf("   Erro ao decodificar JSON de estado: %v", err)
		return
	}

	// Atualiza a variável global de estado
	stateMutex.Lock()
	trancaEstaAberta = payload.TrancaAberta
	stateMutex.Unlock()
	log.Printf("   Estado interno da API foi atualizado para: %v\n", trancaEstaAberta)

	// Sinaliza que a primeira mensagem (o estado retido) foi recebida
	select {
	case initialMessageReceived <- true:
	default:
	}
}

// connectHandler é chamado na conexão/reconexão com o broker
var connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
	log.Println("Conectado ao Broker MQTT!")
	// Ao conectar, a API se inscreve para ouvir o estado do ESP32.
	if token := client.Subscribe(mqttTopicState, 1, messageHandler); token.Wait() && token.Error() != nil {
		log.Printf("Erro ao assinar o tópico: %s\n", token.Error())
		return
	}
	log.Printf("Assinatura ao tópico '%s' realizada com sucesso!\n", mqttTopicState)
}

var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	log.Printf("Conexão com o Broker perdida: %v\n", err)
}

// --- FUNÇÃO MAIN ---
func main() {
	// --- CARREGA A CONFIGURAÇÃO UMA ÚNICA VEZ ---
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Não foi possível carregar o arquivo .env")
	}

	// Preenche as variáveis globais de configuração
	mqttBroker = os.Getenv("MQTT_BROKER")
	mqttClientID = os.Getenv("MQTT_CLIENT_ID")
	mqttUser = os.Getenv("MQTT_USER")
	mqttPass = os.Getenv("MQTT_PASS")
	mqttTopicCommands = os.Getenv("MQTT_TOPIC_COMMANDS")
	mqttTopicState = os.Getenv("MQTT_TOPIC_STATE")
	
	if mqttBroker == "" || mqttUser == "" || mqttPass == "" || mqttClientID == "" {
		log.Fatal("ERRO: As variáveis de ambiente (MQTT_BROKER, MQTT_CLIENT_ID, MQTT_USER, MQTT_PASS) devem ser definidas no arquivo .env")
	}

	// --- Conecta ao MQTT com as variáveis carregadas ---
	opts := MQTT.NewClientOptions().AddBroker(mqttBroker).SetClientID(mqttClientID)
	opts.SetUsername(mqttUser).SetPassword(mqttPass)
	opts.SetOnConnectHandler(connectHandler)
	opts.SetConnectionLostHandler(connectLostHandler)
	
	mqttClient = MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Sincroniza o estado inicial
	log.Println("Aguardando estado inicial do broker...")
	select {
	case <-initialMessageReceived:
		log.Println("Estado inicial recebido e sincronizado com sucesso!")
	case <-time.After(5 * time.Second):
		log.Println("Timeout: Nenhuma mensagem de estado retida encontrada. Definindo estado padrão como 'fechada'.")
		stateMutex.Lock()
		trancaEstaAberta = false
		stateMutex.Unlock()
	}

	// Inicia o servidor HTTP
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/toggle", toggleHandler)
	addr := "0.0.0.0:8080"
	log.Printf("API HTTP rodando e escutando em http://%s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor HTTP: %s\n", err)
	}
}

// --- HANDLERS HTTP ---

// toggleHandler agora está limpo e usa a variável global
func toggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido, use POST", http.StatusMethodNotAllowed)
		return
	}
	
	payload := `{"command":"toggle"}`
	token := mqttClient.Publish(mqttTopicCommands, 1, false, payload)
	token.Wait()

	if token.Error() != nil {
		log.Printf("Erro ao publicar no MQTT: %s\n", token.Error())
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	log.Printf("<< COMANDO PUBLICADO | Tópico: %s | Mensagem: %s\n", mqttTopicCommands, payload)
	w.Write([]byte(`{"status": "comando 'toggle' enviado com sucesso"}`))
}

// statusHandler continua o mesmo
func statusHandler(w http.ResponseWriter, r *http.Request) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	resposta := StatusResposta{
		TrancaAberta:    trancaEstaAberta,
		UltimaAtualizacao: time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resposta)
}