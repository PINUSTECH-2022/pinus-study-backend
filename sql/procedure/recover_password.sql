-- Procedure for password recovery

-- Drop the recover_password procedure if exists
DROP PROCEDURE IF EXISTS recover_password;

-- Create the recover_password procedure
CREATE OR REPLACE PROCEDURE recover_password(recovery_id INT, requested_secret_code CHAR(32), new_password VARCHAR, new_salt VARCHAR, 
											 OUT is_exist BOOLEAN, OUT is_expired BOOLEAN, OUT is_match BOOLEAN, OUT used BOOLEAN)
LANGUAGE plpgsql
AS $$
BEGIN
	-- Get is_exist, is_expired, is_match, and is_used
	SELECT COUNT(*) > 0 INTO is_exist
	FROM password_recoveries
	WHERE id = recovery_id;
	
	SELECT secret_code = requested_secret_code INTO is_match
	FROM password_recoveries
	WHERE id = recovery_id;
	
	SELECT expired_at :: time < CURRENT_TIME INTO is_expired
	FROM password_recoveries
	WHERE id = recovery_id;
	
	SELECT is_used INTO used
	FROM password_recoveries
	WHERE id = recovery_id;
	
	IF is_exist AND NOT is_expired AND is_match AND NOT used THEN
		-- Update the password_recoveries link to become used
		UPDATE password_recoveries
		SET is_used = TRUE
		WHERE id = recovery_id;
	
		-- Update the user's password and salt
		UPDATE users
		SET password = new_password, salt = new_salt
		WHERE id = (SELECT user_id FROM password_recoveries WHERE id = recovery_id);
	END IF;
	
	-- Change is_exist, is_expired, is_match, and used if NULL
	is_exist = COALESCE(is_exist, FALSE);
	is_expired = COALESCE(is_expired, FALSE);
	is_match = COALESCE(is_match, FALSE);
	used = COALESCE(used, FALSE);
END;
$$;
