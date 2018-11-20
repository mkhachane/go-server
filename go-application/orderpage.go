package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

import _ "github.com/go-sql-driver/mysql"

type Details struct {
	Price string
	Name  string
	Storage string
}

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:root@/packages")
	if err != nil {
		panic(err.Error())
	}
	return db
}


func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selectAllData, err := db.Query("SELECT * FROM details ORDER BY price desc")
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	emp := Details{}
	res := []Details{}
	for selectAllData.Next() {
		var price, name, storage string
		err = selectAllData.Scan(&price, &name, &storage)
		if err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		emp.Price = price
		emp.Name = name
		emp.Storage = storage
		res = append(res, emp)
	}
	t, _ := template.ParseFiles("order_page.html")
	t.Execute(w, res)
	defer db.Close()
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login_page.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}

func details(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	//nPrice := r.URL.Query().Get("price")
	selectAllData, err := db.Query("SELECT * FROM details ORDER BY price desc")
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	emp := Details{}
	for selectAllData.Next() {
		var price, name, storage string
		err = selectAllData.Scan(&price, &name, &storage)
		if err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		emp.Price = price
		emp.Name = name
		emp.Storage = storage
	}
	t, _ := template.ParseFiles("details_page.html")
	t.Execute(w, emp)
	defer db.Close()
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)        // start page
	http.HandleFunc("/login", login)
	http.HandleFunc("/details", details)
	http.ListenAndServe(":8080", nil)  // set port for server
}