package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Welcome to school")
	})
	mux.HandleFunc("/class", (getClassInfo))
	mux.HandleFunc("/student/", checkAuth(getStudentInfo))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error happened", err.Error())
		return
	}
}

type Student struct {
	ID    int
	Name  string
	Marks []int
}

type Class struct {
	Name     string
	Students []Student
}

var Minnie Student = Student{ID: 1, Name: "Minnie Mouse", Marks: []int{5, 5, 5}}
var Donald Student = Student{ID: 2, Name: "Donald Duck", Marks: []int{4, 3, 5}}
var Mickey Student = Student{ID: 3, Name: "Mickey Mouse", Marks: []int{5, 4, 5}}

//var listOfStudents = []Student{
//	{ID: 1, Name: "Minnie Mouse", Marks: []int{5, 5, 5}},
//	{ID: 2, Name: "Donald Duck", Marks: []int{4, 3, 5}},
//	{ID: 3, Name: "Mickey Mouse", Marks: []int{5, 4, 5}},
//}

var class1 = Class{
	Name: "Disney class",
	//	Students: listOfStudents,
}

type Authorisation struct {
	UserName string
	Password string
}

var Teacher = Authorisation{
	UserName: "Maria",
	Password: "yes",
}

var ReadOnly = Authorisation{
	UserName: "George",
	Password: "no",
}

func checkAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if Teacher.UserName != username || Teacher.Password != password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func getClassInfo(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(class1)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getStudentInfo(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.TrimPrefix(r.URL.Path, "/student/")
	id, err := strconv.Atoi(urlPath)
	if err != nil {
		fmt.Println("Error converting id:", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var student Student
	switch id {
	case 1:
		student = Minnie
	case 2:
		student = Donald
	case 3:
		student = Mickey
	default:
		fmt.Println("Student not found for id:", id)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		fmt.Println("Error encoding student:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
