-- Create the reviews table
CREATE TABLE IF NOT EXISTS reviews (
    moduleId VARCHAR(20),
    userId INTEGER,
    workload INTEGER NOT NULL CHECK (workload >= 1 AND workload <= 5),
    expectedGrade VARCHAR(3),
    actualGrade VARCHAR(3),
    difficulty INTEGER NOT NULL CHECK (difficulty >= 1 AND difficulty <= 5),
    semesterTaken CHAR(20) NOT NULL,
    lecturer VARCHAR(256),
    content TEXT NOT NULL,
    suggestion TEXT,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (userId, moduleId),
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (moduleId) REFERENCES modules(id)
);
