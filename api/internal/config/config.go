package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
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

	JWTSecret  string
	JWTExpiry  time.Duration

	DB         *sql.DB
	MaxHistory int
}

func Load() *Config {
	_ = godotenv.Load()

	// Pega a string de conexão do .env
	dsn := must("POSTGRES_ENDPOINT")

	// Abre conexão com Postgres
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Erro ao abrir conexão com Postgres: %v", err)
	}

	// Testa conexão
	if err := db.Ping(); err != nil {
		log.Fatalf("Banco de dados inacessível: %v", err)
	}

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
		DB:                db,
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
		if n > 0 {
			return n
		}
	}
	return def
}

func getDuration(key string, defMinutes int) time.Duration {
	if v := os.Getenv(key); v != "" {
		var n int
		_, _ = fmt.Sscanf(v, "%d", &n)
		if n > 0 {
			return time.Duration(n) * time.Minute
		}
	}
	return time.Duration(defMinutes) * time.Minute
}
