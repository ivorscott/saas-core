package mauth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Token struct {
	ID          string    `db:"ma_token_id" json:"id"`
	AccessToken string    `db:"token" json:"access_token"`
	Created     time.Time `db:"created" json:"created"`
}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

type NewInvite struct {
	Emails []string `json:"emailList"`
}

type InvitedUser struct {
	Email  string `json:"email"`
	UserId string `json:"user_id"`
}

const DatabaseConnection = "Username-Password-Authentication"
