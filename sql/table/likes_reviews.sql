CREATE TABLE IF NOT EXISTS likes_reviews (
    userId INTEGER,
    moduleId VARCHAR(20),
    state BOOLEAN NOT NULL,
    PRIMARY KEY (userId, moduleId),
    FOREIGN KEY (moduleId, userid) REFERENCES reviews(moduleid, userid)
  );