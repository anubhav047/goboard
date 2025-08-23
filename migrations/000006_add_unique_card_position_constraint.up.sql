-- Add unique constraint for card positions within a list
ALTER TABLE cards ADD CONSTRAINT unique_card_position_per_list UNIQUE (list_id, position);