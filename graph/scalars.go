package graph

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalTime converts time.Time to a GraphQL string (RFC3339Nano)
func MarshalTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconvQuote(t.Format(time.RFC3339Nano)))
	})
}

// UnmarshalTime parses an incoming GraphQL value to time.Time
func UnmarshalTime(v any) (time.Time, error) {
	s, ok := v.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("time must be a string")
	}
	return time.Parse(time.RFC3339Nano, s)
}

// small helper to avoid importing strconv in the marshaler body repeatedly
func strconvQuote(s string) string {
	// wrap in quotes for valid JSON string
	return "\"" + s + "\""
}
