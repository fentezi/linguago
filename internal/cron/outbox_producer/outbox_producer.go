package outboxproducer

import (
	"context"
	"fmt"
	"github.com/fentezi/translator/internal/entities"
	"github.com/fentezi/translator/internal/kafka"
	"github.com/fentezi/translator/internal/repositories"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log/slog"
)

type OutboxProducer struct {
	log      *slog.Logger
	db       *repositories.PostgreSQLRepository
	producer *kafka.Producer
}

func New(
	log *slog.Logger, producer *kafka.Producer, db *repositories.PostgreSQLRepository,
) *OutboxProducer {
	return &OutboxProducer{
		log:      log,
		producer: producer,
		db:       db,
	}
}

func (o *OutboxProducer) ProduceMessage(ctx context.Context, topic string) (err error) {
	const op = "outboxProducer.ProduceMessage"
	db := o.db.DB()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
		SELECT outbox.event_id, words.text, words.translation
		FROM outbox
		JOIN words ON outbox.word_id = words.word_id
		WHERE sent IS FALSE
		ORDER BY outbox.word_id DESC
		LIMIT 50
	`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			return
		}
	}()

	var eventIDs []uuid.UUID
	for rows.Next() {
		msg := entities.OutboxMessage{}
		if err := rows.Scan(&msg.EventID, &msg.Word, &msg.Translation); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		eventIDs = append(eventIDs, msg.EventID)
		err = o.producer.Produce(
			entities.KafkaMessage{
				TopicPartition: topic,
				Value: []byte(fmt.Sprintf(
					`{"event_id": "%v", "word": "%s", "translation": "%s"}`, msg.EventID, msg.Word,
					msg.Translation,
				)),
			},
		)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if len(eventIDs) == 0 {
		return tx.Rollback()
	}

	query = `UPDATE outbox SET sent = TRUE WHERE event_id = ANY($1)`
	_, err = tx.ExecContext(ctx, query, pq.Array(eventIDs))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (o *OutboxProducer) Close() {
	o.producer.Close()
}
