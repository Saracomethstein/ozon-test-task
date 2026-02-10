package comment

type commentService struct{}

func New() UseCase {
	return &commentService{}
}
