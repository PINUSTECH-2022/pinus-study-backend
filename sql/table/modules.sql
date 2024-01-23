-- Create the modules table
CREATE TABLE IF NOT EXISTS modules (
    id VARCHAR(20) PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL
);
