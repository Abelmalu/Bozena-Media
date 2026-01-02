-- posts table --

CREATE TABLE posts (
    id SERIAL,
    CONSTRAINT pk_users PRIMARY KEY (id)
    title VARCHAR(100) NOT NULL,
    content TEXT,
    user_id INT NOT NULL,
    CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
      REFERENCES users(id)
      ON DELETE CASCADE
);
