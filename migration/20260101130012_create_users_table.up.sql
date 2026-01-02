
CREATE TABLE users(
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(100),
    username VARCHAR(100) UNIQUE,
    `password` VARCHAR(255),
    `email` VARCHAR(100) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()

);