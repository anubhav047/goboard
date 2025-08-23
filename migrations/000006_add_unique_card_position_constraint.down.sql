-- Remove the unique constraint
ALTER TABLE cards DROP CONSTRAINT IF EXISTS unique_card_position_per_list;