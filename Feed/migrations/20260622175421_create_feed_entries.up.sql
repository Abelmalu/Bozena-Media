CREATE TABLE feed_entries(

    id SERIAL PRIMARY KEY,
    user_id  INT,
    post_id INT NOT NULL,
    owner_id INT NOT NULL
);