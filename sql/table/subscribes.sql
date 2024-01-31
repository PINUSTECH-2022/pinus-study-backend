-- Create the subscribes table
CREATE TABLE IF NOT EXISTS subscribes (
    moduleId VARCHAR(10),
    userId INTEGER,
    PRIMARY KEY (userId, moduleId),
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (moduleId) REFERENCES modules(id)
);
