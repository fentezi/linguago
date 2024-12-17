package repositories

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPostgreSQLRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	key := "hello"
	expectedTranslation := "привет"

	ctx := context.Background()
	repo := NewPostgreSQLRepository(db, ctx)

	mock.ExpectQuery(`SELECT translation FROM words WHERE text = \$1`).
		WithArgs(key).
		WillReturnRows(sqlmock.NewRows([]string{"translation"}).AddRow(expectedTranslation))

	translation, err := repo.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, expectedTranslation, translation)

	mock.ExpectQuery(`SELECT translation FROM words WHERE text = \$1`).
		WithArgs(key).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.Get(key)
	assert.ErrorIs(t, err, ErrNotFound)

}

func TestPostgreSQLRepository_Set(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repo := NewPostgreSQLRepository(db, ctx)

	key := "hello"
	value := "привет"

	mock.ExpectExec(`INSERT INTO words \(text, translation\) VALUES \(\$1, \$2\)`).
		WithArgs(key, value).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Set(key, value)
	assert.NoError(t, err)
}

func TestPostgreSQLRepository_Gets(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repo := NewPostgreSQLRepository(db, ctx)

	query := regexp.QuoteMeta(`SELECT word_id, text, translation FROM words`)

	rows := sqlmock.NewRows([]string{"word_id", "text", "translation"}).
		AddRow(1, "hello", "привет").
		AddRow(2, "world", "мир")

	mock.ExpectQuery(query).WillReturnRows(rows)

	words, err := repo.Gets()
	assert.NoError(t, err)
	assert.Len(t, words, 2)

	assert.Equal(t, 1, words[0].ID)
	assert.Equal(t, "hello", words[0].Text)
	assert.Equal(t, "привет", words[0].Translation)

	assert.Equal(t, 2, words[1].ID)
	assert.Equal(t, "world", words[1].Text)
	assert.Equal(t, "мир", words[1].Translation)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgreSQLRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repo := NewPostgreSQLRepository(db, ctx)

	key := "hello"
	query := regexp.QuoteMeta(`DELETE FROM words WHERE text = $1`)

	mock.ExpectExec(query).
		WithArgs(key).
		WillReturnResult(sqlmock.NewResult(0, 1))
	err = repo.Delete(key)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
