
CREATE TABLE follows (
    id BIGSERIAL PRIMARY KEY,

    follower_id BIGINT NOT NULL,
    following_id BIGINT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_follower
        FOREIGN KEY (follower_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_following
        FOREIGN KEY (following_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT unique_follow
        UNIQUE (follower_id, following_id),

    CONSTRAINT no_self_follow
        CHECK (follower_id <> following_id)
);