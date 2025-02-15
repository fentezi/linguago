package repositories

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
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

	id := uuid.New()
	key := "hello"
	value := "привет"

	mock.ExpectExec(`INSERT INTO words \(word_id, text, translation\) VALUES \(\$1, \$2, \$3\)`).
		WithArgs(id, key, value).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Set(id, key, value)
	assert.NoError(t, err)

}

func TestPostgreSQLRepository_Gets(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repo := NewPostgreSQLRepository(db, ctx)

	id1 := uuid.New()
	id2 := uuid.New()

	query := regexp.QuoteMeta(`SELECT word_id, text, translation FROM words`)

	rows := sqlmock.NewRows([]string{"word_id", "text", "translation"}).
		AddRow(id1, "hello", "привет").
		AddRow(id2, "world", "мир")

	mock.ExpectQuery(query).WillReturnRows(rows)

	words, err := repo.Gets()
	assert.NoError(t, err)
	assert.Len(t, words, 2)

	assert.Equal(t, id1, words[0].ID)
	assert.Equal(t, "hello", words[0].Word)
	assert.Equal(t, "привет", words[0].Translation)

	assert.Equal(t, id2, words[1].ID)
	assert.Equal(t, "world", words[1].Word)
	assert.Equal(t, "мир", words[1].Translation)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgreSQLRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	repo := NewPostgreSQLRepository(db, ctx)

	key := "1"
	query := regexp.QuoteMeta(`DELETE FROM words WHERE word_id = $1`)

	mock.ExpectExec(query).
		WithArgs(key).
		WillReturnResult(sqlmock.NewResult(0, 1))
	err = repo.Delete(key)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
