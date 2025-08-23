CREATE TABLE lists (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    board_id INTEGER NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for faster queries when finding lists by board
CREATE INDEX idx_lists_board_id ON lists(board_id);

-- Index for ordering lists by position within a board
CREATE INDEX idx_lists_board_position ON lists(board_id, position);