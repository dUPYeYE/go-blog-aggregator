-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id)
                       ON DELETE CASCADE
                       ON UPDATE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_fetched_at TIMESTAMP,

    UNIQUE (url)
);

-- +goose Down
DROP TABLE feeds;
