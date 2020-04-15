// File: main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pborman/getopt"
	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	ID       string
	Username string
	Password string
}

var db *sql.DB

func init() { //wird immer ausgefuehrt

	db, _ = sql.Open("sqlite3", "./user.db") //Datenbank-Connect

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS people (ID TEXT PRIMARY KEY, Username TEXT, Password TEXT)") //Create Table
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec() // Create Table ausfuehren

	//fuer Testing
	statement, err = db.Prepare("INSERT INTO people (ID,Username, Password) VALUES (?, ?, ?)") //Insert-Statement
	if err != nil {
		log.Fatal(err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	statement.Exec("user@tgm.ac.at", "user", hashedPassword) //Insert ausfuehren
	printUsers()
}

func personCreate(w http.ResponseWriter, r *http.Request) {
	var p Person //Person struct als variable

	err := json.NewDecoder(r.Body).Decode(&p) //den json body dem struct p zuweisen
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if (p.Username == "") || (p.Password == "") || (p.ID == "") {
		http.Error(w, "Fuer die Registrierung muessen alle Daten (ID, Username, Password) gesendet werden.", 400)
		return
	}

	if (strings.Contains(p.ID, `'`)) || (strings.Contains(p.ID, ";")) || (strings.Contains(p.ID, "--")) || (strings.Contains(p.Username, `'`)) || (strings.Contains(p.Username, ";")) || (strings.Contains(p.Username, "--")) || (strings.Contains(p.Password, `'`)) || (strings.Contains(p.Password, ";")) || (strings.Contains(p.Password, "--")) {
		http.Error(w, "No SQL-Injection pls :(", 400)
		return
	}

	if !(strings.Contains(p.ID, `@`)) {
		http.Error(w, "ID should be an e-mail", 400)
		return
	}

	rows, err := db.Query("SELECT * FROM people WHERE ID = '" + p.ID + "'") //Schauen ob der Name schon in der Datenbank existiert
	if err != nil {
		log.Fatal(err)
		return
	}
	var exist bool
	exist = false
	var ID string
	var Username string
	var Password string
	for rows.Next() {
		rows.Scan(&ID, &Username, &Password)
		if ID == p.ID {
			exist = true
		}
	}
	if exist {
		http.Error(w, "Failure User "+p.ID+" already exists", 400)
		return
	} else {
		statement, err := db.Prepare("INSERT INTO people (ID, Username,Password) VALUES (?, ?, ?)") //Wenn der Name nicht existiert, wird ein neuer User erstellt.
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost) //save as hash
		statement.Exec(p.ID, p.Username, hashedPassword)
	}
	fmt.Fprintf(w, "Successfull User "+p.ID+" has been created")
}

func personLogin(w http.ResponseWriter, r *http.Request) {
	var p Person

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if (strings.Contains(p.ID, `'`)) || (strings.Contains(p.ID, ";")) || (strings.Contains(p.ID, "--")) || (strings.Contains(p.Username, `'`)) || (strings.Contains(p.Username, ";")) || (strings.Contains(p.Username, "--")) || (strings.Contains(p.Password, `'`)) || (strings.Contains(p.Password, ";")) || (strings.Contains(p.Password, "--")) {
		http.Error(w, "No SQL-Injection pls", 400)
		return
	}

	var (
		dbid       string
		dbpassword string
	)
	query := `SELECT ID,Password FROM people WHERE ID = ?`
	err = db.QueryRow(query, p.ID).Scan(&dbid, &dbpassword)
	if err != nil {
		http.Error(w, "Failure User doesn't exist", 400)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(p.Password)) //compare hash with password

	if err != nil {
		http.Error(w, "Failure Password is wrong", 400)
		return
	}
	fmt.Fprintf(w, "Successfull Willkommen")

}

func printUsers() {
	rows, err := db.Query("SELECT * FROM people")
	if err != nil {
		log.Fatal(err)
	}

	var id string
	var username string
	var Password string

	for rows.Next() {
		rows.Scan(&id, &username, &Password)
		fmt.Println(id + " " + username + " " + Password)
	}
}

func main() {

	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)

	port := getopt.StringLong("port", 'p', "", "Portnumber")
	ip := getopt.StringLong("ip", 'i', "", "IP Adress")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	if *port == "" {
		println("Server running on " + *ip + ":8080")
		log.Fatal(http.ListenAndServe(*ip+":8080", nil))
	} else {
		println("Server running on " + *ip + ":" + *port)
		log.Fatal(http.ListenAndServe(*ip+":"+*port, nil))
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		personLogin(w, r)
	default:
		w.Header().Set("Allow", "POST")
		http.Error(w, "Available methods for /login are: POST", http.StatusMethodNotAllowed)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		personCreate(w, r)
	default:
		w.Header().Set("Allow", "POST")
		http.Error(w, "Available methods for /register are: POST", http.StatusMethodNotAllowed)
	}
}
