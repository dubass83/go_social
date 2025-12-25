-- Drop the existing constraint
ALTER TABLE comments
DROP CONSTRAINT comments_post_id_fkey;

-- Add it back with CASCADE
ALTER TABLE comments
ADD CONSTRAINT comments_post_id_fkey
FOREIGN KEY (post_id)
REFERENCES posts(id)
ON DELETE CASCADE;
