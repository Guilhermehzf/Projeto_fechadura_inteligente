package config

import (
	"log"
	"os"
	"time"
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr          string

	MQTTBroker        string
	MQTTClientID      string
	MQTTUser          string
	MQTTPass          string
	MQTTTopicState    string
	MQTTTopicCommands string

	JWTSecret         string
	JWTExpiry         time.Duration

	MaxHistory        int
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		HTTPAddr:          getenv("HTTP_ADDR", "0.0.0.0:8088"),
		MQTTBroker:        must("MQTT_BROKER"),
		MQTTClientID:      must("MQTT_CLIENT_ID"),
		MQTTUser:          must("MQTT_USER"),
		MQTTPass:          must("MQTT_PASS"),
		MQTTTopicState:    must("MQTT_TOPIC_STATE"),
		MQTTTopicCommands: getenv("MQTT_TOPIC_COMMANDS", "fechadura/comandos"),
		JWTSecret:         getenv("JWT_SECRET", "dev-secret-change-me"),
		JWTExpiry:         getDuration("JWT_EXP_MINUTES", 60),
		MaxHistory:        getInt("MAX_HISTORY", 200),
	}

	return cfg
}

func must(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("variável de ambiente obrigatória não definida: %s", key)
	}
	return v
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		var n int
		_, _ = fmt.Sscanf(v, "%d", &n)
		if n > 0 { return n }
	}
	return def
}

func getDuration(key string, defMinutes int) time.Duration {
	if v := os.Getenv(key); v != "" {
		var n int
		_, _ = fmt.Sscanf(v, "%d", &n)
		if n > 0 { return time.Duration(n) * time.Minute }
	}
	return time.Duration(defMinutes) * time.Minute
}
