package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// Using a mock object for database

type UserMap struct {
	users map[string]string
}

func NewUserMap() UserMap {
	var newUserMap UserMap
	newUserMap.users = make(map[string]string)

	return newUserMap
}

func (db UserMap) EditUser(username string, user_json_string string) error {
	db.users[username] = user_json_string
	return nil
}

func (db UserMap) CreateUser(username string, user_json_string string) error {
	db.users[username] = user_json_string
	return nil
}

func (db UserMap) GetUser(username string) (string, error) {
	user, exists := db.users[username]
	if exists == false {
		return "", errors.New("User not found")
	} else {
		return user, nil
	}
}

func (db UserMap) DeleteUser(username string) error {
	_, err := db.GetUser(username)
	delete(db.users, username)

	return err
}

func Router(user_db UserMap) *mux.Router {
	handler := NewHandler(user_db)

	router := mux.NewRouter()

	router.HandleFunc("/user", handler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/", handler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{username}", handler.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{username}/", handler.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{username}", handler.DeleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/user/{username}/", handler.DeleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/user/{username}", handler.EditUser).Methods(http.MethodPut)
	router.HandleFunc("/user/{username}/", handler.EditUser).Methods(http.MethodPut)

	return router
}
func Test_Create(t *testing.T) {
	user_db := NewUserMap()

	test_user := []byte(`{"username": "billy2000", "fullname": "Bob Bobson", "email": "Bob@bobmail.bob", "address": {"name": "Bob", "Line 1": "44 Bobstreet", "region": "Bobville", "country": "Bobland"}}`)

	request, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(test_user))
	response := httptest.NewRecorder()
	Router(user_db).ServeHTTP(response, request)

	if (response.Code) != 201 {
		t.Logf("Unexpected response code: %d", response.Code)
		t.Fail()
	}

	if response.Body.String() != "Created user" {
		t.Logf("Unexpected response body: %s", response.Body)
		t.Fail()
	}

	if user_db.users["billy2000"] != string(test_user) {
		t.Logf("Input corrupted")
		t.Fail()
	}
}

func Test_Get(t *testing.T) {
	user_db := NewUserMap()

	test_user := []byte(`{"username": "billy2000", "fullname": "Bob Bobson", "email": "Bob@bobmail.bob", "address": {"name": "Bob", "Line 1": "44 Bobstreet", "region": "Bobville", "country": "Bobland"}}`)

	user_db.users["billy2000"] = string(test_user)

	// Test get here
	request, _ := http.NewRequest("GET", "/user/billy2000/", nil)
	response := httptest.NewRecorder()
	Router(user_db).ServeHTTP(response, request)

	if response.Code != 200 {
		t.Logf("Expected: %d\nGot %d\n", 200, response.Code)
		t.Fail()
	}

	if response.Body.String() != string(test_user) {
		t.Logf("Expected: %s\nGot: %s\n", test_user, response.Body)
		t.Fail()
	}
}

func Test_Delete(t *testing.T) {
	user_db := NewUserMap()

	test_user := []byte(`{"username": "billy2000", "fullname": "Bob Bobson", "email": "Bob@bobmail.bob", "address": {"name": "Bob", "Line 1": "44 Bobstreet", "region": "Bobville", "country": "Bobland"}}`)

	user_db.users["billy2000"] = string(test_user)

	// Test get here
	request, _ := http.NewRequest("DELETE", "/user/billy2000/", nil)
	response := httptest.NewRecorder()
	Router(user_db).ServeHTTP(response, request)

	_, exists := user_db.users["billy2000"]

	if exists != false {
		t.Logf("User was not deleted")
		t.Fail()
	}
}

func Test_Edit(t *testing.T) {
	user_db := NewUserMap()

	test_user := []byte(`{"username":"billy2000","fullname":"Bob Bobson","email":"Bob@bobmail.bob","address":{"name":"Bob","Line 1":"44 Bobstreet","region":"Bobville","country":"Bobland"}}`)
	test_user_mod := []byte(`{"username":"billy2000","fullname":"Robert Newname","email":"Bob@bobmail.bob","address":{"name":"Bob","line 1":"44 Bobstreet","region":"Bobville","country":"Bobland"}}`)

	user_db.users["billy2000"] = string(test_user)

	// Test get here
	request, _ := http.NewRequest("PUT", "/user/billy2000/", bytes.NewBuffer(test_user_mod))
	response := httptest.NewRecorder()
	Router(user_db).ServeHTTP(response, request)

	if response.Body.String() != "Created user" {
		t.Logf("Expected: %s\nGot: %s\n", "Created user", response.Body.String())
		t.Fail()
	}

	if user_db.users["billy2000"] != string(test_user_mod) {
		t.Logf("User was not updated correctly\nExpected: %s\nGot: %s\n", user_db.users["billy2000"], string(test_user_mod))
		t.Fail()
	}
}
