package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/couchbase/gocb"
)

type Session struct {
	User User
}

// User Struct Type
type User struct {
	Type      string `json:"_type"`
	ID        string `json:"_id"`
	CreatedOn string `json:"createdON"`
	Name      struct {
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"name"`
	Address struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		State   string `json:"state"`
		Zip     string `json:"zip"`
		Country string `json:"country"`
	} `json:"address"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Company  string `json:"company"`
	Active   bool   `json:"active"`
}

func (u *Session) Login(email string, password string) (*User, error) {
	// Setup Query and Execute
	myQuery := gocb.NewN1qlQuery("SELECT * FROM `comply` " +
		"WHERE _type = 'User' and email='" + email + "' ")
	rows, err := bucket.ExecuteN1qlQuery(myQuery, nil)
	var wrapUser struct {
		User User `json:"comply"`
	}
	err = rows.One(&wrapUser)
	if err != nil {
		return nil, errors.New("User Not Found")
	}
	if wrapUser.User.Password == password {
		return &wrapUser.User, nil
	}
	fmt.Println("DEBUG:PASS:", wrapUser.User.Password)
	return nil, errors.New("Password is invalid")
}

func (u *Session) Create() (*User, error) {
	u.User.Type = "User"
	u.User.ID = u.User.Email
	u.User.CreatedOn = time.Now().Format(time.RFC3339)
	u.User.Active = true

	_, err := bucket.Upsert(u.User.Email, u.User, 0)
	if err != nil {
		return nil, err
	}
	return &u.User, nil
}

func (u *Session) Retrieve(id string) (*User, error) {
	_, err := bucket.Get(id, &u.User)
	if err != nil {
		return nil, err
	}
	return &u.User, nil
}

func (u *Session) RetrieveAll() ([]User, error) {
	myQuery := gocb.NewN1qlQuery("SELECT * FROM `comply` " +
		"WHERE _type = 'User'")
	rows, err := bucket.ExecuteN1qlQuery(myQuery, nil)
	if err != nil {
		return nil, err
	}

	type wrapUser struct {
		User User `json:"comply"`
	}
	var row wrapUser
	var curUsers []User

	for rows.Next(&row) {
		curUsers = append(curUsers, row.User)
	}
	return curUsers, nil
}
