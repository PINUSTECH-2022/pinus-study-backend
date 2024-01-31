-- Create the tokens table
CREATE TABLE IF NOT EXISTS tokens (
    userId INTEGER PRIMARY KEY,
    token VARCHAR(256) NOT NULL,
    FOREIGN KEY (userId) REFERENCES users(id)
);
