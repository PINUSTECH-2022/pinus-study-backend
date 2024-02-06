-- Procedure for storing the password recovery code in password_recoveries table

-- Drop the make_password_recovery procedure if exists
DROP PROCEDURE IF EXISTS make_password_recovery;

-- Create the make_password_recovery procedure
CREATE OR REPLACE PROCEDURE make_password_recovery(user_id INT, secret_code VARCHAR,
                        OUT exist BOOLEAN, OUT verified BOOLEAN, OUT recovery_id INT, OUT user_email VARCHAR)
LANGUAGE plpgsql
AS $$
BEGIN
    -- Check whether user_id valid
    SELECT COUNT(*) > 0 INTO exist
    FROM users
    WHERE id = user_id;
    
    IF exist THEN
        -- Get verified and email from users table
        SELECT is_verified INTO verified
        FROM users
        WHERE id = user_id;
        
		SELECT email INTO user_email
		FROM users
		WHERE id = user_id;
		
        IF verified THEN
            -- Insert secret code to password_recovery
            INSERT INTO password_recoveries(id, user_id, secret_code)
            VALUES ((SELECT COUNT(*) FROM password_recoveries) + 1, user_id, secret_code)
			RETURNING id INTO recovery_id;
        END IF;
    END IF;
    
    -- Change verified, email and recovery_id if NULL
    verified = COALESCE(verified, false);
	user_email = COALESCE(user_email, '');
	recovery_id = COALESCE(recovery_id, 0);
END;
$$;
