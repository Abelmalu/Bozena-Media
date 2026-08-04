package dto

type MessageRequest struct {
	Message    string `json:"message"`
	ReceiverID int    `json:"receiver_id"`
}

