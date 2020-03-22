CREATE TABLE users(
    "id"            BIGSERIAL PRIMARY KEY,
    "first_name"    TEXT NOT NULL,
    "last_name"     TEXT,
    "username"      TEXT,
    "language_code" TEXT,
    "created_at"    TIMESTAMP NOT NULL DEFAULT NOW()
);
