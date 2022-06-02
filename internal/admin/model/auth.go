package model

// AuthCredentials represent login credentials.
type AuthCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
