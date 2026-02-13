package repository

type Container struct {
	Post    PostUC
	Comment CommentUC
}

func New(
	postRepo PostUC,
	commentRepo CommentUC,
) *Container {
	return &Container{
		Post:    postRepo,
		Comment: commentRepo,
	}
}
