package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
	dsn := os.Getenv("POSTGRES_ENDPOINT")
	if dsn == "" {
		log.Fatal("Variável POSTGRES_ENDPOINT não definida")
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Erro parseando config: %v", err)
	}

	// Timeout de conexão
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	Pool, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("Erro conectando ao PostgreSQL: %v", err)
	}

	if err := Pool.Ping(ctx); err != nil {
		log.Fatalf("Erro no ping ao PostgreSQL: %v", err)
	}

	log.Println("Conexão com PostgreSQL estabelecida com sucesso!")
}
