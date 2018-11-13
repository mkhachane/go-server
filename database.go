package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

import _ "github.com/go-sql-driver/mysql"


type Employee struct {
	Id    int
	Name  string
	Dept string
	City string
	Sal int
}


func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:root@/goServer")
	if err != nil {
		panic(err.Error())
	}
	return db
}


func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM employee ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id, sal int
		var name, city, dept string
		err = selDB.Scan(&id, &name, &dept, &city, &sal)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Dept = dept
		emp.City = city
		emp.Sal = sal
		res = append(res, emp)
	}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id, sal int
		var name, city, dept string
		err = selDB.Scan(&id, &name, &dept, &city, &sal)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Dept = dept
		emp.City = city
		emp.Sal = sal
	}
	t, _ := template.ParseFiles("show.html")
	t.Execute(w, emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("new.html")
	t.Execute(w, nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		dept := r.FormValue("dept")
		city := r.FormValue("city")
		sal := r.FormValue("sal")
		insForm, err := db.Prepare("INSERT INTO employee(name, dept, city, sal) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, dept, city, sal)
		log.Println("INSERT: Name: " + name + " | Dept: " + dept + " | City: " + city + " | Sal: " + sal)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id, sal int
		var name, city, dept string
		err = selDB.Scan(&id, &name, &dept, &city, &sal)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Dept = dept
		emp.City = city
		emp.Sal = sal
	}
	t, _ := template.ParseFiles("update.html")
	t.Execute(w, emp)
	defer db.Close()
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		dept := r.FormValue("dept")
		city := r.FormValue("city")
		sal := r.FormValue("sal")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE employee SET name=?, dept=?, city=?, sal=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, dept, city, sal, id)
		log.Println("INSERT: Name: " + name + " | Dept: " + dept + " | City: " + city + " | Sal: " + sal)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM employee WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("COLUMN DELETED")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}