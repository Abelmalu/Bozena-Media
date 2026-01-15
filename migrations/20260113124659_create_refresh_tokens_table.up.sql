
CREATE TABLE refresh_tokens(
    id SERIAL,
    user_id INT,
    token_text TEXT UNIQUE,
    expires_at TIMESTAMP,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
     
-- constraints 
    CONSTRAINT fk_users FOREIGN KEY(user_id) REFERENCES users(id)




);