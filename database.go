package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

import _ "github.com/go-sql-driver/mysql"


type Employee struct {
	Id    int
	Name  string
	Dept string
	City string
	Sal int
}

//database connection
func dbConn() (*sql.DB, error) {
	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	return db, err
}


func Index(w http.ResponseWriter, r *http.Request) {
	db, err := dbConn()
	defer db.Close()
	if err != nil {
		http.Error(w, "Database Connection Error!", 500)
		return
	}
	selectAllData, err := db.Query("SELECT * FROM employee ORDER BY id DESC")
	if err != nil {
		http.Error(w, "You have an error in your SQL syntax", 502)
		return
	}else {
		emp := Employee{}
		res := []Employee{}
		for selectAllData.Next() {
			var id, sal int
			var name, city, dept string
			err = selectAllData.Scan(&id, &name, &dept, &city, &sal)
			if err != nil {
				http.Error(w, err.Error(), 500)
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
	}
}

func Show(w http.ResponseWriter, r *http.Request) {
	db, err := dbConn()
	defer db.Close()
	if err != nil {
		http.Error(w, "Database Connection Error!", 500)
		return
	}
	nId := r.URL.Query().Get("id")
	showRecords, err := db.Query("SELECT * FROM employee WHERE id=?", nId)
	if err != nil {
		http.Error(w, "You have an error in your SQL syntax", 502)
		return
	}else {
		emp := Employee{}
		for showRecords.Next() {
			var id, sal int
			var name, city, dept string
			err = showRecords.Scan(&id, &name, &dept, &city, &sal)
			if err != nil {
				http.Error(w, err.Error(), 502)
			}
			emp.Id = id
			emp.Name = name
			emp.Dept = dept
			emp.City = city
			emp.Sal = sal
		}
		t, _ := template.ParseFiles("show.html")
		t.Execute(w, emp)
	}
}

func New(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("new.html")
	t.Execute(w, nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db, err := dbConn()
	defer db.Close()
	if err != nil {
		http.Error(w, "Database Connection Error!", 500)
		return
	}else {
		if r.Method == "POST" {
			name := r.FormValue("name")
			dept := r.FormValue("dept")
			city := r.FormValue("city")
			sal := r.FormValue("sal")
			insertRecord, err := db.Prepare("INSERT INTO employee(name, dept, city, sal) VALUES(?,?,?,?)")
			if err != nil {
				http.Error(w, "You have an error in your SQL syntax", 502)
				return
			}
			insertRecord.Exec(name, dept, city, sal)
			log.Println("INSERT: Name: " + name + " | Dept: " + dept + " | City: " + city + " | Sal: " + sal)
		}else {
			fmt.Fprintf(w, "Method Not Found")
		}
		http.Redirect(w, r, "/", 301)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db, err := dbConn()
	defer db.Close()
	if err != nil {
		http.Error(w, "Database Connection Error!", 500)
		return
	}else {
		nId := r.URL.Query().Get("id")
		records, err := db.Query("SELECT * FROM employee WHERE id=?", nId)
		if err != nil {
			http.Error(w, "You have an error in your SQL syntax", 502)
			return
		}
		emp := Employee{}
		for records.Next() {
			var id, sal int
			var name, city, dept string
			err = records.Scan(&id, &name, &dept, &city, &sal)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			emp.Id = id
			emp.Name = name
			emp.Dept = dept
			emp.City = city
			emp.Sal = sal
		}
		t, _ := template.ParseFiles("update.html")
		t.Execute(w, emp)
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	db, err := dbConn()
	defer db.Close()
	if err != nil {
		http.Error(w, "Database Connection Error!", 500)
		return
	}else {
		if r.Method == "POST" {
			name := r.FormValue("name")
			dept := r.FormValue("dept")
			city := r.FormValue("city")
			sal := r.FormValue("sal")
			id := r.FormValue("uid")
			updateRecords, err := db.Prepare("UPDATE employee SET name=?, dept=?, city=?, sal=? WHERE id=?")
			if err != nil {
				http.Error(w, "You have an error in your SQL syntax", 502)
				return
			}
			updateRecords.Exec(name, dept, city, sal, id)
			log.Println("INSERT: Name: " + name + " | Dept: " + dept + " | City: " + city + " | Sal: " + sal)
		}else {
			fmt.Fprintf(w, "Method Not Found")
		}
		http.Redirect(w, r, "/", 301)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db, err := dbConn()
	defer db.Close()
	if err != nil {
		http.Error(w, "Database Connection Error!", 500)
		return
	}else {
		emp := r.URL.Query().Get("id")
		deleteRecords, err := db.Prepare("DELETE FROM employee WHERE id=?")
		if err != nil {
			http.Error(w, "You have an error in your SQL syntax", 502)
			return
		}
		deleteRecords.Exec(emp)
		http.Redirect(w, r, "/", 301)
	}
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)	       // start page
	http.HandleFunc("/show", Show)     // showing records from database
	http.HandleFunc("/new", New)       // Create new records
	http.HandleFunc("/insert", Insert) // Insert new records
	http.HandleFunc("/edit", Edit)     // Edit details of specific records
	http.HandleFunc("/update", Update) // update details of specifice records
	http.HandleFunc("/delete", Delete) // delete records from database
	http.ListenAndServe(":8080", nil) // set port for server
}