--Create all tables
BEGIN TRANSACTION;
  CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username VARCHAR(100) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    salt VARCHAR(200) NOT NULL UNIQUE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE
  );


  CREATE TABLE IF NOT EXISTS modules (
    id VARCHAR(20) PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL
  );


  CREATE TABLE IF NOT EXISTS threads (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    moduleId VARCHAR(20) REFERENCES modules,
    authorId INTEGER REFERENCES users,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    likes_count INTEGER NOT NULL DEFAULT 0 CHECK ((likes_count) >= 0),
    dislikes_count INTEGER NOT NULL DEFAULT 0 CHECK ((dislikes_count) >= 0),
    comments_count INTEGER NOT NULL DEFAULT 0 CHECK ((comments_count) >= 0),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
  );


  CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY,
    content TEXT NOT NULL,
    threadId INTEGER,
    authorId INTEGER,
    parentId INTEGER,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (threadId) REFERENCES threads(id),
    FOREIGN KEY (authorId) REFERENCES users(id),
    FOREIGN KEY (parentId) REFERENCES comments(id)
  );


  CREATE TABLE IF NOT EXISTS subscribes (
    moduleId VARCHAR(10),
    userId INTEGER,
    PRIMARY KEY (userId, moduleId),
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (moduleId) REFERENCES modules(id)
  );


  CREATE TABLE IF NOT EXISTS tokens (
    userId INTEGER PRIMARY KEY,
    token VARCHAR(256) NOT NULL,
    FOREIGN KEY (userId) REFERENCES users(id)
  );


  CREATE TABLE IF NOT EXISTS likes_comments (
    userId INTEGER,
    commentId INTEGER,
    state BOOLEAN NOT NULL,
    PRIMARY KEY (userId, commentId),
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (commentId) REFERENCES comments(id)
  );


  CREATE TABLE IF NOT EXISTS likes_threads (
    userId INTEGER REFERENCES users,
    threadId INTEGER REFERENCES threads,
    state BOOLEAN NOT NULL,
    PRIMARY KEY (userId, threadId)
  );


  CREATE TABLE IF NOT EXISTS tags (
    tagId SERIAL PRIMARY KEY,
    tagDesc TEXT NOT NULL
  );


  CREATE TABLE IF NOT EXISTS thread_tags (
    threadId INTEGER NOT NULL REFERENCES threads,
    tagId INTEGER NOT NULL REFERENCES tags,
    PRIMARY KEY (threadId, tagId)
  );


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


  CREATE TABLE IF NOT EXISTS bookmark_threads (
    thread_id INTEGER,
    user_id INTEGER,
    PRIMARY KEY (user_id, thread_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (thread_id) REFERENCES threads(id)
  );


  CREATE TABLE IF NOT EXISTS email_verifications (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    email TEXT NOT NULL,
    secret_code CHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expired_at TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '15 minutes',
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (email) REFERENCES users(email)
  );
COMMIT;


-- Create Procedure update_all_thread_count
BEGIN TRANSACTION;
  DROP PROCEDURE IF EXISTS update_all_thread_count;


  CREATE OR REPLACE PROCEDURE update_all_thread_count()
  LANGUAGE plpgsql
  AS $$
  BEGIN
    UPDATE threads t
    SET likes_count = (SELECT COUNT(*) FROM likes_threads l WHERE l.threadid = t.id AND l.state),
    dislikes_count = (SELECT COUNT(*) FROM likes_threads l WHERE l.threadid = t.id AND NOT l.state),
    comments_count = (SELECT COUNT(*) FROM comments c WHERE c.threadid = t.id AND NOT c.is_deleted);
  END;
  $$;
COMMIT;


-- Create Trigger thread_likes_count_update_trigger and Function update_thread_likes_count
BEGIN TRANSACTION;
  DROP TRIGGER IF EXISTS thread_likes_count_update_trigger ON likes_threads;
  DROP FUNCTION IF EXISTS update_thread_likes_count;


  CREATE OR REPLACE FUNCTION update_thread_likes_count()
  RETURNS TRIGGER AS $$
  BEGIN
    IF TG_OP = 'INSERT' THEN
    UPDATE threads SET likes_count = likes_count + 1 WHERE NEW.threadid = id AND NEW.state;
    UPDATE threads SET dislikes_count = dislikes_count + 1 WHERE NEW.threadid = id AND NOT NEW.state;
    ELSIF TG_OP = 'DELETE' THEN
    UPDATE threads SET likes_count = likes_count - 1 WHERE OLD.threadid = id AND OLD.state;
    UPDATE threads SET dislikes_count = dislikes_count - 1 WHERE OLD.threadid = id AND NOT OLD.state;
    END IF;
    RETURN NULL;
  END;
  $$ LANGUAGE plpgsql;


  CREATE TRIGGER thread_likes_count_update_trigger
  AFTER INSERT OR DELETE ON likes_threads
  FOR EACH ROW
  EXECUTE FUNCTION update_thread_likes_count();
COMMIT;


-- Create Trigger thread_comments_count_update_trigger and Function update_thread_comments_count
BEGIN TRANSACTION;
  DROP TRIGGER IF EXISTS thread_comments_count_update_trigger ON comments;
  DROP FUNCTION IF EXISTS update_thread_comments_count;


  CREATE OR REPLACE FUNCTION update_thread_comments_count()
  RETURNS trigger
  LANGUAGE 'plpgsql'
  AS $$
  BEGIN
    IF TG_OP = 'INSERT' THEN
    UPDATE threads SET comments_count = comments_count + 1 WHERE NEW.threadid = id;
    ELSIF TG_OP = 'DELETE' THEN
    UPDATE threads SET comments_count = comments_count - 1 WHERE OLD.threadid = id;
    END IF;
    RETURN NULL;
  END;
  $$;


  CREATE OR REPLACE TRIGGER thread_comments_count_update_trigger
  AFTER INSERT OR DELETE ON comments
  FOR EACH ROW
  EXECUTE FUNCTION update_thread_comments_count();
COMMIT;


-- Create Procedure signup
BEGIN TRANSACTION;
  DROP PROCEDURE IF EXISTS signup;


  CREATE OR REPLACE PROCEDURE signup(candidate_username VARCHAR, candidate_email VARCHAR, password VARCHAR, salt VARCHAR, secret_code CHAR(32),
                     OUT user_id INT, OUT email_id INT, OUT is_email_exist BOOLEAN, OUT is_username_exist BOOLEAN)
  LANGUAGE plpgsql
  AS $$
  BEGIN
    -- Check if email already exist
    SELECT COUNT(*) > 0 INTO is_email_exist FROM users WHERE email = candidate_email;


    -- Check if username already exist
    SELECT COUNT(*) > 0 INTO is_username_exist FROM users WHERE username = candidate_username;


    IF NOT is_email_exist AND NOT is_username_exist THEN
      -- Insert user into users table
      INSERT INTO users (id, email, username, password, salt)
      VALUES ((SELECT COUNT(*) FROM users) + 1, candidate_email, candidate_username, password, salt)
      RETURNING id INTO user_id;


      -- Store secret code to email_verifications table
      INSERT INTO email_verifications (id, user_id, email, secret_code)
      VALUES ((SELECT COUNT(*) FROM email_verifications) + 1, user_id, candidate_email, secret_code)
      RETURNING id INTO email_id;
    END IF;


    -- Change user_id to -1 instead of NULL
    user_id := COALESCE(user_id, -1);


    -- Change email_id to -1 instead of NULL
    email_id := COALESCE(email_id, -1);
  END;
  $$;
COMMIT;


-- Create Procedure make_verification
BEGIN TRANSACTION;
  DROP PROCEDURE IF EXISTS make_verification;


  CREATE OR REPLACE PROCEDURE make_verification(user_id INT, secret_code VARCHAR,
                          OUT exist BOOLEAN, OUT email_id INT, OUT user_email VARCHAR, OUT uname VARCHAR, OUT verified BOOLEAN)
  LANGUAGE plpgsql
  AS $$
  BEGIN
    -- Check whether user_id valid
    SELECT COUNT(*) > 0 INTO exist
    FROM users
    WHERE id = user_id;


    IF exist THEN
      -- Get user_email, uname and verified from users table
      SELECT email INTO user_email
      FROM users
      WHERE id = user_id;


      SELECT username INTO uname
      FROM users
      WHERE id = user_id;


      SELECT is_verified INTO verified
      FROM users
      WHERE id = user_id;


      IF NOT verified THEN
        -- Insert secret code to email_verifications
        INSERT INTO email_verifications (id, user_id, email, secret_code)
        VALUES ((SELECT COUNT(*) FROM email_verifications) + 1, user_id, user_email, secret_code)
        RETURNING id INTO email_id;
      END IF;
    END IF;


    -- Change user_email, uname, verified, email_id if NULL
    user_email = COALESCE(user_email, '');
    uname = COALESCE(uname, '');
    verified = COALESCE(verified, false);
    email_id = COALESCE(email_id, -1);
  END;
  $$;
COMMIT;


-- Create procedure verify_email
BEGIN TRANSACTION;
  DROP PROCEDURE IF EXISTS verify_email;


  CREATE OR REPLACE PROCEDURE verify_email(email_id INT, requested_secret_code CHAR(32), OUT is_verified BOOLEAN, OUT is_expired BOOLEAN, OUT is_match BOOLEAN)
  LANGUAGE plpgsql
  AS $$
  BEGIN
    SELECT e.secret_code = requested_secret_code INTO is_match
    FROM email_verifications e
    WHERE e.id = email_id;


    SELECT e.expired_at :: time < CURRENT_TIME INTO is_expired
    FROM email_verifications e
    WHERE e.id = email_id;


    SELECT u.is_verified INTO is_verified
    FROM email_verifications e, users u
    WHERE e.user_id = u.id AND e.id = email_id;


    IF is_match AND NOT is_expired THEN
      UPDATE users
      SET is_verified = TRUE
      WHERE id = (SELECT user_id FROM email_verifications WHERE id = email_id);
    END IF;
  END;
  $$;
COMMIT;

