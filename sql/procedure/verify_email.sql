-- Procedure for verifying email

-- Drop the verify_email procedure if exists
DROP PROCEDURE IF EXISTS verify_email;

-- Create the verify_email procedure
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
