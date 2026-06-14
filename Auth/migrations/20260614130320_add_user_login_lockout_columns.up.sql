-- ADD user lockout columns 

ALTER TABLE users 
ADD COLUMN failed_attempts INT DEFAULT 0,
ADD COLUMN  is_permanently_locked BOOLEAN DEFAULT FALSE,
ADD COLUMN  temporary_locked_until TIMESTAMP;