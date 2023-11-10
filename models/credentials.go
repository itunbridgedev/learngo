package models

type Credentials struct {
	Username string `json:"username"` // or Email if you use email for login
	Password string `json:"password"`
}
