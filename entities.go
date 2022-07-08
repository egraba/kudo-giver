package main

type Person struct {
	ID        int32  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	NbKudos   int32  `json:"nbKudos"`
}

type Kudo struct {
	ID         int32  `json:"id"`
	SenderID   int32  `json:"senderId"`
	ReceiverID int32  `json:"receiverId"`
	Message    string `json:"message"`
}
