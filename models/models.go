package models

type Message struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type JoinClientResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type LeaveClientResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type SendMessageResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type GetMessagesResponse struct {
	Id               string     `json:"id"`
	Messages         []*Message `json:"messages"`
	MessageIndicator string     `json:"messageindicator,omitempty"`
}

type PingResponse struct {
	Message string `json:"message"`
}
