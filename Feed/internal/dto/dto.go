package dto

type PostOwner struct {
	ID       int64

}

type UserFeed struct {

	ID int
	PostID      int64
	PostTitle   string
	PostContent string
	PostOwnerID int64
	UserName string
	Name     string
	Image *string
	LikeCount int

}

type PaginatedResponse struct {
	UserFeeds  []*UserFeed
	Cursor string
    HasNext bool
}

type FeedEntry struct {

	UserID int `db:"user_id"`
	OwnerID int `db:"owner_id"`
	PostID int `db:"post_id"`
}


type PostCache struct {

	UserID int `db:"user_id"`
	PostID int `db:"post_id"`

}


type UserCachePostsResponse struct {

	CachePosts []*PostCache


}
