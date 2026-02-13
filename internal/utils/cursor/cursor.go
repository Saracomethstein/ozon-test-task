package cursor

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	cursorSeparator = "|"
)

func Encode(createdAt string, id int64) string {
	encoded := createdAt + cursorSeparator + strconv.FormatInt(id, 10)
	return base64.RawURLEncoding.EncodeToString([]byte(encoded))
}

func Decode(cursor string) (createdAt string, id int64, err error) {
	decoded, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return "", 0, errors.New("invalid cursor encoding")
	}

	parts := strings.SplitN(string(decoded), cursorSeparator, 2)
	if len(parts) != 2 {
		return "", 0, errors.New("invalid cursor format")
	}

	id, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, errors.New("invalid cursor id")
	}

	return parts[0], id, nil
}
