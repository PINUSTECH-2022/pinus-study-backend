CREATE TABLE IF NOT EXISTS comments_review (
    id INTEGER PRIMARY KEY,
    content TEXT NOT NULL,
    userId INTEGER,
	moduleId varchar(20),
    authorId INTEGER,
    parentId INTEGER,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
	likes INTEGER not NULL check (likes >= 0),
	dislikes INTEGER NOT null check (dislikes >= 0),
    FOREIGN KEY (userId, moduleId) REFERENCES reviews(userId, moduleId),
    FOREIGN KEY (authorId) REFERENCES users(id),
    FOREIGN KEY (parentId) REFERENCES comments_review(id)
);