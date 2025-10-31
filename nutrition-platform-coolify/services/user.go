package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserData struct {
	Users    []User                 `json:"users"`
	Metadata map[string]interface{} `json:"metadata"`
}

// CreateUser creates a new user
func CreateUser(email, password, firstName, lastName string) (*User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
		Role:      "user",
		IsAdmin:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Check if user already exists
	existingUser, _ := GetUserByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Add user to storage
	err = AppendJSON("users.json", user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retrieves user by email
func GetUserByEmail(email string) (*User, error) {
	var userData UserData
	err := ReadJSON("users.json", &userData)
	if err != nil {
		return nil, err
	}

	for _, user := range userData.Users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

// GetUserByID retrieves user by ID
func GetUserByID(id string) (*User, error) {
	var userData UserData
	err := ReadJSON("users.json", &userData)
	if err != nil {
		return nil, err
	}

	for _, user := range userData.Users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

// AuthenticateUser authenticates user credentials
func AuthenticateUser(email, password string) (*User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// UpdateUser updates user information
func UpdateUser(userID string, updates map[string]interface{}) (*User, error) {
	var userData UserData
	err := ReadJSON("users.json", &userData)
	if err != nil {
		return nil, err
	}

	for i, user := range userData.Users {
		if user.ID == userID {
			// Update fields
			if firstName, ok := updates["first_name"].(string); ok {
				userData.Users[i].FirstName = firstName
			}
			if lastName, ok := updates["last_name"].(string); ok {
				userData.Users[i].LastName = lastName
			}
			if email, ok := updates["email"].(string); ok {
				userData.Users[i].Email = email
			}
			userData.Users[i].UpdatedAt = time.Now()

			err = WriteJSON("users.json", userData)
			if err != nil {
				return nil, err
			}

			return &userData.Users[i], nil
		}
	}

	return nil, errors.New("user not found")
}

// DeleteUser deletes a user
func DeleteUser(userID string) error {
	var userData UserData
	err := ReadJSON("users.json", &userData)
	if err != nil {
		return err
	}

	for i, user := range userData.Users {
		if user.ID == userID {
			// Remove user from slice
			userData.Users = append(userData.Users[:i], userData.Users[i+1:]...)

			// Update metadata
			if userData.Metadata != nil {
				if count, ok := userData.Metadata["total_count"].(float64); ok {
					userData.Metadata["total_count"] = count - 1
				}
				userData.Metadata["last_updated"] = time.Now().Format(time.RFC3339)
			}

			return WriteJSON("users.json", userData)
		}
	}

	return errors.New("user not found")
}

// GenerateSalt generates a random salt
func GenerateSalt() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashPassword hashes password with salt
func HashPassword(password, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return hex.EncodeToString(hash[:])
}
