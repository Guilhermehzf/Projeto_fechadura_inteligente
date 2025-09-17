package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Estrutura para a resposta do endpoint /status
type StatusResposta struct {
	TrancaAberta    bool   `json:"tranca_aberta"`
	UltimaAtualizacao string `json:"ultima_atualizacao"`
}

// --- CONSTANTES DO PROJETO ---
// IMPORTANTE: Use as mesmas credenciais do seu ESP32
const MQTT_BROKER = "mqtts://0dfb02c89e7f487f9a2f8e5e29729297.s1.eu.hivemq.cloud:8883"
const MQTT_CLIENT_ID = "api-go-client-final"
const MQTT_USER = "esp32-device" // <<< CORRIGIDO: Usando o usuário que você já tem
const MQTT_PASS = "Gatitcha1!"   // <<< CORRIGIDO: Use a senha correta para este usuário

const MQTT_TOPIC_COMMANDS = "fechadura/comandos"
const MQTT_TOPIC_STATE = "fechadura/estado"

// --- Estado Global e Sincronização ---
var stateMutex sync.Mutex
var trancaEstaAberta bool
var mqttClient MQTT.Client
var initialMessageReceived = make(chan bool, 1) // Canal para sinalizar o recebimento do estado inicial

// messageHandler é chamado quando a API recebe uma mensagem do ESP32 sobre seu estado
var messageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Printf(">> MENSAGEM DE ESTADO RECEBIDA | Tópico: %s | Mensagem: %s\n", msg.Topic(), msg.Payload())
	
	// 1. Cria uma estrutura para receber os dados do JSON
	type StatePayload struct {
		TrancaAberta bool `json:"tranca_aberta"`
	}

	var payload StatePayload

	// 2. Tenta decodificar o payload da mensagem na nossa estrutura
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Printf("   Erro ao decodificar JSON de estado: %v. A mensagem era: %s", err, msg.Payload())
		return
	}

	// 3. Protege o acesso e atualiza a variável global
	stateMutex.Lock()
	defer stateMutex.Unlock()

	trancaEstaAberta = payload.TrancaAberta
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
	if token := client.Subscribe(MQTT_TOPIC_STATE, 1, messageHandler); token.Wait() && token.Error() != nil {
		log.Printf("Erro ao assinar o tópico: %s\n", token.Error())
		return
	}
	log.Printf("Assinatura ao tópico '%s' realizada com sucesso!\n", MQTT_TOPIC_STATE)
}

var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	log.Printf("Conexão com o Broker perdida: %v\n", err)
}

// --- FUNÇÃO MAIN ---
func main() {
	opts := MQTT.NewClientOptions().AddBroker(MQTT_BROKER).SetClientID(MQTT_CLIENT_ID)
	opts.SetUsername(MQTT_USER).SetPassword(MQTT_PASS)
	opts.SetOnConnectHandler(connectHandler)
	opts.SetConnectionLostHandler(connectLostHandler)
	
	mqttClient = MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// <<< LÓGICA DE SINCRONIZAÇÃO CORRIGIDA ---
	// Espera pela primeira mensagem de estado chegar ou um timeout de 5 segundos
	log.Println("Aguardando estado inicial do broker...")
	select {
	case <-initialMessageReceived:
		log.Println("Estado inicial recebido e sincronizado com sucesso!")
	case <-time.After(5 * time.Second):
		// Se não houver mensagem retida, definimos um padrão seguro (trancado)
		log.Println("Timeout: Nenhuma mensagem de estado retida encontrada. Definindo estado padrão como 'fechada'.")
		stateMutex.Lock()
		trancaEstaAberta = false
		stateMutex.Unlock()
	}

	// Inicia o servidor HTTP somente APÓS a sincronização
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/toggle", toggleHandler)

	addr := "0.0.0.0:8080"
	log.Printf("API HTTP rodando e escutando em http://%s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor HTTP: %s\n", err)
	}
}

// --- HANDLERS HTTP ---
// O toggleHandler e o statusHandler continuam exatamente iguais
func toggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido, use POST", http.StatusMethodNotAllowed)
		return
	}
	
	payload := `{"command":"toggle"}`
	token := mqttClient.Publish(MQTT_TOPIC_COMMANDS, 1, false, payload)
	token.Wait()

	if token.Error() != nil {
		log.Printf("Erro ao publicar no MQTT: %s\n", token.Error())
		http.Error(w, "Erro interno ao enviar comando", http.StatusInternalServerError)
		return
	}

	log.Printf("<< COMANDO PUBLICADO | Tópico: %s | Mensagem: %s\n", MQTT_TOPIC_COMMANDS, payload)
	w.Write([]byte(`{"status": "comando 'toggle' enviado com sucesso"}`))
}

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