package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gocommerce/models"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationHandler struct {
	DB *sql.DB
}

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func (h *AuthenticationHandler) validateCredentials(username, password string) (float64, error) {
	var hashedPassword string
	var userID float64
	err := h.DB.QueryRow("SELECT id, password_hash FROM users WHERE username = $1", username).Scan(&userID, &hashedPassword)
	if err != nil {
		// Handle sql.ErrNoRows separately if needed (user not found)
		// Return 0 for the userID and the error
		return 0, err
	}

	// Compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Return 0 for the userID and a specific error indicating password mismatch
		return 0, fmt.Errorf("password mismatch")
	}

	// Return the user ID and nil error if the credentials are valid
	return userID, nil
}

func (h *AuthenticationHandler) GenerateToken(userID float64) (string, string, error) {
	// Log the user ID being used
	log.Printf("Generating token for user ID: %v", userID)

	// Generate Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(), // short-lived access token
	})

	// Log the Access Token Claims
	log.Printf("Access Token Claims: %+v", accessToken.Claims)

	at, err := accessToken.SignedString(jwtSecretKey)
	if err != nil {
		log.Printf("Error signing access token: %v", err)
		return "", "", err
	}

	// Generate Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // long-lived refresh token
	})

	// Log the Refresh Token Claims
	log.Printf("Refresh Token Claims: %+v", refreshToken.Claims)

	rt, err := refreshToken.SignedString(jwtSecretKey)
	if err != nil {
		log.Printf("Error signing refresh token: %v", err)
		return "", "", err
	}

	return at, rt, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// validateUserData checks the new user data for validity.
func (h *AuthenticationHandler) validateUserData(user *models.User) error {
	// Validate email format
	if err := validateEmailFormat(user.Email); err != nil {
		return err
	}

	// Check if username or email is already in use
	if err := h.checkUserExists(user.Username, user.Email); err != nil {
		return err
	}

	// Check password complexity
	if err := validatePasswordComplexity(user.Password); err != nil {
		return err
	}

	return nil
}

// validateEmailFormat validates the format of the email.
func validateEmailFormat(email string) error {
	// Regular expression for validating an email
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// checkUserExists checks if the username or email is already in use.
func (h *AuthenticationHandler) checkUserExists(username, email string) error {
	var id int
	err := h.DB.QueryRow("SELECT id FROM users WHERE username = $1 OR email = $2", username, email).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return errors.New("username or email already in use")
	}
	return nil
}

// validatePasswordComplexity checks if the password meets complexity requirements.
func validatePasswordComplexity(password string) error {
	// Example: Minimum 8 characters
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

func (h *AuthenticationHandler) createUserRecord(user *models.User) error {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword

	// Insert user record into the database
	sqlStatement := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`
	_, err = h.DB.Exec(sqlStatement, user.Username, user.Email, user.PasswordHash)
	return err
}

func (h *AuthenticationHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate credentials
	userID, err := h.validateCredentials(creds.Username, creds.Password)
	if err != nil {
		if err.Error() == "password mismatch" {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			// Handle other errors, like database errors
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error validating credentials: %v", err)
		}
		return
	}
	// Check if userID is valid
	if userID == 0 {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate token
	accessToken, refreshToken, err := h.GenerateToken(userID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		log.Printf("Error generating token: %v", err)
		return
	}

	response := map[string]string{
		"token":        accessToken,
		"refreshToken": refreshToken,
	}

	// Encoding the response as JSON and sending it
	json.NewEncoder(w).Encode(response)
}

func (h *AuthenticationHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate user data
	if err := h.validateUserData(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.PasswordHash = hashedPassword

	if err := h.createUserRecord(&user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	// Respond with success or user data (excluding sensitive information like password)
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message": "User successfully registered",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	}
	json.NewEncoder(w).Encode(response)
}

// validateRefreshToken checks if the refresh token is valid.
func validateRefreshToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method in token")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if token is valid and not expired
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// You can also check for the "exp" claim here manually, but jwt-go does it for you
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return nil, errors.New("token is expired")
			}
			return token, nil
		}
	}
	return nil, errors.New("invalid token")
}

func (h *AuthenticationHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var tokenDetails models.TokenDetails
	if err := json.NewDecoder(r.Body).Decode(&tokenDetails); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the refresh token
	token, err := validateRefreshToken(tokenDetails.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Extract user ID or other needed info from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to parse token claims", http.StatusInternalServerError)
		return
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Failed to retrieve user ID from token", http.StatusInternalServerError)
		return
	}

	// Generate new tokens
	accessToken, refreshToken, err := h.GenerateToken(userID)
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// Return the new tokens
	json.NewEncoder(w).Encode(models.TokenDetails{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		// Set expiration times as needed
	})
}
