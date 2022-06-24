package main

type Person struct {
	ID        int32  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
