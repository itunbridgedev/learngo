package middleware

import (
	"context"
	"errors"
	"gocommerce/constants"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecretKey = os.Getenv("JWT_SECRET_KEY")

// getUserIDFromToken decodes the JWT token and extracts the user ID.
func getUserIDFromToken(tokenString string) (int, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method in token")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return 0, err
	}

	// Check if token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user ID from token claims
		userID, ok := claims["user_id"].(float64) // Make sure the type assertion is correct; it's float64 by default
		if !ok {
			return 0, errors.New("error extracting user ID from token")
		}
		return int(userID), nil
	}

	return 0, errors.New("invalid token")
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip middleware for certain routes
		if r.URL.Path == "/api/auth/register" || r.URL.Path == "/api/auth/login" {
			next.ServeHTTP(w, r)
			return
		}

		// Try to get the token from the Authorization header
		tokenString := getTokenFromHeader(r)

		// If not found in the header, try to get it from a cookie
		if tokenString == "" {
			tokenString = getTokenFromCookie(r)
		}
		// Decode token and get user ID
		userID, err := getUserIDFromToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user ID to context and proceed
		ctx := context.WithValue(r.Context(), constants.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	// Format: "Bearer <token>"
	if strings.HasPrefix(bearerToken, "Bearer ") {
		return strings.TrimPrefix(bearerToken, "Bearer ")
	}
	return ""
}

// getTokenFromCookie extracts the token from a specific cookie
func getTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("token_cookie_name") // Replace with your token cookie's name
	if err != nil {
		return ""
	}
	return cookie.Value
}
