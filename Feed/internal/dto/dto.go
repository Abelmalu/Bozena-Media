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

}

type PaginatedResponse struct {
	UserFeeds  []*UserFeed
	Cursor string
    HasNext bool
}
