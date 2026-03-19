-- +goose Up
CREATE TABLE users
(
    "id"            UUID PRIMARY KEY,
    "first_name"    TEXT,
    "last_name"     TEXT,
    "username"      TEXT NOT NULL,
    "password"      TEXT NOT NULL,
    "email"         TEXT NOT NULL,
    "is_active"     BOOLEAN NOT NULL DEFAULT FALSE,
    "updated_at"    TIMESTAMP WITH TIME ZONE NOT NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_users__username ON "users" ("username");
CREATE INDEX idx_users__email ON "users" ("email");
CREATE INDEX idx_users__created_at ON "users" ("created_at");

CREATE TABLE user_codes
(
    "id"            UUID PRIMARY KEY,
    "user_id"       UUID NOT NULL,
    "value"         TEXT NOT NULL,
    "status"        TEXT NOT NULL,
    "type"          TEXT NOT NULL,
    "expired_time"  TIMESTAMP WITH TIME ZONE NOT NULL,
    "updated_at"    TIMESTAMP WITH TIME ZONE NOT NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT fk__user_codes__users FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

CREATE INDEX idx_user_codes__user_id_type_created_at_status ON "user_codes" ("user_id", "type", "created_at", "status");
CREATE INDEX idx_user_codes__user_id_type_expired_time ON "user_codes" ("user_id", "type", "expired_time");
CREATE INDEX idx_user_codes__user_id_type_status ON "user_codes" ("user_id", "type", "status");

-- +goose Down
DROP TABLE "user_codes";
DROP TABLE "users";