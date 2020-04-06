package model

import (
	"errors"
	"strings"
)

type User struct {
	Name     string
	LastName string
	Orders   []int
}

type UserDB struct {
	Usr User
	db  map[string]*User
}

func NewUserDB() UserDB {
	return UserDB{User{}, make(map[string]*User)}
}

func NewUser(name, lastName string) (User, error) {
	var u User

	if err := checkName(name, lastName); err != nil {
		return u, err
	}

	u.Name = name
	u.LastName = lastName

	return u, nil
}

func (udb *UserDB) AddToDB() error {
	if err := checkName(udb.Usr.Name, udb.Usr.LastName); err != nil {
		return  err
	}

	key := generateKeyForUser(&udb.Usr)
	if _, ok := udb.db[key]; ok {
		return errors.New("user already exists")
	}

	udb.db[key] = &udb.Usr

	return nil
}

func (udb *UserDB) Find(key string) (interface{}, error) {
	if usr, ok := udb.db[key]; ok {
		return usr, nil
	}

	return nil, errors.New("user not found")
}

func (udb *UserDB) Delete(key string) bool {
	if _, ok := udb.db[key]; ok {
		delete(udb.db, key)
		return true
	}

	return false
}

func checkName(name, lastName string ) error {
	if name == "" || lastName == "" {
		return errors.New("name or last name cannot be empty")
	}

	return nil
}

func generateKeyForUser(u *User) string {
	return strings.ToLower(u.Name) + "-" + strings.ToLower(u.LastName)
}
