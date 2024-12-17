package repositories

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository(t *testing.T) {
	miniRedis, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	defer miniRedis.Close()

	client := redis.NewClient(
		&redis.Options{
			Addr: miniRedis.Addr(),
		},
	)

	repo := NewRedisRepository(client, context.Background())

	t.Run("Set and Get", func(t *testing.T) {
		err := repo.Set("key1", "value1")
		assert.NoError(t, err)

		value, err := repo.Get("key1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", value)
	})

	t.Run("Get non-existing key", func(t *testing.T) {
		_, err := repo.Get("non-existing-key")
		assert.Error(t, err)
	})

	t.Run("Delete key", func(t *testing.T) {
		err := repo.Set("key2", "value2")
		assert.NoError(t, err)

		err = repo.Delete("key2")
		assert.NoError(t, err)

		_, err = repo.Get("key2")
		assert.Error(t, err)
	})

	t.Run("Get and twice Set", func(t *testing.T) {
		err := repo.Set("key3", "value3")
		assert.NoError(t, err)

		value, err := repo.Get("key3")
		assert.NoError(t, err)
		assert.Equal(t, "value3", value)

		err = repo.Set("key3", "value33")
		assert.NoError(t, err)

		value, err = repo.Get("key3")
		assert.NoError(t, err)
		assert.Equal(t, "value33", value)
	})
}
