package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type History struct {
	ID        string
	Action    string
	Timestamp time.Time
	Method    string
}

// Insere novo registro no histórico
func InsertHistory(ctx context.Context, pool *sql.DB, action, method string) error {
	_, err := pool.ExecContext(ctx,
		`INSERT INTO lock_history (action, method) VALUES ($1, $2)`,
		action, method,
	)
	return err
}

// Lista histórico (limit)
func GetHistory(ctx context.Context, pool *sql.DB, limit int) ([]History, error) {
	rows, err := pool.QueryContext(ctx,
		`SELECT id, action, timestamp, method
		 FROM lock_history
		 ORDER BY timestamp DESC
		 LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hist []History
	for rows.Next() {
		var h History
		if err := rows.Scan(&h.ID, &h.Action, &h.Timestamp, &h.Method); err != nil {
			return nil, err
		}
		hist = append(hist, h)
	}
	return hist, rows.Err()
}
