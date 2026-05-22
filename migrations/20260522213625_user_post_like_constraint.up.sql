ALTER TABLE likes
ADD CONSTRAINT unique_user_post_like UNIQUE (user_id, post_id);
