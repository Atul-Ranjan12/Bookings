package dbrepo

import (
	"database/sql"
	"github/Atul-Ranjan12/booking/internal/config"
	"github/Atul-Ranjan12/booking/internal/repository"
)

type PostgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &PostgresDBRepo{
		App: a,
		DB:  conn,
	}
}
