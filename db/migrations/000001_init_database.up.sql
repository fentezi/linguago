CREATE TABLE IF NOT EXISTS "words" (
    word_id SERIAL PRIMARY KEY,
    text VARCHAR(255) NOT NULL,
    translation TEXT NOT NULL,
    UNIQUE (text, translation)
);