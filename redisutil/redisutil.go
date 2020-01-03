package redisutil

import (
	"os"
	"strconv"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis"
)

type RedisHashConn struct {
	client              *redis.Client
	locker              *redislock.Client
	data_ttl            int
	timeout_threshold   int
	persisting_filepath string
}

// TODO: return err
func NewRedisHashConn(address string, password string, database int, maxretries int, data_ttl int, persisting_filepath string) (RedisHashConn, error) {
	var new_redis_conn RedisHashConn

	new_client := redis.NewClient(&redis.Options{
		Addr:            address,
		Password:        password,
		DB:              database,
		MaxRetries:      maxretries,
		MinRetryBackoff: 500,
		MaxRetryBackoff: 1000,
	})

	new_redis_conn.client = new_client

	err := new_client.Ping().Err()
	if err != nil {
		return new_redis_conn, err
	}

	new_redis_conn.locker = redislock.New(new_client)
	new_redis_conn.data_ttl = data_ttl
	new_redis_conn.persisting_filepath = persisting_filepath

	return new_redis_conn, nil
}

func (db RedisHashConn) GetUser(username string) (string, error) {
	value, err := db.client.HGet("users", username).Result()

	return value, err
}

func (db RedisHashConn) CreateUser(username string, user_json_string string) error {

	lock, _ := db.locker.Obtain(username, 300*time.Second, nil)
	defer lock.Release()
	// Critical path here, set user, and timestamp of modification
	err := db.client.HSet("users", username, user_json_string).Err()
	time_of_modification_string := strconv.FormatInt(time.Now().UnixNano(), 10)
	db.client.HSet("modified_user_time", username, time_of_modification_string)

	go db.expire(username, time_of_modification_string)

	return err
}

func (db RedisHashConn) DeleteUser(user string) error {
	_, err := db.GetUser(user)
	if err == nil {
		db.client.HDel("users", user)
	}

	return err
}

func (db RedisHashConn) expire(username string, time_of_modification_string string) {

	time.Sleep(time.Duration(db.data_ttl) * time.Millisecond)

	lock, _ := db.locker.Obtain(username, 300*time.Second, nil)
	defer lock.Release()

	value, _ := db.client.HGet("modified_user_time", username).Result()

	if value == time_of_modification_string {
		user_data, _ := db.GetUser(username)
		db.client.HDel("modified_user_time", username)
		db.DeleteUser(username)

		// Todo: logging
		file, _ := os.Create(db.persisting_filepath + "/" + username + "-" + time_of_modification_string + ".json")
		file.WriteString(user_data)
		file.Close()
	}
}
