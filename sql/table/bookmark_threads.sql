-- Create the bookmark_threads table
CREATE TABLE IF NOT EXISTS bookmark_threads (
    thread_id INTEGER,
    user_id INTEGER,
    PRIMARY KEY (user_id, thread_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (thread_id) REFERENCES threads(id)
);
