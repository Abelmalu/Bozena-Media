-- posts table --

CREATE TABLE posts (
    id SERIAL,
   
    title VARCHAR(100) NOT NULL,
    content TEXT,
    user_id INT NOT NULL,
     CONSTRAINT pk_users PRIMARY KEY (id),
    CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
      REFERENCES users(id)
      ON DELETE CASCADE
);
