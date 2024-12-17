package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/fentezi/translator/config"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repository interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
}

func NewRedis(cfg *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.NumberDB,
	})

	return client
}

func NewPostgreSQL(cfg *config.Config) (*sql.DB, error) {
	psql := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)
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
