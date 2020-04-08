package model

import (
	"errors"
	"strings"
)

// User holds every data related to a user
type User struct {
	Name     string
	LastName string
	Balance  float32 // Current balance of the user
	Orders   []int   // Orders of the user
}

// UserDB holds the needed data for every DB operation to run
type UserDB struct {
	Usr User             // user whom the operations will be on
	db  map[string]*User // holds every User created
}

// NewUserDB creates a UserDB struct with empty initial values and returns it.
func NewUserDB() UserDB {
	return UserDB{User{}, make(map[string]*User)}
}

// NewUser creates a User with the given name and last name and returns it.
// Every other data other than name and last name are set to their initial values.
func NewUser(name, lastName string) (User, error) {
	var u User

	if err := checkName(name, lastName); err != nil {
		return u, err
	}

	u.Name = name
	u.LastName = lastName

	return u, nil
}

// UpdateBalance function adds the given balance to the given user in UserDB struct.
func (udb *UserDB) UpdateBalance(balance float32) (float32, error) {
	key := generateKeyForUser(&udb.Usr)
	if _, ok := udb.db[key]; !ok {
		return 0, errors.New("user not found")
	}

	usr := udb.db[key]
	usr.Balance = usr.Balance + balance

	return usr.Balance, nil
}

// AddToDB function adds user in UserDB to the DB.
func (udb *UserDB) AddToDB() error {
	if err := checkName(udb.Usr.Name, udb.Usr.LastName); err != nil {
		return err
	}

	key := generateKeyForUser(&udb.Usr)
	if _, ok := udb.db[key]; ok {
		return errors.New("user already exists")
	}

	udb.db[key] = &udb.Usr

	return nil
}

// Find function finds the user from db and returns it.
// An error is returned if key does not exist in DB map.
func (udb *UserDB) Find(key string) (interface{}, error) {
	if usr, ok := udb.db[key]; ok {
		return usr, nil
	}

	return nil, errors.New("user not found")
}

// Delete function removes the user associated with the key in parameter from the DB.
func (udb *UserDB) Delete(key string) bool {
	if _, ok := udb.db[key]; ok {
		delete(udb.db, key)
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
