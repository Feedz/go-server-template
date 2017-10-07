package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"fmt"
)

var tpl *template.Template

const (
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "myDbName"
)

func init() {
	/* DB CONNECT */
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbInfo)
	checkErr(err)
	defer db.Close()
	/* END DB CONNECT */

	getUserInfo(*db)

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)

	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {


	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "login.gohtml", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getUserInfo(db sql.DB) {
	/* ONLY FOR TESTING */
	rows, err := db.Query("SELECT uid,email,username,age FROM users")
	checkErr(err)

	for rows.Next() {
		var uid int
		var email string
		var username string
		var age int
		err = rows.Scan(&uid, &email, &username, &age)
		checkErr(err)
		fmt.Println("uid | email | username | age ")
		fmt.Printf("%3v | %8v | %6v | %6v\n", uid, email, username, age)
	}
}
