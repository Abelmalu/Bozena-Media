-- create table likes --

CREATE TABLE likes(
    id SERIAL NOT NULL PRIMARY KEY,
    post_id INT NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT fk_posts FOREIGN KEY(post_id) REFERENCES posts_cache(post_id),
    CONSTRAINT fk_users FOREIGN KEY(user_id) REFERENCES users_cache(user_id),
    CONSTRAINT  unique_user_post_like UNIQUE (user_id,post_id)
)