package repository

import "database/sql"



type FeedRepository struct {

	DB *sql.DB
}





func NewFeedRepository(db *sql.DB) *FeedRepository {



	return &FeedRepository{
		DB: db,
	}
}