-- Create the threads table
CREATE TABLE IF NOT EXISTS threads (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    moduleId VARCHAR(20) REFERENCES modules,
    authorId INTEGER REFERENCES users,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    likes_count INTEGER NOT NULL DEFAULT 0 CHECK ((likes_count) >= 0),
    dislikes_count INTEGER NOT NULL DEFAULT 0 CHECK ((dislikes_count) >= 0),
    comments_count INTEGER NOT NULL DEFAULT 0 CHECK ((comments_count) >= 0),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);
