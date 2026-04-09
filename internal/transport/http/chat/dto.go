package chat

//go:generate easyjson -all dto.go

type SendMessageRequest struct {
	Text string `json:"text"`
}

type SendMessageResponse struct {
	MessageID string `json:"message_id"`
}
