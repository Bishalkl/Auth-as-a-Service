-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL, -- Example: 'admin', 'user', etc.
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(user_id, role)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_roles;
-- +goose StatementEnd
