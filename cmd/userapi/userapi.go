package main

import (
	//	"encoding/json"

	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Haelium/BloomReachTest/handlers"
	"github.com/Haelium/BloomReachTest/redisutil"
)

func main() {
	redisAddrPtr := flag.String("redis_address", "localhost", "Redis server address")
	redisPortPtr := flag.String("redis_port", "6379", "Redis server port")
	redisPassPtr := flag.String("redis_password", "", "Redis server password")
	redisDBIndexPtr := flag.Int("redis_db_index", 0, "Redis database index")
	redisMaxRetries := flag.Int("redis_max_retries", 5, "Number of times redis will retry a command")

	appListenPortPtr := flag.String("listen_port", "8080", "Port which service listens on")
	appDataTTLSeconds := flag.Int("data_ttl", 60, "time before data expires (seconds)")
	appDataPersistPath := flag.String("persist_path", "/opt/userapidata/", "path to directory to save data")

	flag.Parse()

	user_db, err := redisutil.NewRedisHashConn(
		(*redisAddrPtr)+":"+(*redisPortPtr), *redisPassPtr, *redisDBIndexPtr, *redisMaxRetries, *appDataTTLSeconds, *appDataPersistPath,
	)
	if err != nil {
		log.Panicf("Exit: %s", err)
	}

	handler := handlers.NewHandler(user_db)

	router := mux.NewRouter()

	router.HandleFunc("/user", handler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{username}", handler.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{username}", handler.DeleteUser).Methods(http.MethodDelete)
	//	router.PathPrefix("/").Handler(catchAllHandler)

	log.Fatal(http.ListenAndServe(":"+(*appListenPortPtr), router))
}
