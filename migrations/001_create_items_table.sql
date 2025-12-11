-- +goose Up
CREATE TABLE IF NOT EXISTS items
(
    id         UUID PRIMARY KEY,
    type       VARCHAR(16) NOT NULL CHECK (type IN ('income', 'expense')),
    amount     INT         NOT NULL CHECK (amount > 0),
    date       TIMESTAMPTZ NOT NULL,
    category   VARCHAR(32) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_items_date ON items (date);
CREATE INDEX IF NOT EXISTS idx_items_type ON items (type);
CREATE INDEX IF NOT EXISTS idx_items_category ON items (category);
CREATE INDEX IF NOT EXISTS idx_items_date_type ON items (date, type);

-- +goose Down
DROP TABLE IF EXISTS items;
