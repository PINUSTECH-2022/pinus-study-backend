-- Create the tags table
CREATE TABLE IF NOT EXISTS tags (
    tagId SERIAL PRIMARY KEY,
    tagDesc TEXT NOT NULL
);
