-- +goose Up
CREATE TABLE users(
    id UUID primary key,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    username TEXT unique not null,
    email TEXT unique not null,
    hashed_password TEXT not null
);

-- +goose Down
DROP TABLE users;