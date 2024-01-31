-- Procedure for storing the verification code in email_verifications table

-- Drop the make_verification procedure if exists
DROP PROCEDURE IF EXISTS make_verification;

-- Create the make_verification procedure
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
