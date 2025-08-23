CREATE TABLE cards (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    list_id INTEGER NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for faster queries when finding cards by list
CREATE INDEX idx_cards_list_id ON cards(list_id);

-- Index for ordering cards by position within a list
CREATE INDEX idx_cards_list_position ON cards(list_id, position);