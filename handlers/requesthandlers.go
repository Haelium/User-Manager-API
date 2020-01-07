package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Haelium/BloomReachTest/validation"
)

/*
TODO:
config: {
	redis_timeout int
	batch_size int
	redis_backend string
	listen_port int
}
*/

/*
Schema:

user: {
	username: unique
	name
	email
	shipping address
}

POST /user/							- Create user 				- Takes json struct
GET /user/{username}				- Gets user 				- Returns json struct
DELETE /user/{username}				- Deletes user 				- Returns json struct
PUT /user/{username}				- Updates user				- Takes json struct

*/

type DatabaseInterface interface {
	SetUser(string, string) error
	GetUser(string) (string, error)
	DeleteUser(string) error
}

type RequestHandler struct {
	db DatabaseInterface
	// Log level?
	// Log path?
}

func NewHandler(db DatabaseInterface) RequestHandler {
	var handler RequestHandler
	handler.db = db

	return handler
}

func responseErrorBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err)))
	return
}

func responseErrorNotFound(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err)))
	return
}

func responseErrorForbidden(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err)))
	return
}

func (handler RequestHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	username := pathParams["username"]

	user_json_string, err := handler.db.GetUser(username)
	if err != nil {
		responseErrorNotFound(w, err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseErrorBadRequest(w, err)
		return
	}

	new_user_body_json, err := validation.ModifyUser(user_json_string, string(body))
	if err != nil {
		responseErrorBadRequest(w, err)
		return
	}

	err = handler.db.SetUser(username, string(new_user_body_json))
	if err != nil {
		responseErrorBadRequest(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Created user"))
}

func (handler RequestHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseErrorBadRequest(w, err)
		return
	}

	username, err := validation.ValidateUser(string(body))
	if err != nil {
		responseErrorBadRequest(w, err)
		return
	}

	existing_user, _ := handler.db.GetUser(username)
	if existing_user != "" {
		responseErrorForbidden(w, errors.New("User already exists"))
		return
	}

	err = handler.db.SetUser(username, string(body))
	if err != nil {
		responseErrorBadRequest(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Created user"))
}

func (handler RequestHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	username := pathParams["username"]
	user_json_string, err := handler.db.GetUser(username)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "User not found"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(user_json_string))
	}
}

func (handler RequestHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	username := pathParams["username"]
	err := handler.db.DeleteUser(username)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "User not found"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"deleted\":\"" + username + "\"}"))
	}
}
