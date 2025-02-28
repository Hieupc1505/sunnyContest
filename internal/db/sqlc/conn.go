package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	config "go-rest-api-boilerplate/configs/app"
	"log"
)

var (
	DB_NOT_FOUND = pgx.ErrNoRows
)

// Connect to postgrest db
func Conn() *pgxpool.Pool {
	connPool, err := pgxpool.New(context.Background(), config.Envs.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to PostgreSQL database: ", err)
	}
	return connPool
}
