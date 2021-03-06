package main

import (
	"errors"
	"time"

	"github.com/couchbase/gocb"
	hashids "github.com/speps/go-hashids"
)

// Task Struct Type
type Task struct {
	URL         string   `json:"url"`
	Type        string   `json:"_type"`
	ID          string   `json:"_id"`
	CreatedOn   string   `json:"createdON"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Owner       string   `json:"owner"`
	AssignedTo  string   `json:"assignedTo"`
	Status      string   `json:"status"`
	Priority    string   `json:"priority"`
	Tasks       []string `json:"tasks"`
	Users       []string `json:"users"`
	Permalink   string   `json:"permalink"`
	History     []struct {
		Log       string `json:"log"`
		User      string `json:"user"`
		CreatedAt string `json:"createdAt"`
		photos    []struct {
			Filename  string `json:"filename"`
			Extension string `json:"extension"`
		} `json:"photos"`
	} `json:"history"`
}

type SessionTask struct {
	Task Task
}

func (s *SessionTask) Create(projectID string) (interface{}, error) {
	// Generate a faux uuid
	uuid := GenUUID()
	t := int64(time.Now().UnixNano())

	// Use hashids to create a permalink to store in the database, and for easy
	// serve/retrieve from external systems.  Note, this is NOT secure
	hd := hashids.NewData()
	hd.Salt = "testApp"
	h := hashids.NewWithData(hd)
	e, _ := h.EncodeInt64([]int64{t})

	// Create a new task instance.   This method uses the task struct within
	// the SessionTask struct when it's passed in.   It adds the specific
	// item fields not set by the rest endpoint in the JSON body, and then stores
	// within the appropriate bucket.
	s.Task.Type = "Task"
	s.Task.ID = uuid
	s.Task.CreatedOn = time.Now().Format(time.RFC3339)
	s.Task.Status = "active"
	s.Task.Permalink = e
	s.Task.Users = append(s.Task.Users, s.Task.Owner)

	// Store in couchbase, check for error.   I
	_, err := bucket.Upsert(s.Task.ID, s.Task, 0)
	if err != nil {
		return nil, err
	}

	// Now retrieve the current project to add the reference to the task
	var curProject Project
	_, err = bucket.Get(projectID, &curProject)
	if err != nil {
		return nil, err
	}

	// Add the task reference to the project
	curProject.Tasks = append(curProject.Tasks, uuid)
	_, err = bucket.Upsert(curProject.ID, curProject, 0)
	if err != nil {
		return nil, err
	}

	// Setup a new N1QL Query to retrieve and assemble the complete json doc and
	// return to the front end application
	myQuery := gocb.NewN1qlQuery("SELECT c._id,c.createdON,c.name,c.description," +
		"(SELECT _id,_type,active,address,company,createdON,name,`password`,phone " +
		"FROM`comply` USE KEYS c.owner)[0] AS owner,c.status, (SELECT _id,_type," +
		"active,address,company,createdON,name,`password`,phone FROM `comply` " +
		" USE KEYS c.users) AS users, c.permalink from `comply` " +
		" c WHERE c._id=$1 ").Consistency(gocb.RequestPlus)

	// Build the parameters interface to replace the $1 with the correct parameter
	var myParams []interface{}
	myParams = append(myParams, uuid)

	// Execute N1QL Query
	rows, err := bucket.ExecuteN1qlQuery(myQuery, myParams)
	if err != nil {
		return nil, err
	}

	// Interfaces for handling return values
	var retValues interface{}

	// Stream the values returned from the query into an untyped interfaces
	err = rows.One(&retValues)
	if err != nil {
		return nil, err
	}

	return retValues, nil

}

func (s *SessionTask) Retrieve(id string) (interface{}, error) {
	// Retrieve a single task instance from the id.  Uses a N1QL Query
	// against the database and returns the Document to the front end
	// application.

	myQuery := gocb.NewN1qlQuery("SELECT (p._id) AS projectId,(SELECT _id, " +
		"(SELECT _id,_type,active,address,company,createdON,name,`password`,phone " +
		"FROM `comply` USE KEYS c.assignedTo)[0] AS assignedTo, createdON, " +
		"description,(select t.log, t.createdAt, (SELECT _id,_type,active,address, " +
		"company,createdON,name,`password`,phone FROM `comply` USE KEYS t.`user`)[0] " +
		"as `user` from `comply` r UNNEST r.history t where r._id=$1) as history,name, " +
		"(SELECT _id,_type,active,address,company, createdON,name,`password`,phone " +
		"FROM`comply` USE KEYS c.owner)[0] as owner,(SELECT _id,_type,active,address, " +
		"company,createdON,name,`password`,phone FROM `comply` USE KEYS c.users) " +
		"AS users, permalink from `comply` c WHERE c._id=$1)[0] as task FROM `comply` " +
		"p WHERE ANY x IN tasks SATISFIES x=$1 END ")

	// Parameters interface to replace $1 with correct parameter
	var myParams []interface{}
	myParams = append(myParams, id)
	rows, err := bucket.ExecuteN1qlQuery(myQuery, myParams)
	if err != nil {
		return nil, err
	}

	// Interfaces for handling streaming return values
	var retValues interface{}

	// Stream the values returned from the query into an untyped and unstructred
	// array of interfaces
	err = rows.One(&retValues)
	if err != nil {
		return nil, err
	}
	return retValues, nil
}

func (s *SessionTask) RetrieveAssignedToUser(userID string) ([]interface{}, error) {
	// Retrieve a single task instance assigned to a user from the id.  Uses a
	// N1QL Query operation against the database and returns the JSON Document
	// to the front end application.

	// Setup new N1QL Query
	myQuery := gocb.NewN1qlQuery("SELECT _id,(SELECT _id,_type,active," +
		"address,company,createdON,name,`password`,phone FROM `comply` USE KEYS " +
		"c.assignedTo)[0] AS assignedTo, createdON, description,history,name," +
		"(SELECT _id,_type,active,address,company,createdON,name,`password`,phone " +
		"FROM`comply` USE KEYS c.owner)[0] as owner,(SELECT _id,_type,active,address," +
		"company,createdON,name,`password`,phone FROM `comply` USE KEYS c.users) AS " +
		"users, permalink from `comply` c WHERE c.assignedTo=$1")

	// Parameters interface to replace $1 with correct parameter
	var myParams []interface{}
	myParams = append(myParams, userID)

	// Exeute Query
	rows, err := bucket.ExecuteN1qlQuery(myQuery, myParams)
	if err != nil {
		return nil, err
	}
	// Interfaces for handling streaming return values
	var retValues []interface{}
	var row interface{}

	// Stream the values returned from the query into an untyped and unstructred
	// array of interfaces
	for rows.Next(&row) {
		retValues = append(retValues, row)
	}

	return retValues, nil
}
func (s *SessionTask) AddUserToTask(taskID string, userID string) (*User, error) {
	// Retrieve a single User instance from the userID.  Uses a get operation
	// against the database and adds the userID to the task and returns the
	// user instance to the front end application.
	var curUser User
	_, err := bucket.Get(userID, &curUser)
	if err != nil {
		return nil, err
	}
	_, err = bucket.Get(taskID, &s.Task)
	if err != nil {
		return nil, err
	}
	if SliceItemExists(userID, s.Task.Users) {
		return nil, errors.New("User already exists")
	}
	s.Task.Users = append(s.Task.Users, userID)
	_, err = bucket.Upsert(s.Task.ID, s.Task, 0)
	if err != nil {
		return nil, err
	}
	return &curUser, nil
}

func (s *SessionTask) AssignToUser(taskID string, userID string) (*User, error) {
	// Retrieve a single User instance from the userID.  Uses a get operation
	// against the database and assigns the userID to the task and returns the
	// user instance to the front end application.
	var curUser User
	_, err := bucket.Get(userID, &curUser)
	if err != nil {
		return nil, err
	}
	_, err = bucket.Get(taskID, &s.Task)
	if err != nil {
		return nil, err
	}
	s.Task.AssignedTo = userID
	_, err = bucket.Upsert(s.Task.ID, s.Task, 0)
	if err != nil {
		return nil, err
	}
	return &curUser, nil
}
func (s *SessionTask) AddHistoryToTask(taskID string, userID string, log string) (interface{}, error) {
	var curHistory struct {
		Log       string `json:"log"`
		User      string `json:"user"`
		CreatedAt string `json:"createdAt"`
		photos    []struct {
			Filename  string `json:"filename"`
			Extension string `json:"extension"`
		} `json:"photos"`
	}

	curHistory.CreatedAt = time.Now().Format(time.RFC3339)
	curHistory.User = userID
	curHistory.Log = log

	_, err := bucket.Get(taskID, &s.Task)
	if err != nil {
		return nil, err
	}

	s.Task.History = append(s.Task.History, curHistory)
	_, err = bucket.Upsert(s.Task.ID, s.Task, 0)
	if err != nil {
		return nil, err
	}

	// Setup a new N1QL Query to retrieve and assemble the complete json doc and
	// return to the front end application
	myQuery := gocb.NewN1qlQuery("SELECT ($1) AS log, (SELECT _id,_type,active," +
		"address,company,createdON,name,`password`,phone " +
		"FROM`comply` USE KEYS c._id)[0] AS `user`,($2) AS createdAt " +
		" FROM `comply` c WHERE c._id=$3 ").Consistency(gocb.RequestPlus)

	// Build the parameters interface to replace the $1 with the correct parameter
	var myParams []interface{}
	myParams = append(myParams, log)
	myParams = append(myParams, curHistory.CreatedAt)
	myParams = append(myParams, userID)

	// Execute N1QL Query
	rows, err := bucket.ExecuteN1qlQuery(myQuery, myParams)
	if err != nil {
		return nil, err
	}

	// Interfaces for handling return values
	var retValues interface{}

	// Stream the values returned from the query into an untyped interfaces
	err = rows.One(&retValues)
	if err != nil {
		return nil, err
	}

	return retValues, nil

}
