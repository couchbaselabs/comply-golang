package main

import (
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

func main() {
	// Cluster connection and bucket for couchbase
	cluster, _ := gocb.Connect("couchbase://192.168.99.100")
	bucket, _ = cluster.OpenBucket("comply", "")
	// Http Routing
	r := mux.NewRouter()
	r.HandleFunc("/api/user/login/{email}/{pass}", LoginHandler).Methods("GET")
	r.HandleFunc("/api/user/getAll", RetrieveAllUserHandler).Methods("GET")
	r.HandleFunc("/api/user/get/{userId}", RetrieveUserHandler).Methods("GET")
	r.HandleFunc("/api/user/create", CreateLoginHandler).Methods("POST")

	r.HandleFunc("/api/company/get/{companyId}", RetrieveCompanyHandler).Methods("GET")
	r.HandleFunc("/api/company/getAll", RetrieveAllCompanyHandler).Methods("GET")
	r.HandleFunc("/api/company/create", CreateCompanyHandler).Methods("POST")

	//r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	//http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Printf("Starting server on :3000\n")
	http.ListenAndServe(":3000", r)
}
