package middleware

import (
	"context"
	"net/http"
	"time"

	myLoader "github.com/Saracomethstein/ozon-test-task/internal/graphql/dataloader"
	"github.com/graph-gophers/dataloader"
)

func DataloaderMiddleware(commentLoader myLoader.CommentLoader) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loader := dataloader.NewBatchedLoader(
				commentLoader.BatchGetChildren,
				dataloader.WithWait(2*time.Millisecond),
				dataloader.WithBatchCapacity(50),
			)

			ctx := context.WithValue(r.Context(), myLoader.Key, loader)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
