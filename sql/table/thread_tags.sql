-- Create the thread_tags table
CREATE TABLE IF NOT EXISTS thread_tags (
    threadId INTEGER NOT NULL REFERENCES threads,
    tagId INTEGER NOT NULL REFERENCES tags,
    PRIMARY KEY (threadId, tagId)
);
