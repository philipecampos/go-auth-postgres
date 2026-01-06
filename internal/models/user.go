package models

import (
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	var errorMessages []string

	if len(u.Username) < 6 {
		errorMessages = append(errorMessages, "username must be at least 6 characters")
	}

	if len(u.Email) < 6 {
		errorMessages = append(errorMessages, "email address must be at least 6 characters")
	}

	if len(u.Password) < 6 {
		errorMessages = append(errorMessages, "password must be at least 6 characters")
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("%s", strings.Join(errorMessages, ", "))
	}

	return nil
}

func NewUser(username, password, email string) (*User, error) {
	user := &User{
		ID:        0,
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

type UserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
