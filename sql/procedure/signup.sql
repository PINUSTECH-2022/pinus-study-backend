-- Procedure for user to signup

-- Drop the signup procedure if exists
DROP PROCEDURE IF EXISTS signup;

-- Create the signup procedure
CREATE OR REPLACE PROCEDURE signup(candidate_username VARCHAR, candidate_email VARCHAR, password VARCHAR, salt VARCHAR, secret_code CHAR(32),
								   OUT user_id INT, OUT email_id INT, OUT is_email_exist BOOLEAN, OUT is_username_exist BOOLEAN, OUT is_verified BOOLEAN)
LANGUAGE plpgsql
AS $$
BEGIN
	-- Check if email already exist
	SELECT COUNT(*) > 0 INTO is_email_exist FROM users WHERE email = candidate_email;
	
	-- Check if username already exist
	SELECT COUNT(*) > 0 INTO is_username_exist FROM users WHERE username = candidate_username AND email <> candidate_email;

    -- Check if email already verified
    SELECT COUNT(*) > 0 INTO is_verified FROM users u WHERE u.email = candidate_email AND u.is_verified = true;
    
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

	IF is_email_exist AND NOT is_verified THEN
		SELECT id INTO user_id FROM users WHERE email = candidate_email;
	END IF;
END;
$$;