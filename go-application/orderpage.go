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

//This for get all details of packages from database
type Details struct {
	Price   string
	Name    string
	Storage string
}

//This for get all users from database
type User struct {
	Id       int
	Username string
	Password string
	Email    string
}

//Connection of mysql database
func dbConn() (*sql.DB, error) {
	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	return db, err
}

//Global variable for get selected price from order_page
var price = ""

//Show all packages which have
func Index(w http.ResponseWriter, r *http.Request) {
	db, err := dbConn()
	defer db.Close()
	if err != nil {
		http.Error(w, "Database Connection Error!", 500)
		return
	}else {
		selectAllData, err := db.Query("SELECT * FROM details ORDER BY price desc")
		if err != nil {
			http.Error(w, "You have an error in your SQL syntax", 500)
			return
		}
		emp := Details{}
		res := []Details{}
		for selectAllData.Next() {
			var price, name, storage string
			err = selectAllData.Scan(&price, &name, &storage)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return 
			}
			emp.Price = price
			emp.Name = name
			emp.Storage = storage
			res = append(res, emp)
		}
		t, _ := template.ParseFiles("order_page.html")
		t.Execute(w, res)
	}
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
	name := r.FormValue("username")
	pass := r.FormValue("password")
	redirectTarget := "/details"
	db, err := dbConn()
	defer db.Close()
	if err != nil {
		http.Error(w, "Database Connection Error!", 500)
		return
	}else {
		rows, err := db.Query("SELECT COUNT(*) FROM userprofile where username=? and password=? ", name, pass)
		if err != nil {
			http.Error(w, "You have an error in your SQL syntax", 502)
			return
		}
		var count int
		for rows.Next() {
			if err := rows.Scan(&count); err != nil {
				http.Error(w, err.Error(), 502)
			}
		}
		if count <= 0 {
			redirectTarget = "/login_page"
		}
		http.Redirect(w, r, redirectTarget, 302)
	}
}

//New Entry template
func newEntry(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("register_form.html")
	t.Execute(w, nil)
}

//Add new user to database
func adduser(w http.ResponseWriter, r *http.Request) {
	db, err := dbConn()
	defer db.Close()
	if err != nil {
		http.Error(w, "Database Connection Error!", 500)
		return
	}else {
		redirectTarget := "/temp"
		if r.Method == "POST" {
			uname := r.FormValue("uname")
			pwd := r.FormValue("pwd")
			email := r.FormValue("email")
			insertRecord, err := db.Prepare("INSERT INTO userprofile(username, password, email) VALUES(?,?,?)")
			if err != nil {
				http.Error(w, "Sql syntax errorYou have an error in your SQL syntax", 500)
				return
			}
			res, err := insertRecord.Exec(uname, pwd, email)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			log.Println("Inserted records: Name: " + uname + "| Email: " + email)
			fmt.Printf("\nres: %v", res)
		}else {
			http.Error(w,"Method Not Found!", 500)
		}
		log.Printf("redirect to: %s", redirectTarget)
		http.Redirect(w, r, redirectTarget, 301)
	}
<<<<<<< HEAD
=======
	defer db.Close()
	http.Redirect(w,r,redirectTarget,302)
>>>>>>> d5583a05395668ee553fbb1e9d7e0a79d4c55448

}

//details of selected price from oreder page and display on details page
func orderDetails(w http.ResponseWriter, r *http.Request) {
	db, err := dbConn()
	defer db.Close()
	if err != nil {
		http.Error(w, "Database Connection Error!", 500)
		return
	}else {
		nPrice := price
		records, err := db.Query("SELECT * FROM details WHERE price=?", nPrice)
		if err != nil {
			http.Error(w, "You have an error in your SQL syntax", 500)
			return
		}
		emp := Details{}
		for records.Next() {
			var price, name, storage string
			err = records.Scan(&price, &name, &storage)
			if err != nil {
				http.Error(w, "Something wrong in data", 500)
				return 
			}
			emp.Price = price
			emp.Name = name
			emp.Storage = storage
		}
		t, _ := template.ParseFiles("details_page.html")
		t.Execute(w, emp)
	}
}


func last_page(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("last_page.html")
	t.Execute(w, r)
}

//Main function for handle all functions
func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)               //start page
	http.HandleFunc("/login", loginHandler)   //Login authentication
	http.HandleFunc("/login_page", loginPage) //Login page template
	http.HandleFunc("/temp", logTemp)         //redirect Login template after adduser
	http.HandleFunc("/new", newEntry)         //For new entry template
	http.HandleFunc("/insert", adduser)       //Add user function
	http.HandleFunc("/details", orderDetails) //details of orders
	http.HandleFunc("/last", last_page)       //last page
	http.ListenAndServe(":8080", nil)
}
