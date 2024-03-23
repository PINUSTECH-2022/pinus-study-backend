-- Procedure for updating / fixing all counts in threads table

-- Drop the update_all_thread_count procedure if exists
DROP PROCEDURE IF EXISTS update_all_comment_count;

-- Create the update_all_thread_count procedure
CREATE OR REPLACE PROCEDURE update_all_comment_count()
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE comments c
    SET likes_count = (SELECT COUNT(*) FROM likes_comments l WHERE l.commentid = c.id AND l.state),
    dislikes_count = (SELECT COUNT(*) FROM likes_comments l WHERE l.commentid = c.id AND NOT l.state),
	comments_count = (SELECT COUNT(*) FROM comments c1 WHERE c1.parentid = c.id AND NOT c1.is_deleted);
END;
$$;
