-- +goose Up
ALTER TABLE "users" ADD COLUMN "telegram_chat_id" TEXT;

CREATE TABLE "user_stock_alerts"
(
    "id"                 UUID PRIMARY KEY NOT NULL,
    "user_id"            UUID NOT NULL,
    "symbol"             VARCHAR(10) NOT NULL,
    "high_threshold"     NUMERIC(15, 4),
    "low_threshold"      NUMERIC(15, 4),
    "last_notified_at"   TIMESTAMP WITH TIME ZONE,
    "is_active"          BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at"         TIMESTAMP WITH TIME ZONE NOT NULL,
    "updated_at"         TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT fk_usa_user FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
    CONSTRAINT fk_usa_stock FOREIGN KEY ("symbol") REFERENCES "stocks" ("symbol")
);

CREATE INDEX idx_user_stock_alerts_user_id ON "user_stock_alerts" ("user_id");
CREATE INDEX idx_user_stock_alerts_symbol ON "user_stock_alerts" ("symbol");

-- +goose Down
DROP TABLE "user_stock_alerts";
ALTER TABLE "users" DROP COLUMN "telegram_chat_id";
