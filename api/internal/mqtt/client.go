package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"

	"smartlock/internal/config"
	"smartlock/internal/db"
	"smartlock/internal/state"
)

type Client struct {
	cfg   *config.Config
	store *state.Store

	client  paho.Client
	Initial chan struct{}
}

func New(cfg *config.Config, st *state.Store) *Client {
	return &Client{cfg: cfg, store: st, Initial: make(chan struct{}, 1)}
}

func (c *Client) Start() {
	host, _ := os.Hostname()
	clientID := c.cfg.MQTTClientID + "-" + host + "-" + strconv.Itoa(os.Getpid())

	opts := paho.NewClientOptions().
		AddBroker(c.cfg.MQTTBroker).
		SetClientID(clientID).
		SetUsername(c.cfg.MQTTUser).
		SetPassword(c.cfg.MQTTPass).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetKeepAlive(60 * time.Second).
		SetPingTimeout(20 * time.Second).
		SetConnectTimeout(15 * time.Second).
		SetWriteTimeout(15 * time.Second).
		SetMaxReconnectInterval(30 * time.Second).
		SetOrderMatters(false)

	opts.SetOnConnectHandler(func(cl paho.Client) {
		log.Println("Conectado ao broker MQTT")
		if tok := cl.Subscribe(c.cfg.MQTTTopicState, 1, c.onMessage); tok.Wait() && tok.Error() != nil {
			log.Printf("Erro ao assinar tópico %s: %v", c.cfg.MQTTTopicState, tok.Error())
			return
		}
		log.Printf("Assinado em %q", c.cfg.MQTTTopicState)
		c.store.SetConnected(true)
	})

	opts.SetConnectionLostHandler(func(cl paho.Client, err error) {
		log.Printf("Conexão MQTT perdida: %v", err)
		c.store.SetConnected(false)
	})

	c.client = paho.NewClient(opts)
	if tok := c.client.Connect(); tok.Wait() && tok.Error() != nil {
		log.Fatalf("Erro conectando ao broker: %v", tok.Error())
	}

	// Aguarda retained inicial por até 5s
	log.Println("Aguardando estado inicial (retained) por até 5s...")
	select {
	case <-c.Initial:
		log.Println("Estado inicial recebido.")
	case <-time.After(5 * time.Second):
		log.Println("Sem retained inicial; seguindo em frente.")
	}
}

func (c *Client) onMessage(_ paho.Client, msg paho.Message) {
	log.Printf(">> MQTT mensagem | tópico=%s qos=%d retained=%v | payload=%s\n",
		msg.Topic(), msg.Qos(), msg.Retained(), string(msg.Payload()))

	// Atualiza estado interno
	c.store.UpdateFromMQTT(msg.Topic(), msg.Qos(), msg.Retained(), msg.Payload())

	// Se for o tópico de estado, registra no histórico
	if msg.Topic() == c.cfg.MQTTTopicState {
		var payload map[string]any
		if err := json.Unmarshal(msg.Payload(), &payload); err == nil {
			action := "locked"
			method := "auto"

			if v, ok := payload["tranca_aberta"].(bool); ok && v {
				action = "unlocked"
			}
			if v, ok := payload["method"].(string); ok && v != "" {
				method = v
			}

			// Insere no histórico
			if err := db.InsertHistory(context.Background(), c.cfg.DB, action, method); err != nil {
				fmt.Printf("erro ao salvar histórico MQTT: %v\n", err)
			}
		}
	}

	// Marca retained inicial
	select {
	case c.Initial <- struct{}{}:
	default:
	}
}

func (c *Client) PublishToggle() error {
	payload := `{"command":"toggle"}`
	tok := c.client.Publish(c.cfg.MQTTTopicCommands, 1, false, payload)
	tok.Wait()
	return tok.Error()
}
