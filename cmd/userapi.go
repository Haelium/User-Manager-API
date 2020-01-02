package main

import (
	//	"encoding/json"

	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Haelium/BloomReachTest/handlers"
	"github.com/Haelium/BloomReachTest/redisutil"
)

// TDOD: input args & config struct

func main() {
	user_db, err := redisutil.NewRedisHashConn("localhost:6379", "", 0)
	if err != nil {
		log.Panicf("Exit: %s", err)
	}

	handler := handlers.NewHandler(user_db)

	router := mux.NewRouter()

	router.HandleFunc("/user", handler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{username}", handler.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{username}", handler.DeleteUser).Methods(http.MethodDelete)
	//	router.PathPrefix("/").Handler(catchAllHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
