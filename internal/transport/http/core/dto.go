//go:generate easyjson -all dto.go

package core

type ErrorResponse struct {
	Code      int    `json:"code"`
	ErrorText string `json:"error"`
}
