-- Migrations run at server start (see internal/db/db.go) and must be idempotent.
--
-- Objects are stored as JSON blobs; fields that are searched, sorted or
-- indexed are exposed as generated columns computed from the blob. See the
-- "Sql DB, NoSql Schema" blog post and the sqlc docs page.

CREATE TABLE IF NOT EXISTS contacts (
    id         TEXT PRIMARY KEY,
    blob       JSONB NOT NULL,
    first_name TEXT GENERATED ALWAYS AS (blob->>'first_name') STORED,
    last_name  TEXT GENERATED ALWAYS AS (blob->>'last_name') STORED,
    email      TEXT GENERATED ALWAYS AS (blob->>'email') STORED
);
CREATE INDEX IF NOT EXISTS contacts_last_name_idx ON contacts (lower(last_name), lower(first_name));
CREATE UNIQUE INDEX IF NOT EXISTS contacts_email_idx ON contacts (lower(email));
