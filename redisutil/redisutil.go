package redisutil

import (
	"github.com/go-redis/redis"
)

type RedisHashConn struct {
	client *redis.Client
}

// TODO: return err
func NewRedisHashConn(address string, password string, database int, maxretries int) (RedisHashConn, error) {
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

	return new_redis_conn, nil
}

func (db RedisHashConn) GetUser(username string) (string, error) {
	value, err := db.client.HGet("users", username).Result()

	return value, err
}

func (db RedisHashConn) CreateUser(username string, user_json_string string) error {
	err := db.client.HSet("users", username, user_json_string).Err()

	return err
}

func (db RedisHashConn) DeleteUser(user string) error {
	_, err := db.GetUser(user)
	if err == nil {
		db.client.HDel("users", user)
	}

	return err
}
