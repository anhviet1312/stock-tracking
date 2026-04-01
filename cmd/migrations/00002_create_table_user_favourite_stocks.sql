-- +goose Up
CREATE TABLE user_favourite_stocks (
    id         UUID PRIMARY KEY NOT NULL,
    user_id    UUID NOT NULL,
    symbol     VARCHAR(10) NOT NULL,
    status     VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT fk_ufs_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_ufs_stock FOREIGN KEY (symbol) REFERENCES stocks (symbol)
);

CREATE INDEX idx_user_favourite_stocks_user_id ON user_favourite_stocks (user_id);
CREATE INDEX idx_user_favourite_stocks_symbol ON user_favourite_stocks (symbol);

-- +goose Down
DROP TABLE user_favourite_stocks;
