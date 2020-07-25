package user

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User type to represent a user
type User struct {
	Email    string
	Name     string
	password []byte
}

// PasswordMatching validates if password matching with users hashPassword
func (u *User) PasswordMatching(password []byte) bool {
	err := bcrypt.CompareHashAndPassword(u.password, password)
	return err == nil
}

var userDB map[string]User

func init() {
	userDB = make(map[string]User)
	NewUser("jaimejimenezisi@gmail.com", "Jaime Jimenez", "12345678") // creating my account
}

// GetUser returns a user by email if exist or nil isn't exist
func GetUser(email string) *User {
	if u, ok := userDB[email]; ok {
		return &u
	}
	return nil
}

// NewUser adds a new user to memory DB and returns it, if it isn't exist
func NewUser(email, name, password string) (*User, error) {
	if email == "" || name == "" || password == "" {
		return nil, errors.New("There are empty fields")
	}
	if GetUser(email) != nil {
		return nil, fmt.Errorf("User with email %s, already exist", email)
	}
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newUser := User{email, name, cryptedPassword}
	userDB[email] = newUser
	return &newUser, nil
}
