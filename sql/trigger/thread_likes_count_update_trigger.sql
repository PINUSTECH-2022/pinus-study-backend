-- Drop the thread_likes_count_update_trigger trigger and update_thread_likes_count function if exists
DROP TRIGGER IF EXISTS thread_likes_count_update_trigger ON likes_threads;
DROP FUNCTION IF EXISTS update_thread_likes_count;

-- Create the update_thread_likes_count function
CREATE OR REPLACE FUNCTION update_thread_likes_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE threads SET likes_count = likes_count + 1 WHERE NEW.threadid = id AND NEW.state;
        UPDATE threads SET dislikes_count = dislikes_count + 1 WHERE NEW.threadid = id AND NOT NEW.state;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE threads SET likes_count = likes_count - 1 WHERE OLD.threadid = id AND OLD.state;
        UPDATE threads SET dislikes_count = dislikes_count - 1 WHERE OLD.threadid = id AND NOT OLD.state;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create the thread_likes_count_update_trigger trigger
CREATE TRIGGER thread_likes_count_update_trigger
AFTER INSERT OR DELETE ON likes_threads
FOR EACH ROW
EXECUTE FUNCTION update_thread_likes_count();
