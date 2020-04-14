Markus Frank 5CHIT
# SEW EK REST Login + Registrierung

## Design und Beschreibung

Ich verwende fuer die EK die Programmiersprache golang. Da diese fuer Webapplikationen ideal ist, da es dafür in go kein Framework braucht. Zusätzlich kann man go compilieren im Gegensatz zu anderen Sprachen. Fuer den Webserver wird das Package http von go verwendet. Dann wurde eine einfache REST-Schnittstelle und eine SQLite Datenbank erstellt. Zuletzt wurden die Passwoerter verschluesselt und das Programm getestet.



## Implementierung

#### main.go

REST-Webserver in der `main()` definieren: 

HandleFunc weißt `/register` der Funktion personCreate zu.

HandleFunc weißt `/login` der Funktion personLogin zu.

Der Port ist 8080.

`http.NewServeMux()`

```go
func main() {
    http.HandleFunc("/register", personCreate)
    http.HandleFunc("/login", personLogin)
    err := http.ListenAndServe(":8080", mux)
    log.Fatal(err)
}
```

**struct**

Ein struct ist eine Art Interface oder Objekt in denen man Daten speichert. Diese eigen sich ideal zum Umwandeln in json oder speichern in einer Datenbank.

```go
type Person struct {
    ID string
    Username string
    Password  string
}
```

**init()**

Init-Funktion ist die Funktion die immer (auch beim Testing) ausgeführt wird. Sie eignet sich gut um fürs Initialisieren von Objekten, die gebraucht werden und um Datenbankverbindungen aufzubauen.

```go
var db *sql.DB
func init() {
    db, _ = sql.Open("sqlite3", "./user.db")
    statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS people (Name TEXT PRIMARY KEY, Password TEXT)")
    if err != nil {
        log.Fatal(err)
    }
    statement.Exec()
}
```

**Default User in init()**

Dieser User wird für die Test-Cases verwendet und wird daher in der init erstellt.

User = "user"

Password = hashed "password"

```go
statement, err = db.Prepare("INSERT INTO people (ID,Username, Password) VALUES (?, ?, ?)")
if err != nil {
	log.Fatal(err)
}
hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
statement.Exec("user@tgm.ac.at","user", hashedPassword )
printUsers()
```



**personCreate**

`json.NewDecoder(r.Body).Decode(&p)`holt den Request body und weist sie einem struct Objekt `var p Person`  zu. 

Dann noch die Daten in die Datenbank eintragen.

```go
statement, err = db.Prepare("INSERT INTO people (ID,Username, Password) VALUES (?,?,?)")
if err != nil {
	log.Fatal(err)
}
statement.Exec(p.ID,p.Username,p.Password )
```

mit `fmt.Fprintf(w, "Person: %+v", p)`antworte der Server in der Rest Schnittstelle.



**personLogin**

Das json holt sich diese Funktion genau wie die Funktion personCreate. Zusätzlich wird die Datenbank nach einer Person durchsucht, die das selbe Passwort hat.

```go
var (
    dbid  string
    dbpassword  string
)
query := `SELECT ID,Password FROM people WHERE ID = ?`
err = db.QueryRow(query, p.ID).Scan( &dbid, &dbpassword)
if err != nil {
    http.Error(w, "Failure User doesn't exist", 400)
    return
}
err = bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(p.Password))
fmt.Fprintf(w, "Successfull Willkommen") //Willkommensnachricht
```



Verschlüsselung:

```
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
```

**Passwort mit Hash vergleichen:**

p.password kommt von dem Json

dbpasswort kommt von der SQLite Database

```
err = bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(p.Password))
```




### Build tool Dependencies installieren

`go.mod`Datei erstellen und einen Modulnamen auswaehlen:

```
module Go-Rest-Login 
```

mit `go get [packagename]`die Packages in die `go.mod` hinzu generieren.

```
go 1.14

require (
	github.com/mattn/go-sqlite3 v2.0.1+incompatible
	github.com/pborman/getopt v0.0.0-20190409184431-ee0cd42419d3
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413
)
```

Sodass `go build` auch die Dependencies installiert.



### Testing

Json erstellen und Request erstellen

```
var jsonStr = []byte(`{
	"ID":"user@tgm.ac.at",
	"Password":"password"
}`)


req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
if err != nil {
	t.Fatal(err)
}
```

Check Status Code:

Bei positiven Ergebnis wird 200 erwartet und bei negativen 400.

```
if status := rr.Code; status != 200 {
	t.Errorf("handler returned wrong status code: got %v want %v",status, 200)
}
```

Check Output:

```
expected := `Successfull`
if !strings.Contains(rr.Body.String(),expected) {
	t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
}
```

### Integration Tests(travis-ci):

**.travis.yml**
```
language: go
go: "1.13.5"
```

### System Testing mit curl

Register:

```
$ curl -d '{"ID":"user@tgm.ac.at","Username":"user","Password":"password"}' -H "Content-Type: application/json" http://localhost:8080/register
Failure User user already exists
```

```
$ curl -d '{"ID":"mfrank02@tgm.ac.at","Username":"mfrank02","Password":"test"}' -H "Content-Type: application/json" http://localhost:8080/register
Successfull User mfrank02@tgm.ac.at has been created
```

Login

```
$ curl -d '{"ID":"user@tgm.ac.at"r","Password":"password"}' -H "Content-Type: application/json" http://localhost:8080/login
Successfull
```

```
$ curl -d '{"Name":666,"Password":42}' -H "Content-Type: application/json" http://localhost:8080/login
json: cannot unmarshal number into Go struct field Person.Name of type string
```




## Deployment

#### Requirements

* golang

Installation (git)

```bash
git clone https://github.com/mfrank02-tgm/Go-Rest-Login.git
```
Change Directory
```bash
cd Go-Rest-Login
```

Build + Dependencies installieren und ausfuehren

```
go build -o Go-Rest-Login
./Go-Rest-Login
```

simples ausführen

```bash
go run main.go
```

REST Testing

```bash
go test -v
```

Binary Ausfuehren
```bash
./Go-Rest-Login  #Unix and Unix-Like
Go-Rest-Login.exe  #Windoof
```

Binary Ausfuehren mit command-line-flags
```bash
./Go-Rest-Login --help
./Go-Rest-Login -p 8080 -i 127.0.0.1
./Go-Rest-Login --port=8080 --ip=127.0.0.1
```

Testen ob der Server laeuft:

```bash
curl -d '{"ID":"user@tgm.ac.at"r","Password":"password"}' -H "Content-Type: application/json" http://localhost:8080/login
```


**Troubleshooting**:

Wenn noch nicht alle Pakete (beim import zu sehen) installiert sind:

```
go get [packagename]
```

