CREATE TABLE notifications(

    id SERIAL PRIMARY KEY,
    recipient_id INT NOT NULL,
    actor_id    INT NOT NULL,
    is_read     BOOLEAN NOT NULL DEFAULT FALSE,
    message   TEXT DEFAULT 'started following you' NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    read_at TIMESTAMP NULL,


    CONSTRAINT fk_actor_id 
    FOREIGN KEY(actor_id)
    REFERENCES users_cache(user_id),


    CONSTRAINT fk_recipient_id
    FOREIGN KEY(recipient_id)
    REFERENCES users_cache(user_id)



);

CREATE INDEX idx_recipent_id ON notifications(recipient_id);
CREATE INDEX idx_created_at ON notifications(created_at DESC);