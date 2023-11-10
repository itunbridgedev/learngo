package models

import "time"

// TokenDetails holds the information about a user's authentication token
type TokenDetails struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"` // Optional, based on your auth strategy
	AtExpires    time.Time `json:"at_expires"`              // Access token expiration time
	RtExpires    time.Time `json:"rt_expires,omitempty"`    // Refresh token expiration time
	UserID       int       `json:"user_id"`                 // User ID associated with the token
	// Add any other relevant fields
}
