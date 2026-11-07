-- Queries for keeping contacts in Postgres with sqlc.
--
-- The template runs with the in-memory Store so nothing uses these yet.
-- To switch: add the Postgres driver (go get github.com/jackc/pgx/v5), set
-- DSN in .env, run `make generate` to produce the query package from these
-- queries, and implement the Store methods with it. See the docs:
-- Persisting Data and sqlc.
--
-- The contacts table (internal/db/migrations/migrations.sql) stores each
-- contact as a JSON blob with generated columns for the fields we search,
-- sort and index on.

-- name: CreateContact :exec
INSERT INTO contacts (id, blob) VALUES ($1, $2);

-- name: GetContact :one
SELECT blob FROM contacts WHERE id = $1;

-- name: UpdateContact :exec
UPDATE contacts SET blob = $2 WHERE id = $1;

-- name: DeleteContact :exec
DELETE FROM contacts WHERE id = $1;

-- name: SearchContacts :many
SELECT blob FROM contacts
WHERE last_name ILIKE '%' || sqlc.arg(term)::text || '%'
   OR first_name ILIKE '%' || sqlc.arg(term)::text || '%'
   OR email ILIKE '%' || sqlc.arg(term)::text || '%'
ORDER BY lower(last_name), lower(first_name)
LIMIT sqlc.arg(page_size)::int OFFSET sqlc.arg(page_offset)::int;

-- name: CountContacts :one
SELECT count(*) FROM contacts
WHERE last_name ILIKE '%' || sqlc.arg(term)::text || '%'
   OR first_name ILIKE '%' || sqlc.arg(term)::text || '%'
   OR email ILIKE '%' || sqlc.arg(term)::text || '%';

-- name: EmailTaken :one
SELECT EXISTS (SELECT 1 FROM contacts WHERE lower(email) = lower(sqlc.arg(email)::text) AND id <> sqlc.arg(except_id)::text);
