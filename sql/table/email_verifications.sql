-- Create the email_verifications table
CREATE TABLE IF NOT EXISTS email_verifications (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    email TEXT NOT NULL,
    secret_code CHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expired_at TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '15 minutes',
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (email) REFERENCES users(email)
);
