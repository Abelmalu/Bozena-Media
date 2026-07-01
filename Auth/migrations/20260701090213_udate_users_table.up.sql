

UPDATE users 
SET follower_count = 0, 
    following_count = 0 
WHERE follower_count IS NULL 
   OR following_count IS NULL;
