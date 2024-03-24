DROP TRIGGER IF EXISTS review_likes_count_update_trigger ON likes_threads;
DROP FUNCTION IF EXISTS update_review_likes_count;

CREATE OR REPLACE FUNCTION update_review_likes_count()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
  UPDATE reviews SET likes_count = likes_count + 1 WHERE NEW.reviewid = id AND NEW.state;
  UPDATE reviews SET dislikes_count = dislikes_count + 1 WHERE NEW.reviewid = id AND NOT NEW.state;
  ELSIF TG_OP = 'DELETE' THEN
  UPDATE reviews SET likes_count = likes_count - 1 WHERE OLD.reviewid = id AND OLD.state;
  UPDATE reviews SET dislikes_count = dislikes_count - 1 WHERE OLD.reviewid = id AND NOT OLD.state;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER review_likes_count_update_trigger
AFTER INSERT OR DELETE ON likes_threads
FOR EACH ROW
EXECUTE FUNCTION update_review_likes_count();