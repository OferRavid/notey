-- +goose Up
CREATE TABLE notes(
    id UUID primary key,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    title TEXT not null,
    content TEXT not null,
    user_id UUID not null REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE notes;