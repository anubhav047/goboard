-- name: CreateUser :one
INSERT INTO users (
  name,
  email,
  hashed_password
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- ================================
-- BOARD QUERIES
-- ================================

-- name: CreateBoard :one
INSERT INTO boards (
  name,
  description,
  created_by
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetBoardByID :one
SELECT * FROM boards
WHERE id = $1 LIMIT 1;

-- name: GetBoardsByUser :many
SELECT * FROM boards
WHERE created_by = $1
ORDER BY created_at DESC;

-- name: UpdateBoard :one
UPDATE boards
SET name = $1, description = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: DeleteBoard :exec
DELETE FROM boards
WHERE id = $1;

-- ================================
-- LIST QUERIES
-- ================================

-- name: CreateList :one
INSERT INTO lists (
  name,
  board_id,
  position
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetListsByBoard :many
SELECT * FROM lists
WHERE board_id = $1
ORDER BY position ASC;

-- name: GetListByID :one
SELECT * FROM lists
WHERE id = $1 LIMIT 1;

-- name: UpdateList :one
UPDATE lists
SET name = $1, position = $2, updated_at = NOW()
where id = $3
RETURNING *;

-- name: DeleteList :exec
DELETE FROM lists
where id = $1;

-- ================================
-- CARD QUERIES
-- ================================

-- name: CreateCard :one
INSERT INTO cards (
  title,
  description,
  list_id,
  position
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetCardsByList :many
SELECT * FROM cards
WHERE list_id = $1
ORDER BY position ASC;

-- name: GetCardByID :one
SELECT * FROM cards
WHERE id = $1 LIMIT 1;

-- name: UpdateCard :one
UPDATE cards
SET title = $1, description = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: MoveCard :one
UPDATE cards
SET list_id = $1, position = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: DeleteCard :exec
DELETE FROM cards
WHERE id = $1;