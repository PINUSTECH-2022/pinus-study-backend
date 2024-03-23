-- Create the comments table
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY,
    content TEXT NOT NULL,
    threadId INTEGER,
    authorId INTEGER,
    parentId INTEGER,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    likes_count INTEGER NOT NULL DEFAULT 0 CHECK ((likes_count) >= 0),
    dislikes_count INTEGER NOT NULL DEFAULT 0 CHECK ((dislikes_count) >= 0),
    comments_count INTEGER NOT NULL DEFAULT 0 CHECK ((comments_count) >= 0),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (threadId) REFERENCES threads(id),
    FOREIGN KEY (authorId) REFERENCES users(id),
    FOREIGN KEY (parentId) REFERENCES comments(id)
);
