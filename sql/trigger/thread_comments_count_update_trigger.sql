-- Drop the thread_comments_count_update_trigger trigger and update_thread_comments_count function if exists
DROP TRIGGER IF EXISTS thread_comments_count_update_trigger ON comments;
DROP FUNCTION IF EXISTS update_thread_comments_count;

-- Create the update_thread_comments_count function
CREATE OR REPLACE FUNCTION update_thread_comments_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE threads SET comments_count = comments_count + 1 WHERE NEW.threadid = id AND NOT NEW.parentid IS NULL;
    ELSIF TG_OP = 'DELETE' THEN
		UPDATE threads SET comments_count = comments_count - 1 WHERE OLD.threadid = id AND NEW.parentid IS NOT NULL AND NOT OLD.is_deleted;
    ELSIF TG_OP = 'UPDATE' AND OLD.is_deleted <> NEW.is_deleted THEN
		IF OLD.is_deleted THEN
			UPDATE threads SET comments_count = comments_count + 1 WHERE OLD.threadid = id AND NEW.parentid IS NOT NULL;
		ELSE
			UPDATE threads SET comments_count = comments_count - 1 WHERE OLD.threadid = id AND NEW.parentid IS NOT NULL;
		END IF;
	END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create the thread_comments_count_update_trigger trigger
CREATE TRIGGER thread_comments_count_update_trigger
AFTER INSERT OR UPDATE OR DELETE ON comments
FOR EACH ROW
EXECUTE FUNCTION update_thread_comments_count();
