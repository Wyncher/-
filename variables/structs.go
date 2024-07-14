package variables

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// struct of newMessageCounter
type MessagesCounter struct {
	MCounter int
	UserID   int
}

// struct of Message
type Message struct {
	Text      string
	File      []byte
	FileName  string
	Date      time.Time
	DateStr   string
	Pos       string
	MessageID int
	Favourite bool
}

// struct of favourite Message
type FavouriteMessages struct {
	Text      string
	File      []byte
	FileName  string
	Date      time.Time
	DateStr   string
	Pos       string
	FromID    int
	MessageID int
	UserName  string
	Favourite bool
}

// struct of user connects
type User struct {
	Id       int
	Username string
	Logo     string
}

// structure of JWT claim(contains username)
type JwtCustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}
