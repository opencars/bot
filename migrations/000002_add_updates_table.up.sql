CREATE TABLE updates(
    "id"            BIGSERIAL PRIMARY KEY,
    "user_id"       BIGSERIAL REFERENCES users("id"),
    "text"          TEXT NOT NULL,
    "time"          TIMESTAMP NOT NULL,
    "created_at"    TIMESTAMP NOT NULL DEFAULT NOW()
);
