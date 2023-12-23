package models

type Credentials struct {
	UserId   float64 `json:user_id`
	Username string  `json:"username"` // or Email if you use email for login
	Password string  `json:"password"`
}
