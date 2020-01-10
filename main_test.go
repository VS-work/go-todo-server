package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize("./db/todos_test.db")

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM todos")
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS todos
(
	priority INT NOT NULL,
	content VARCHAR(50) NOT NULL,
	completed INT NOT NULL
)`

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addTodos(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO todos(content, completed, priority) VALUES('%s', 0, 0)", ("Todo " + strconv.Itoa(i+1)))
		_, err := a.DB.Exec(statement)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/todos", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentUser(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/todo/111", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string

	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Todo not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Todo not found'. Got '%s'", m["error"])
	}
}

func TestCreateTodo(t *testing.T) {
	clearTable()

	payload := []byte(`{"content": "test todo"}`)

	req, _ := http.NewRequest("POST", "/todo", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["content"] != "test todo" {
		t.Errorf("Expected todo content to be 'test todo'. Got '%v'", m["content"])
	}

	if m["priority"] != 1.0 {
		t.Errorf("Expected todo Priority to be '1'. Got '%v'", m["priority"])
	}

	if m["rowid"] != 1.0 {
		t.Errorf("Expected todo ID to be '1'. Got '%v'", m["rowid"])
	}
}

func TestGetTodo(t *testing.T) {
	clearTable()

	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateTodo(t *testing.T) {
	clearTable()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)

	var originalTodo map[string]interface{}

	json.Unmarshal(response.Body.Bytes(), &originalTodo)

	payload := []byte(`{"content": "test todo - updated content"}`)

	req, _ = http.NewRequest("PUT", "/todo/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalTodo["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalTodo["id"], m["id"])
	}

	if m["content"] == originalTodo["content"] {
		t.Errorf("Expected the content to change from '%v' to '%v'. Got '%v'", originalTodo["content"], m["content"], m["content"])
	}
}

func TestDeleteTodo(t *testing.T) {
	clearTable()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/todo/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	req, _ = http.NewRequest("GET", "/todo/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
