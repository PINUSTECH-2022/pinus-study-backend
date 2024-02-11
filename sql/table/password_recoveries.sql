-- Create the password_recoveries table
CREATE TABLE IF NOT EXISTS password_recoveries (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    secret_code CHAR(32) NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expired_at TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '15 minutes',
    FOREIGN KEY (user_id) REFERENCES users(id)
);
