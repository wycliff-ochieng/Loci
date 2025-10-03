-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd


CREATE TABLE loci (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL CHECK (length(message) > 0 AND length(message) <= 280),
    -- GEOGRAPHY type is better than GEOMETRY for lat/lon calculations in meters/km
    location GEOGRAPHY(Point, 4326) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    view_count INT NOT NULL DEFAULT 0,
    replies_count INT NOT NULL DEFAULT 0,
    visibility_score FLOAT NOT NULL DEFAULT 0.0
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
