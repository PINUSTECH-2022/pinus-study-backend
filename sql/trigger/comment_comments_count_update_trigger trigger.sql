-- Drop the comment_comments_count_update_trigger trigger and update_comment_comments_count function if exists
DROP TRIGGER IF EXISTS comment_comments_count_update_trigger ON comments;
DROP FUNCTION IF EXISTS update_comment_comments_count;

-- Create the update_comment_comments_count function
CREATE OR REPLACE FUNCTION update_comment_comments_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE comments SET comments_count = comments_count + 1 WHERE NEW.parentid = id;
    ELSIF TG_OP = 'DELETE' THEN
		UPDATE comments SET comments_count = comments_count - 1 WHERE OLD.parentid = id AND NOT OLD.is_deleted;
    ELSIF TG_OP = 'UPDATE' AND OLD.is_deleted <> NEW.is_deleted THEN
		IF OLD.is_deleted THEN
			UPDATE comments SET comments_count = comments_count + 1 WHERE OLD.parentid = id;
		ELSE
			UPDATE comments SET comments_count = comments_count - 1 WHERE OLD.parentid = id;
		END IF;
	END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create the comment_comments_count_update_trigger trigger
CREATE TRIGGER comment_comments_count_update_trigger
AFTER INSERT OR UPDATE OR DELETE ON comments
FOR EACH ROW
EXECUTE FUNCTION update_comment_comments_count();
