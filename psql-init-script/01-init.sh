#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE post_db;
    CREATE DATABASE follow_db;
    CREATE DATABASE like_db;
    CREATE DATABASE feed_db;
    CREATE DATABASE notification_db;
    CREATE DATABASE chat_db;

EOSQL
