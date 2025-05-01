CREATE TABLE IF NOT EXISTS "outbox"
(
    outbox_id serial PRIMARY KEY,
    event_id  UUID                            NOT NULL,
    word_id   UUID REFERENCES words (word_id) ON DELETE CASCADE NOT NULL,
    sent      BOOLEAN DEFAULT FALSE,
    UNIQUE (event_id)
);
