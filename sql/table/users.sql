-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username VARCHAR(100) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    salt VARCHAR(200) NOT NULL UNIQUE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE
);
