package limit

import (
	"net/http"
	"strconv"
)

const (
	DefaultLimit = 50
	MaxLimit     = 500
	MinLimit     = 1
)

func GetFromRequest(r *http.Request) int {
	var limitInt int
	if r.Method == http.MethodGet {
		q := r.URL.Query()
		limitStr := q.Get("limit")
		limitInt, _ = strconv.Atoi(limitStr)
	}

	limitInt = CheckLimit(limitInt)

	return limitInt
}

func CheckLimit(l int) int {
	if l > MaxLimit {
		l = MaxLimit
	}
	if l < MinLimit {
		l = DefaultLimit
	}
	return l
}
