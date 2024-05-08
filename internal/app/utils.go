package app

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"theater/config"
)

func ConnectToPostgres(cfg *config.Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.PG.Host, cfg.PG.User, cfg.PG.Password, cfg.PG.DBName, cfg.PG.Port, cfg.PG.SSLMode)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error occurred while connecting to db: %w", err)
	}

	return db, nil
}
