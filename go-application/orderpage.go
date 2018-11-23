package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

import _ "github.com/go-sql-driver/mysql"

//This for get all details of packages from database
type Details struct {
	Price string
	Name  string
	Storage string
}

//This for get all users from database
type User struct {
	Id    int
	Username  string
	Password string
	Email string
}

//Connection of mysql database
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "packages"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

//Global variable for get selected price from order_page
var price = " "

//Show all packages which have
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

//This for get price of package which selected when order
func loginPage(w http.ResponseWriter, r *http.Request) {
	price = r.URL.Query().Get("price")
	t, _ := template.ParseFiles("login_page.html")
	t.Execute(w, nil)

}

//This for login_page template for update user
func logTemp(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("login_page.html")
	t.Execute(w, nil)

}

//login authentication check user 
func loginHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	redirectTarget := "/details"
	db := dbConn()
	rows, err := db.Query("SELECT COUNT(*) FROM userprofile where username=? and password=? ", name,pass)
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
	}
	if count <= 0 {
		redirectTarget = "/login_page"
	}
	http.Redirect(w, r, redirectTarget, 302)
}

//New Entry template
func newEntry(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("register_form.html")
	t.Execute(w, nil)
}


//Add new user to database
func adduser(w http.ResponseWriter, r *http.Request)  {
	db := dbConn()
	redirectTarget := "/temp"
	if r.Method == "POST" {
		uname := r.FormValue("uname")
		pwd   := r.FormValue("pwd")
		email := r.FormValue("email")
		insertRecord, err  := db.Prepare("INSERT INTO userprofile(username, password, email) VALUES(?,?,?)")
		if err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		res, err :=insertRecord.Exec(uname, pwd, email)
		if err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		log.Println("Inserted records: Name: " + uname + "| Email: "+ email)
		fmt.Printf("\nres: %v", res)
	}
	defer db.Close()
	http.Redirect(w,r,redirectTarget,302)

}


//details of selected price from oreder page and display on details page
func orderDetails(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nPrice := price
	records, err := db.Query("SELECT * FROM details WHERE price=?", nPrice)
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	emp := Details{}
	for records.Next() {
		var price, name, storage string
		err = records.Scan(&price, &name, &storage)
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

func last_page(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("last_page.html")
	t.Execute(w,r)
}


//Main function for handle all functions 
func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)               //start page
	http.HandleFunc("/login", loginHandler)   //Login authentication
	http.HandleFunc("/login_page", loginPage) //Login page template
	http.HandleFunc("/temp", logTemp)		  //redirect Login template after adduser
	http.HandleFunc("/new", newEntry)         //For new entry template
	http.HandleFunc("/insert", adduser)       //Add user function
	http.HandleFunc("/details", orderDetails) //details of orders
	http.HandleFunc("/last", last_page)       //last page 
	http.ListenAndServe(":8080", nil)
}
