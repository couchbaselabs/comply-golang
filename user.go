package main

import (
	"errors"
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
	// Login a single user.   Uses a N1QL query to retrieve a single instance
	// and return the user object back to the front end application.
	myQuery := gocb.NewN1qlQuery("SELECT * FROM `comply` " +
		"WHERE _type = 'User' and email='" + email + "' ")
	rows, err := bucket.ExecuteN1qlQuery(myQuery, nil)

	// Wrapper struct needed for parsing results from N1QL
	// The results will always come wrapped in the "bucket" name
	var wrapUser struct {
		User User `json:"comply"`
	}

	// Use the .One method to retrieve the first result.   As
	// this is retrieving a specific instance, we only need the first result
	err = rows.One(&wrapUser)
	if err != nil {
		return nil, errors.New("User Not Found")
	}

	// Check the password, compare.  If correct return the user object.
	// NOTE: for demonstration purposes, this approach is NOT SECURE
	if wrapUser.User.Password == password {
		return &wrapUser.User, nil
	}
	return nil, errors.New("Password is invalid")
}

func (u *Session) Create() (*User, error) {
	// Create a new user instance.   This method uses the User struct within the
	// Session struct when it's passed in.   It adds the specific items fields
	// Not set by the rest endpoint in the JSON body, and then stores within
	// the appropriate bucket.
	u.User.Type = "User"
	u.User.ID = u.User.Email
	u.User.CreatedOn = time.Now().Format(time.RFC3339)
	u.User.Active = true

	// Store in couchbase, check for error.   If no errors, return the user object
	// back to the front end application.
	_, err := bucket.Upsert(u.User.Email, u.User, 0)
	if err != nil {
		return nil, err
	}
	return &u.User, nil
}

func (u *Session) Retrieve(id string) (*User, error) {
	// Retrieve a single user instance from the id.  Uses a get operation against
	// the database and returns the User object to the front end application.
	_, err := bucket.Get(id, &u.User)
	if err != nil {
		return nil, err
	}
	return &u.User, nil
}

func (u *Session) RetrieveAll() ([]User, error) {
	// Retrieves all users from the database.   Does not implement filtering
	// and returns the array of users back to the front end application.
	myQuery := gocb.NewN1qlQuery("SELECT * FROM `comply` " +
		"WHERE _type = 'User'")
	rows, err := bucket.ExecuteN1qlQuery(myQuery, nil)
	if err != nil {
		return nil, err
	}

	// Wrapper struct needed for parsing results from N1QL
	// The results will always come wrapped in the "bucket" name
	type wrapUser struct {
		User User `json:"comply"`
	}

	// Temporary variables to parse the results.
	var row wrapUser
	var curUsers []User

	// Parse the n1ql results and build the array of Users to return to the
	// front end application
	for rows.Next(&row) {
		curUsers = append(curUsers, row.User)
	}
	return curUsers, nil
}
