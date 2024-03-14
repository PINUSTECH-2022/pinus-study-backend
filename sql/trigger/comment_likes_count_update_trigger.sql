DROP TRIGGER IF EXISTS comment_likes_count_update_trigger ON Likes_Comments;
DROP FUNCTION IF EXISTS update_comment_likes_count;

-- Create the update_thread_likes_count function
CREATE OR REPLACE FUNCTION update_comment_likes_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE Comments SET likes = likes + 1 WHERE NEW.threadid = id AND NEW.state;
        UPDATE Comments SET dislikes = dislikes + 1 WHERE NEW.threadid = id AND NOT NEW.state;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE Comments SET likes = likes - 1 WHERE OLD.threadid = id AND OLD.state;
        UPDATE Comments SET dislikes = dislikes - 1 WHERE OLD.threadid = id AND NOT OLD.state;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create the thread_likes_count_update_trigger trigger
CREATE TRIGGER comment_likes_count_update_trigger
AFTER INSERT OR DELETE ON Likes_Comments
FOR EACH ROW
EXECUTE FUNCTION update_comment_likes_count();