package cursor

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

type Cursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int64     `json:"id"`
}

func (c *Cursor) EncodeCursor() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func DecodeCursor(s string) (*Cursor, error) {
	if s == "" {
		return nil, nil
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	var c Cursor
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func GetFromRequest(r *http.Request) string {
	var cStr string
	if r.Method == http.MethodGet {
		q := r.URL.Query()
		cStr = q.Get("cursor")
	}
	return cStr
}
