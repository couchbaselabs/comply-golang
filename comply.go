package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/couchbase/gocb"
	"github.com/gorilla/mux"
)

// bucket reference
var bucket *gocb.Bucket

func LoginHandler(w http.ResponseWriter, req *http.Request) {
	var session Session
	vars := mux.Vars(req)
	curUser, err := session.Login(vars["email"], vars["pass"])
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curUser)
}

func CreateLoginHandler(w http.ResponseWriter, req *http.Request) {
	var session Session
	_ = json.NewDecoder(req.Body).Decode(&session.User)
	curUser, err := session.Create()
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curUser)
}

func RetrieveUserHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	var session Session
	curUser, err := session.Retrieve(vars["userId"])
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curUser)
}

func RetrieveAllUserHandler(w http.ResponseWriter, req *http.Request) {
	var session Session
	curUser, err := session.RetrieveAll()
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curUser)
}

func CreateCompanyHandler(w http.ResponseWriter, req *http.Request) {
	var sessionCompany SessionCompany
	_ = json.NewDecoder(req.Body).Decode(&sessionCompany.Company)
	curCompany, err := sessionCompany.Create()
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curCompany)
}

func RetrieveCompanyHandler(w http.ResponseWriter, req *http.Request) {
	var sessionCompany SessionCompany
	vars := mux.Vars(req)
	curCompany, err := sessionCompany.Retrieve(vars["companyId"])
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curCompany)
}

func RetrieveAllCompanyHandler(w http.ResponseWriter, req *http.Request) {
	var sessionCompany SessionCompany
	curCompany, err := sessionCompany.RetrieveAll()
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curCompany)
}

func RetrieveAllProjectsHandler(w http.ResponseWriter, req *http.Request) {
	var sessionProject SessionProject
	curProject, err := sessionProject.RetrieveAll()
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curProject)
}

func RetrieveProjectHandler(w http.ResponseWriter, req *http.Request) {
	var sessionProject SessionProject
	vars := mux.Vars(req)
	curProject, err := sessionProject.Retrieve(vars["projectId"])
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curProject)
}

func RetrieveProjectMemberOfHandler(w http.ResponseWriter, req *http.Request) {
	var sessionProject SessionProject
	vars := mux.Vars(req)
	curProject, err := sessionProject.RetrieveMemberOf(vars["userId"])
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curProject)
}

func RetrieveProjectOwnerHandler(w http.ResponseWriter, req *http.Request) {
	var sessionProject SessionProject
	vars := mux.Vars(req)
	curProject, err := sessionProject.RetrieveOwned(vars["ownerId"])
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curProject)
}

func CreateProjectHandler(w http.ResponseWriter, req *http.Request) {
	var sessionProject SessionProject
	_ = json.NewDecoder(req.Body).Decode(&sessionProject.Project)
	curProject, err := sessionProject.Create()
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curProject)
}

func AddUserToProjectHandler(w http.ResponseWriter, req *http.Request) {
	var p struct {
		ProjectID string `json:"projectId"`
		Email     string `json:"email"`
	}
	_ = json.NewDecoder(req.Body).Decode(&p)
	var sessionProject SessionProject
	curUser, err := sessionProject.AddUserToProject(p.ProjectID, p.Email)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curUser)
}

func CreateTaskHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	var sessionTask SessionTask
	_ = json.NewDecoder(req.Body).Decode(&sessionTask.Task)
	fmt.Printf("%+v", sessionTask.Task)
	curTask, err := sessionTask.Create(vars["projectId"])
	fmt.Printf("%+v", curTask)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curTask)
}

func RetrieveTaskHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	var sessionTask SessionTask
	curTask, err := sessionTask.Retrieve(vars["taskId"])
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curTask)
}

func RetrieveTaskAssignedToUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	var sessionTask SessionTask
	curTask, err := sessionTask.RetrieveAssignedToUser(vars["userId"])
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if curTask == nil {
		// Special cae to not pass "null" back to front end framework
		bytes := []byte(`{"":""}`)
		w.Write(bytes)
	} else {
		json.NewEncoder(w).Encode(curTask)
	}
}

func AddDefaultCompany() {
	var t interface{}
	var defaultCompany SessionCompany
	defaultCompany.Company.Active = true
	defaultCompany.Company.Name = "Couchbase"
	defaultCompany.Company.Address.Street = "2440 W El Camino Real #101"
	defaultCompany.Company.Address.City = "Mountain View"
	defaultCompany.Company.Address.State = "California"
	defaultCompany.Company.Address.Zip = 94040
	defaultCompany.Company.Phone = "650-417-7500"
	defaultCompany.Company.Website = "www.couchbase.com"

	if _, err := bucket.Get(defaultCompany.Company.Website, &t); err == nil {
		fmt.Println("Default Company already created.")
		return
	}
	curCompany, err := defaultCompany.Create()
	if err != nil {
		fmt.Println("Error Creating Default Company:", err)
		return
	}
	fmt.Println("Default Company Created.")
	fmt.Println(curCompany)
}
func AddPrimaryIndex() {

}

func GenUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}

func SliceItemExists(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	// Cluster connection and bucket for couchbase
	cluster, _ := gocb.Connect("couchbase://192.168.99.100")
	bucket, _ = cluster.OpenBucket("comply", "")

	// Add a Default Company
	AddDefaultCompany()

	// Http Routing
	r := mux.NewRouter()

	// User Routes
	r.HandleFunc("/api/user/login/{email}/{pass}", LoginHandler).Methods("GET")
	r.HandleFunc("/api/user/getAll", RetrieveAllUserHandler).Methods("GET")
	r.HandleFunc("/api/user/get/{userId}", RetrieveUserHandler).Methods("GET")
	r.HandleFunc("/api/user/create", CreateLoginHandler).Methods("POST")

	// Company Routes
	r.HandleFunc("/api/company/get/{companyId}", RetrieveCompanyHandler).Methods("GET")
	r.HandleFunc("/api/company/getAll", RetrieveAllCompanyHandler).Methods("GET")
	r.HandleFunc("/api/company/create", CreateCompanyHandler).Methods("POST")

	// Project Routes
	r.HandleFunc("/api/project/getOther/{userId}", RetrieveProjectMemberOfHandler).Methods("GET")
	r.HandleFunc("/api/project/getAll/{ownerId}", RetrieveProjectOwnerHandler).Methods("GET")
	r.HandleFunc("/api/project/getAll", RetrieveAllProjectsHandler).Methods("GET")
	r.HandleFunc("/api/project/get/{projectId}", RetrieveProjectHandler).Methods("GET")
	r.HandleFunc("/api/project/create", CreateProjectHandler).Methods("POST")
	r.HandleFunc("/api/project/addUser", AddUserToProjectHandler).Methods("POST")

	// Task Routes
	r.HandleFunc("/api/task/create/{projectId}", CreateTaskHandler).Methods("POST")
	r.HandleFunc("/api/task/get/{taskId}", RetrieveTaskHandler).Methods("GET")
	r.HandleFunc("/api/task/getAssignedTo/{userId}", RetrieveTaskAssignedToUser).Methods("GET")

	// TO BE ADDED LATER
	//r.HandleFunc("/api/task/addUser").Methods("POST")
	//r.HandleFunc("/api/task/assignUser").Methods("POST")
	//r.HandleFunc("/api/task/addHistory").Methods("POST")
	//r.HandleFunc("/api/task/addPhoto").Methods("POST")

	// Static Directories for Angular 2.0 APP
	p := http.StripPrefix("/", http.FileServer(http.Dir("./public/src/")))
	n := http.StripPrefix("/node_modules", http.FileServer(http.Dir("./node_modules/")))
	r.PathPrefix("/node_modules/").Handler(n)
	r.PathPrefix("/").Handler(p)

	fmt.Printf("Starting server on :3000\n")
	http.ListenAndServe(":3000", r)
}
