-- posts table --

CREATE TABLE posts (
    id SERIAL,
   
    title VARCHAR(100) NOT NULL,
    content TEXT,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
     CONSTRAINT pk_users PRIMARY KEY (id),
      CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
      REFERENCES users_cache(user_id)
      ON DELETE CASCADE
);

CREATE INDEX idx_user_id ON posts (user_id);

