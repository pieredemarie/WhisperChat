CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    room_id TEXT NOT NULL,
    display_name TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_messages_room_created ON messages (room_id, created_at);
