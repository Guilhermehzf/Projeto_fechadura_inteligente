package db

import (
    "context"
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"
)

type User struct {
    ID       string
    Email    string
    Password string
    Name     string
}

// SELECT * FROM users WHERE email=...
func GetUserByEmail(ctx context.Context, pool *sql.DB, email string) (*User, error) {
    row := pool.QueryRowContext(ctx,
        "SELECT id, email, password, name FROM users WHERE email=$1 LIMIT 1", email)

    u := User{}
    if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Name); err != nil {
        return nil, err
    }
    return &u, nil
}