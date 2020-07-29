package app
 import (
	 "github.com/dgrijalva/jwt-go"
 )

//COMMENTS
type Event struct {
	Pag         string `json:"Pag"` //
	ID          string `json:"ID"`
	Page        string `json:"Page"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

// NewEvent constructs a event object
func NewEvent() *Event {
	return &Event{}
}

type AllEvents []Event

func GetEvents() *AllEvents {
	return &AllEvents{}
}

///Content
type ContentPage struct {
	Page    string `json:"Page"`
	Title   string `json:"Title"`
	Content string `json:"Content"`
}

func NewContentPage() *ContentPage {
	return &ContentPage{}
}

//TITLES
type Title struct {
	Title string `json:"Title"`
}

func NewTitle() *Title {
	return &Title{}
}

type AllTitles []Title

func GetTitles() *AllTitles {
	return &AllTitles{}
}

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

func NewToken() *Token {
	return &Token{}
}

//a struct to rep user account
type Account struct {
	ID		 uint    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

func NewAccount() *Account {
	return &Account{}
}

