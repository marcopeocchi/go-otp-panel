package dto

type MessagePublishRequest struct {
	Message   string `json:"message"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}
