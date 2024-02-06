-- Procedure for updating / fixing all counts in threads table

-- Drop the update_all_thread_count procedure if exists
DROP PROCEDURE IF EXISTS update_all_comment_count;

-- Create the update_all_thread_count procedure
CREATE OR REPLACE PROCEDURE update_all_comment_count()
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE Comments t
    SET likes = (SELECT COUNT(*) FROM Likes_Comments l WHERE l.threadid = t.id AND l.state),
    dislikes = (SELECT COUNT(*) FROM Likes_Comments l WHERE l.threadid = t.id AND NOT l.state),
END;
$$;
