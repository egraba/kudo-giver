package main

type Kudo struct {
	ID         int64  `json:"id"`
	SenderID   int64  `json:"senderId"`
	ReceiverID int64  `json:"receiverId"`
	Message    string `json:"message"`
}
