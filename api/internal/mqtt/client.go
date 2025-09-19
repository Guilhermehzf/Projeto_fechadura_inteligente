package mqtt

import (
	"log"
	"os"
	"strconv"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"

	"smartlock/internal/config"
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
	c.store.UpdateFromMQTT(msg.Topic(), msg.Qos(), msg.Retained(), msg.Payload())

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
