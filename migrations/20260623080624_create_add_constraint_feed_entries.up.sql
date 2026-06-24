ALTER TABLE feed_entries

ADD CONSTRAINT fk_user_id  
FOREIGN KEY (user_id)
REFERENCES users_cache (user_id),

ADD CONSTRAINT fk_owner_id
FOREIGN KEY (owner_id)
REFERENCES users_cache (user_id),

ADD CONSTRAINT fk_post_id 
FOREIGN KEY (post_id)
REFERENCES posts_cache (post_id);