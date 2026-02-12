package models

type AddCommentInput struct {
	PostID   string
	ParentID *string
	Author   string
	Text     string
}

type Comment struct {
	ID        int64
	PostID    int64
	ParentID  *int64
	Author    string
	Text      string
	CreatedAt string
}

type CommentConnection struct {
	Edges      []*CommentEdge
	PageInfo   *PageInfo
	TotalCount int32
}

type CommentEdge struct {
	Cursor string
	Node   *Comment
}

type CreatePostInput struct {
	Title         string
	Body          string
	Author        string
	AllowComments *bool
}

type PageInfo struct {
	EndCursor   *string
	HasNextPage bool
}

type Post struct {
	ID            string
	Title         string
	Body          string
	Author        string
	AllowComments bool
	CreatedAt     string
	Comments      *CommentConnection
}

type PostConnection struct {
	Edges      []*PostEdge
	PageInfo   *PageInfo
	TotalCount int32
}

type PostEdge struct {
	Cursor string
	Node   *Post
}
