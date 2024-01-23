-- Procedure for updating / fixing all counts in threads table

-- Drop the update_all_thread_count procedure if exists
DROP PROCEDURE IF EXISTS update_all_thread_count;

-- Create the update_all_thread_count procedure
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
