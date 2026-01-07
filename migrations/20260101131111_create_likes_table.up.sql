-- create table likes --

CREATE TABLE likes(
    id SERIAL NOT NULL PRIMARY KEY,
    post_id INT NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT fk_posts FOREIGN KEY(post_id) REFERENCES posts(id),
    CONSTRAINT fk_users FOREIGN KEY(user_id) REFERENCES users(id)
)