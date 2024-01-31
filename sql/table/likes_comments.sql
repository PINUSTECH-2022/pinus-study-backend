-- Create the likes_comments table
CREATE TABLE IF NOT EXISTS likes_comments (
    userId INTEGER,
    commentId INTEGER,
    state BOOLEAN NOT NULL,
    PRIMARY KEY (userId, commentId),
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (commentId) REFERENCES comments(id)
);
