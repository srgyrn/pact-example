package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// User holds every data related to a user
type User struct {
	Name     string  `json:Name`
	LastName string  `json:LastName`
	Balance  float32 `json:Balance` // Current balance of the user
	Orders   []int   // Orders of the user
}

// UserHandler holds the needed data for every DB operation to run
type UserHandler struct {
	Usr *User            // user whom the operations will be on
	db  map[string]*User // holds every User created
}

// NewUserHandler creates a UserDB struct with empty initial values and returns it.
func NewUserHandler() *UserHandler {
	return &UserHandler{nil, make(map[string]*User)}
}

// NewUser creates a User with the given name and last name and returns it.
// Every other data other than name and last name are set to their initial values.
func NewUser(name, lastName string) (*User, error) {
	if err := checkName(name, lastName); err != nil {
		return nil, err
	}

	return &User{
		Name:     name,
		LastName: lastName,
	}, nil
}

func (uh *UserHandler) BulkInsert(b []byte) error {
	err := json.Unmarshal(b, &uh.db)
	if err != nil {
		return fmt.Errorf("failed to unmarshal\n%s", err)
	}

	return nil
}

// UpdateBalance function adds the given balance to the given user in UserDB struct.
func (u *User) UpdateBalance(balance float32) (float32, error) {
	u.Balance = u.Balance + balance
	return u.Balance, nil
}

// AddToDB function adds user in UserDB to the DB.
func (uh *UserHandler) AddToDB() error {
	if err := checkName(uh.Usr.Name, uh.Usr.LastName); err != nil {
		return err
	}

	key := generateKeyForUser(uh.Usr)
	if _, ok := uh.db[key]; ok {
		return errors.New("user already exists")
	}

	uh.db[key] = uh.Usr

	return nil
}

// Find function finds the user from db and returns it.
// An error is returned if key does not exist in DB map.
func (uh *UserHandler) Find(key string) (interface{}, error) {
	if usr, ok := uh.db[key]; ok {
		return usr, nil
	}

	return nil, errors.New("user not found")
}

// Delete function removes the user associated with the key in parameter from the DB.
func (uh *UserHandler) Delete(key string) bool {
	if _, ok := uh.db[key]; ok {
		delete(uh.db, key)
		return true
	}

	return false
}

func checkName(name, lastName string) error {
	if name == "" || lastName == "" {
		return errors.New("name or last name cannot be empty")
	}

	return nil
}

// generateKeyForUser is a helper function to create a key for the user
func generateKeyForUser(u *User) string {
	return strings.ToLower(u.Name) + "-" + strings.ToLower(u.LastName)
}
