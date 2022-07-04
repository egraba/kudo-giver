package main

type Kudo struct {
	ID         int32  `json:"id"`
	SenderID   int32  `json:"senderId"`
	ReceiverID int32  `json:"receiverId"`
	Message    string `json:"message"`
}
