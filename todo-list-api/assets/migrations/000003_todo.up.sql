CREATE TABLE IF NOT EXISTS todo (
    id SERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    user_id SERIAL NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
