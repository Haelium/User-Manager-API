package redisutil

import (
	"strconv"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis"
)

type RedisHashConn struct {
	client   *redis.Client
	locker   *redislock.Client
	data_ttl int
}

// TODO: return err
func NewRedisHashConn(address string, password string, database int, maxretries int, data_ttl int) (RedisHashConn, error) {
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

	return new_redis_conn, nil
}

func (db RedisHashConn) GetUser(username string) (string, error) {
	value, err := db.client.HGet("users", username).Result()

	return value, err
}

func (db RedisHashConn) CreateUser(username string, user_json_string string) error {
	err := db.client.HSet("users", username, user_json_string).Err()
	go db.expire(username)

	return err
}

func (db RedisHashConn) DeleteUser(user string) error {
	_, err := db.GetUser(user)
	if err == nil {
		db.client.HDel("users", user)
	}

	return err
}

func (db RedisHashConn) expire(username string) {
	time_of_modification := time.Now().UnixNano()
	time_of_modification_string := strconv.FormatInt(time_of_modification, 10)

	db.client.HSet("modified_user_time", username, time_of_modification_string)

	time.Sleep(time.Duration(db.data_ttl))

	value, _ := db.client.HGet("modified_user_time", username).Result()

	if value == time_of_modification_string {
		db.client.HDel("modified_user_time", username)
	}
}
