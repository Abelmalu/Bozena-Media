

ALTER TABLE posts_cache 
ADD COLUMN users_id INT NULL,
ADD CONSTRAINT fk_posts_users 
    FOREIGN KEY (users_id) 
    REFERENCES users_cache(user_id);
