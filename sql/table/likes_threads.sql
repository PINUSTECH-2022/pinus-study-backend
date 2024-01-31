-- Create the likes_threads table
CREATE TABLE IF NOT EXISTS likes_threads (
    userId INTEGER REFERENCES users,
    threadId INTEGER REFERENCES threads,
    state BOOLEAN NOT NULL,
    PRIMARY KEY (userId, threadId)
);
