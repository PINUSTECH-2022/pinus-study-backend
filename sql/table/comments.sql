-- Create the comments table
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY,
    content TEXT NOT NULL,
    threadId INTEGER,
    authorId INTEGER,
    parentId INTEGER,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (threadId) REFERENCES threads(id),
    FOREIGN KEY (authorId) REFERENCES users(id),
    FOREIGN KEY (parentId) REFERENCES comments(id)
);
