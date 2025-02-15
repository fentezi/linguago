package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/fentezi/translator/config"

	_ "github.com/lib/pq"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Repository interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
}

func NewPostgreSQL(cfg *config.Config) (*sql.DB, error) {
	psql := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	fmt.Println(psql)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
