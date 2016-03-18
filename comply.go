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

func AddDefaultCompany() {
	var defaultCompany SessionCompany
	defaultCompany.Company.Active = true
	defaultCompany.Company.Name = "Couchbase"
	defaultCompany.Company.Address.Street = "2440 W El Camino Real #101"
	defaultCompany.Company.Address.City = "Mountain View"
	defaultCompany.Company.Address.State = "California"
	defaultCompany.Company.Address.Zip = 94040
	defaultCompany.Company.Phone = "650-417-7500"
	defaultCompany.Company.Website = "www.couchbase.com"

	curCompany, err := defaultCompany.Create()
	if err != nil {
		fmt.Println("Error Creating Default Company:", err)
		return
	}
	fmt.Println("Default Company Created.")
	fmt.Println(curCompany)
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

	// Static Directories for Angular 2.0 APP
	p := http.StripPrefix("/", http.FileServer(http.Dir("./public/src/")))
	n := http.StripPrefix("/node_modules", http.FileServer(http.Dir("./node_modules/")))
	r.PathPrefix("/node_modules/").Handler(n)
	r.PathPrefix("/").Handler(p)

	fmt.Printf("Starting server on :3000\n")
	http.ListenAndServe(":3000", r)
}
