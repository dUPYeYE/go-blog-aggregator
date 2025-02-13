-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE posts;
