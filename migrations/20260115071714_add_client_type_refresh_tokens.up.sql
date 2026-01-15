ALTER TABLE refresh_tokens
ADD COLUMN client_type VARCHAR(10) NOT NULL
CONSTRAINT refresh_tokens_client_type_check
CHECK (client_type IN ('web', 'mobile'));