package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fentezi/translator/config"
	"github.com/fentezi/translator/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostgreSQLRepository struct {
	db  *sql.DB
	ctx context.Context
}

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

func New(ctx context.Context, cfg *config.Postgres) (*PostgreSQLRepository, error) {
	psql := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.Database,
	)

	db, err := sql.Open("postgres", psql)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgreSQLRepository{
		ctx: ctx,
		db:  db,
	}, nil
}

func (r *PostgreSQLRepository) Close() {
	_ = r.db.Close()
}

func (r *PostgreSQLRepository) Get(wordID uuid.UUID) (string, error) {
	const op = "repositories.PostgreSQLRepository.Get"

	query := `SELECT translation FROM words WHERE text = $1`

	var text string

	err := r.db.QueryRowContext(r.ctx, query, wordID).Scan(&text)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, ErrNotFound)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return text, nil
}

func (r *PostgreSQLRepository) Set(wordID uuid.UUID, key string, value string) error {
	const op = "repositories.PostgreSQLRepository.Set"

	query := `INSERT INTO words (word_id, text, translation) VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(r.ctx, query, wordID, key, value)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PostgreSQLRepository) Gets() ([]models.Word, error) {
	const op = "repositories.PostgreSQLRepository.Gets"

	query := `SELECT word_id, text, translation FROM words`
	rows, err := r.db.QueryContext(r.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		err := rows.Scan(&word.ID, &word.Word, &word.Translation)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		words = append(words, word)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return words, nil
}

func (r *PostgreSQLRepository) Delete(wordID uuid.UUID) error {
	const op = "repositories.PostgreSQLRepository.Delete"

	query := `DELETE FROM words WHERE word_id = $1`

	_, err := r.db.ExecContext(r.ctx, query, wordID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
