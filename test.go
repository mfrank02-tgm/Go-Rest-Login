package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func init() {
	os.Remove("./user.db") //zerstoert Persistierung ist nur fuers Testing
}

func TestLoginSuccess(t *testing.T) {
	var jsonStr = []byte(`{
	"ID":"user@tgm.ac.at",
	"Password":"password"
	}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personLogin)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 200 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 200)
	}

	expected := `Successfull`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestLoginWithSomething(t *testing.T) {
	var jsonStr = []byte(`{
	"ndsfds":"user",
	"psdfsd":"password"
	}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personLogin)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 400)
	}

	expected := `Failure`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestLoginWithNumbers(t *testing.T) {
	var jsonStr = []byte(`{
	"ID":42,
	"Password":666
	}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personLogin)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 400)
	}

	expected := `json: cannot unmarshal number`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	var jsonStr = []byte(`{
	"ID":"user@tgm.ac.at",
	"Password":"falsch"
	}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personLogin)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 400)
	}

	expected := `Failure`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestLoginWrongUsername(t *testing.T) {
	var jsonStr = []byte(`{
	"ID":"falsch@tgm.ac.at",
	"Password":"password"
	}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personLogin)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 400)
	}

	expected := `Failure`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestLoginWithoutJson(t *testing.T) {
	var jsonStr = []byte(``)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personLogin)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 400)
	}

	expected := `EOF`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestLoginWithUserThatDoesntExist(t *testing.T) {
	var jsonStr = []byte(`{
	"ID":"falsch",
	"Password":"falsch"
	}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personLogin)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 400)
	}

	expected := `Failure User doesn't exist`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPersonCreate(t *testing.T) {

	var jsonStr = []byte(`{
	"ID":"user2@tgm.ac.at",
	"Username": "user2",
	"Password":"password2"
	}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personCreate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `Successfull`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPersonCreateDuplicate(t *testing.T) {
	var jsonStr = []byte(`{
	"ID":"user@tgm.ac.at",
	"Username": "user",
	"Password":"password"
	}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personCreate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 400)
	}

	expected := `Failure`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPersonCreateWithNumber(t *testing.T) {
	var jsonStr = []byte(`{
	"ID":666,
	"Username":1337,
	"Password":42
	}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personCreate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 400)
	}

	expected := `json: cannot unmarshal number`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPersonCreateSQLInjection(t *testing.T) {
	var jsonStr = []byte(`{
	"ID":"test'--",
	"Username":"test",
	"Password":"test"
	}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personCreate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 400)
	}

	expected := `No SQL-Injection`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPersonCreateWithSomething(t *testing.T) {
	var jsonStr = []byte(`{
	"gdsg":666,
	"sdg": 42
	}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personCreate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 200)
	}

	expected := `Fuer die Registrierung muessen alle Daten (ID, Username, Password) gesendet werden.`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPersonCreateWithoutID(t *testing.T) {
	var jsonStr = []byte(`{
	"ID":"",
	"Username": "user",
	"Password":"password"
	}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personCreate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 400)
	}

	expected := `Fuer die Registrierung muessen alle Daten (ID, Username, Password) gesendet werden.`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPersonCreateWithoutUsername(t *testing.T) {
	var jsonStr = []byte(`{
	"ID":"user@tgm.ac.at",
	"Username": "",
	"Password":"password"
	}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personCreate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 400)
	}

	expected := `Fuer die Registrierung muessen alle Daten (ID, Username, Password) gesendet werden.`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPersonCreateWithoutPassword(t *testing.T) {
	var jsonStr = []byte(`{
	"ID":"user@tgm.ac.at",
	"Username": "user",
	"Password": ""
	}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personCreate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 400)
	}

	expected := `Fuer die Registrierung muessen alle Daten (ID, Username, Password) gesendet werden.`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
func TestPersonCreateWithoutAt(t *testing.T) {
	var jsonStr = []byte(`{
	"ID":"user.tgm.ac.at",
	"Username": "user",
	"Password": "password"
	}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personCreate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 400)
	}

	expected := `ID should be an e-mail`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPersonCreateWithEmptyJson(t *testing.T) {
	var jsonStr = []byte(`{}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(personCreate)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 400)
	}

	expected := `Fuer die Registrierung muessen alle Daten (ID, Username, Password) gesendet werden.`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
