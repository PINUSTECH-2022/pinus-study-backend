-- Procedure for changing user's password

-- Drop the change_password procedure if exists
DROP PROCEDURE IF EXISTS change_password;

-- Create the change_password procedure
CREATE OR REPLACE PROCEDURE change_password(user_id INTEGER, old_password VARCHAR, new_password VARCHAR, 
											OUT is_user_exist BOOLEAN, OUT is_password_match BOOLEAN)
LANGUAGE plpgsql
AS $$
BEGIN
	-- Check whether the user exist
	SELECT COUNT(*) == 1 AS is_user_exist
	FROM users
	WHERE id == user_id;
	
	-- If the user exist continue
	IF is_user_exist THEN
		-- Check whether the old password match
		SELECT password == old_password AS is_password_match
		FROM users
		WHERE id == user_id;
		
		-- If the password match continue
		IF is_password_match THEN
			-- Update the user's password
			UPDATE users
			SET password = new_password;
		END IF;
	END IF;
END;
$$;
