package redisutil

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
)

var valid_users = map[string]string{
	"bobman12": `{"username": "bobman12", "email": "bob@bobmail.com"}`,
	"jcdenton": `{"username": "jcdenton", "email": "jc@unatco.org"}`,
	"herpderp": `{"username": "herpderp", "email": "herp@derp.io"}`,
}

func Test_GetUser(t *testing.T) {
	// Set up minikube for testing, fail if not working
	miniredis_socket, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer miniredis_socket.Close()

	// Set users in to test on
	for key, val := range valid_users {
		miniredis_socket.HSet("users", key, val)
	}
	// Start client
	redis_client, _ := NewRedisHashConn(miniredis_socket.Addr(), "", 0, 5, 32, ".")

	for key, expected_val := range valid_users {
		actual_val, err := redis_client.GetUser(key)

		if err != nil {
			t.Logf("err: %s", err)
			t.Fail()
		}

		if expected_val != actual_val {
			t.Logf("Expected:\t %s \nGot:\t %s\n", expected_val, actual_val)
			t.Fail()
		}
	}
}

func Test_CreateUser(t *testing.T) {
	// Set up minikube for testing, fail if not working
	miniredis_socket, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer miniredis_socket.Close()

	redis_client, _ := NewRedisHashConn(miniredis_socket.Addr(), "", 0, 5, 32, ".")

	for username, userdata := range valid_users {
		redis_client.CreateUser(username, userdata)
	}

	for key, expected_val := range valid_users {
		actual_val, err := redis_client.GetUser(key)

		if err != nil {
			t.Logf("err: %s", err)
			t.Fail()
		}

		if expected_val != actual_val {
			t.Logf("Expected:\t %s \nGot:\t %s\n", expected_val, actual_val)
			t.Fail()
		}
	}
}

func Test_expire(t *testing.T) {
	miniredis_socket, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer miniredis_socket.Close()

	redis_client, _ := NewRedisHashConn(miniredis_socket.Addr(), "", 0, 5, 5, ".")
	redis_client.CreateUser("bob_should_expire", "junk data")
	time.Sleep(7 * time.Second)

	returned_value, err := redis_client.GetUser("bob_should_expire")
	if returned_value != "" {
		t.Logf("Data is not being expired: %s", returned_value)
		t.Fail()
	}

	redis_client.CreateUser("bob_should_not_expire", "data_still_there")
	time.Sleep(4 * time.Second)

	returned_value, err = redis_client.GetUser("bob_should_not_expire")
	if returned_value != "data_still_there" {
		t.Logf("Data is expired too early: %s", returned_value)
		t.Fail()
	}

}

/*
func Test_DeleteUser(t *testing.T) {
	// Set up minikube for testing, fail if not working
	miniredis_socket, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer miniredis_socket.Close()

	redis_client := NewRedisHashConn(miniredis_socket.Addr(), "", 0)

	for _, val := range valid_users {
		redis_client.CreateUser(val)
	}

	for key, _ := range valid_users {

		, err := redis_client.DeleteUser(key)

		if err != nil {
			t.Logf("err: %s", err)
		}

		returned_user, err := redis_client.GetUser(key)

		fmt.Println(returned_user, err)

		if expected_val != actual_val {
			t.Logf("Expected:\t %s \nGot:\t %s\n", expected_val, actual_val)
		}

	}
}
*/
