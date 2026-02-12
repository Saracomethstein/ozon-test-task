package models

type AddCommentInput struct {
	PostID   string  `db:"postId"`
	ParentID *string `db:"parentId,omitempty"`
	Author   string  `db:"author"`
	Text     string  `db:"text"`
}

type Comment struct {
	ID        int64              `db:"id"`
	PostID    int64              `db:"postId"`
	ParentID  *int64             `db:"parentId,omitempty"`
	Author    string             `db:"author"`
	Text      string             `db:"text"`
	Path      string             `db:"path"`
	CreatedAt string             `db:"createdAt"`
	Children  *CommentConnection `db:"children"`
}

type CommentConnection struct {
	Edges      []*CommentEdge `db:"edges"`
	PageInfo   *PageInfo      `db:"pageInfo"`
	TotalCount int32          `db:"totalCount"`
}

type CommentEdge struct {
	Cursor string   `db:"cursor"`
	Node   *Comment `db:"node"`
}

type CreatePostInput struct {
	Title         string `db:"title"`
	Body          string `db:"body"`
	Author        string `db:"author"`
	AllowComments *bool  `db:"allowComments,omitempty"`
}

type PageInfo struct {
	EndCursor   *string `db:"endCursor,omitempty"`
	HasNextPage bool    `db:"hasNextPage"`
}

type Post struct {
	ID            string             `db:"id"`
	Title         string             `db:"title"`
	Body          string             `db:"body"`
	Author        string             `db:"author"`
	AllowComments bool               `db:"allowComments"`
	CreatedAt     string             `db:"createdAt"`
	Comments      *CommentConnection `db:"comments"`
}

type PostConnection struct {
	Edges      []*PostEdge `db:"edges"`
	PageInfo   *PageInfo   `db:"pageInfo"`
	TotalCount int32       `db:"totalCount"`
}

type PostEdge struct {
	Cursor string `db:"cursor"`
	Node   *Post  `db:"node"`
}
