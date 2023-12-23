package middleware

import (
	"context"
	"fmt"
	"gocommerce/constants"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// getUserIDFromToken decodes the JWT token and extracts the user ID.
func getUserIDFromToken(tokenString string) (int, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecretKey, nil
	})

	if err != nil {
		return 0, fmt.Errorf("error parsing token: %w", err)
	}

	// Check if token is valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token or cannot convert claims")
	}

	// Extract user ID from token claims
	userID, ok := claims["user_id"].(float64) // JWT stores numbers as float64
	if !ok {
		return 0, fmt.Errorf("error extracting user ID from token, type assertion failed: user_id is type %T", claims["user_id"])
	}

	return int(userID), nil
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request path
		log.Printf("Processing request for %s", r.URL.Path)

		// Skip middleware for certain routes
		if r.URL.Path == "/api/auth/register" || r.URL.Path == "/api/auth/login" {
			next.ServeHTTP(w, r)
			return
		}

		// Try to get the token from the Authorization header
		tokenString := getTokenFromHeader(r)

		// Log if the token was found in the header
		if tokenString != "" {
			log.Println("Token found in header")
		} else {
			log.Println("Token not found in header, checking cookie")
			// If not found in the header, try to get it from a cookie
			tokenString = getTokenFromCookie(r)

			if tokenString != "" {
				log.Println("Token found in cookie")
			} else {
				log.Println("Token not found in cookie")
			}
		}

		// Decode token and get user ID
		userID, err := getUserIDFromToken(tokenString)
		if err != nil {
			// Log the error
			log.Printf("Error decoding token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Log user ID
		log.Printf("Authenticated user ID: %v", userID)

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
