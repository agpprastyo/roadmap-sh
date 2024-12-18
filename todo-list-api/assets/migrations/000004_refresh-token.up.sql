CREATE TABLE refresh_tokens (
                                id SERIAL PRIMARY KEY,
                                user_id SERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                token TEXT NOT NULL UNIQUE,
                                expiry TIMESTAMP WITH TIME ZONE NOT NULL,
                                created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                revoked BOOLEAN DEFAULT false,
                                CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_refresh_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_user ON refresh_tokens(user_id);
