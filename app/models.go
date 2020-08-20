package app

import (
	"github.com/dgrijalva/jwt-go"
)

//COMMENTS
type Comment struct {
	Page    string `json:"Page"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
}

// NewComment constructs a Comment object
func NewComment() *Comment {
	return &Comment{}
}

type AllComments []Comment

func GetComments() *AllComments {
	return &AllComments{}
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
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

func NewAccount() *Account {
	return &Account{}
}
