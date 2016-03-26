package main

import (
	"errors"
	"time"

	"github.com/couchbase/gocb"
	hashids "github.com/speps/go-hashids"
)

// Project Struct Type
type Project struct {
	Type        string   `json:"_type"`
	ID          string   `json:"_id"`
	CreatedOn   string   `json:"createdON"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Owner       string   `json:"owner"`
	Status      string   `json:"status"`
	Tasks       []string `json:"tasks"`
	Users       []string `json:"users"`
	Permalink   string   `json:"permalink"`
}

type SessionProject struct {
	Project Project
}

func (s *SessionProject) Create() (*Project, error) {
	// Generate a faux uuid
	uuid := GenUUID()
	t := int64(time.Now().UnixNano())

	// Use hashids to create a permalink to store in the database, and for easy
	// serve/retrieve from external systems.  Note, this is NOT secure
	hd := hashids.NewData()
	hd.Salt = "testApp"
	h := hashids.NewWithData(hd)
	e, _ := h.EncodeInt64([]int64{t})

	// Create a new project instance.   This method uses the Project struct within
	// the SessionProject struct when it's passed in.   It adds the specific
	// item fields not set by the rest endpoint in the JSON body, and then stores
	// within the appropriate bucket.
	s.Project.Type = "Project"
	s.Project.ID = uuid
	s.Project.CreatedOn = time.Now().Format(time.RFC3339)
	s.Project.Status = "active"
	s.Project.Permalink = e
	s.Project.Users = append(s.Project.Users, s.Project.Owner)

	// Store in couchbase, check for error.   If no errors, return the project
	// object back to the front end application.
	_, err := bucket.Upsert(s.Project.ID, s.Project, 0)
	if err != nil {
		return nil, err
	}
	return &s.Project, nil
}

func (s *SessionProject) Retrieve(id string) (interface{}, error) {
	// Retrieve a single Project instance from the id.  Uses a get operation
	// against the database and returns the Project object to the front end
	// application.
	myQuery := gocb.NewN1qlQuery("SELECT _id, createdON, description,name, " +
		"(SELECT _id,_type,active,address, company,createdON,name,`password`,phone " +
		"FROM`comply` USE KEYS c.owner)[0] as owner,(SELECT _id,_type,active, " +
		"address,company,createdON,name,`password`,phone FROM `comply` USE KEYS " +
		"c.users) AS users, (SELECT _id,name,description,owner,assignedTo,Users, " +
		"history,permalink FROM `comply` USE KEYS c.tasks) as tasks, permalink from " +
		" `comply` c WHERE c._id=$1")
	var myParams []interface{}
	myParams = append(myParams, id)

	rows, err := bucket.ExecuteN1qlQuery(myQuery, myParams)
	if err != nil {
		return nil, err
	}
	// Interface for return values
	var retValues interface{}

	// Put the query results into the return interface
	_ = rows.One(&retValues)
	return retValues, nil
}

func (s *SessionProject) RetrieveAll() ([]Project, error) {
	// Retrieves all Projects from the database.   Does not implement filtering
	// and returns the array of users back to the front end application.
	myQuery := gocb.NewN1qlQuery("SELECT * FROM comply " +
		"WHERE _type = 'Project'")
	rows, err := bucket.ExecuteN1qlQuery(myQuery, nil)
	if err != nil {
		return nil, err
	}

	// Wrapper struct needed for parsing results from N1QL
	// The results will always come wrapped in the "bucket" name
	type wrapProject struct {
		Project Project `json:"comply"`
	}

	// Temporary variables to parse the results.
	var row wrapProject
	var curProjects []Project

	// Parse the n1ql results and build the array of Projects to return to the
	// front end application
	for rows.Next(&row) {
		curProjects = append(curProjects, row.Project)
	}
	return curProjects, nil
}

func (s *SessionProject) RetrieveOwned(owner string) ([]Project, error) {
	// Retrieves all Projects from the database.   Includes filtering by owner
	// and returns the array of projects back to the front end application.
	myQuery := gocb.NewN1qlQuery("SELECT * FROM comply " +
		"WHERE _type = 'Project' and owner = $1")
	var myParams []interface{}
	myParams = append(myParams, owner)
	rows, err := bucket.ExecuteN1qlQuery(myQuery, myParams)
	if err != nil {
		return nil, err
	}

	// Wrapper struct needed for parsing results from N1QL
	// The results will always come wrapped in the "bucket" name
	type wrapProject struct {
		Project Project `json:"comply"`
	}

	// Temporary variables to parse the results.
	var row wrapProject
	var curProjects []Project

	// Parse the n1ql results and build the array of Projects to return to the
	// front end application
	for rows.Next(&row) {
		curProjects = append(curProjects, row.Project)
	}
	return curProjects, nil
}

func (s *SessionProject) RetrieveMemberOf(member string) ([]Project, error) {
	// Retrieves all Projects from the database.   Includes filtering by member
	// and returns the array of projects back to the front end application.
	myQuery := gocb.NewN1qlQuery("SELECT * FROM comply " +
		"WHERE _type = 'Project' AND ANY x IN users SATISFIES x = $1 END")
	var myParams []interface{}
	myParams = append(myParams, member)
	rows, err := bucket.ExecuteN1qlQuery(myQuery, myParams)
	if err != nil {
		return nil, err
	}

	// Wrapper struct needed for parsing results from N1QL
	// The results will always come wrapped in the "bucket" name
	type wrapProject struct {
		Project Project `json:"comply"`
	}

	// Temporary variables to parse the results.
	var row wrapProject
	var curProjects []Project

	// Parse the n1ql results and build the array of Projects to return to the
	// front end application
	for rows.Next(&row) {
		curProjects = append(curProjects, row.Project)
	}
	return curProjects, nil
}

func (s *SessionProject) AddUserToProject(projectID string, userID string) (*User, error) {
	// Retrieve a single User instance from the userID.  Uses a get operation
	// against the database and adds the userID to the project and returns the
	// user instance to the front end application.
	var curUser User
	_, err := bucket.Get(userID, &curUser)
	if err != nil {
		return nil, err
	}
	_, err = bucket.Get(projectID, &s.Project)
	if err != nil {
		return nil, err
	}
	if SliceItemExists(userID, s.Project.Users) {
		return nil, errors.New("User already exists")
	}
	s.Project.Users = append(s.Project.Users, userID)
	_, err = bucket.Upsert(s.Project.ID, s.Project, 0)
	if err != nil {
		return nil, err
	}
	return &curUser, nil
}
