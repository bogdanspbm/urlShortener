package utils

import (
	"errors"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-redis/redis/v8"
)

const (
	prefix = "bmadzhuga"
	name   = "main"
	key    = "bmadzhuga::main"
)

type Redis struct {
	Cluster string
	Client  *redis.Client
}

func (client *Redis) Connect() error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     client.Cluster,
		Password: "",
		DB:       0,
	})

	_, err := redisClient.Ping(redisClient.Context()).Result()

	if err != nil {
		return err
	}

	client.Client = redisClient
	return nil
}

func (client *Redis) Push(value string) error {
	if client.Client == nil {
		return errors.New("Empty redis client")
	}

	client.Client.RPush(client.Client.Context(), key, value)

	return nil
}

func (client *Redis) Put(field string, val string) error {
	if client.Client == nil {
		return errors.New("Empty redis client")
	}

	err := client.Client.HSet(client.Client.Context(), key, field, val)

	if err != nil {
		return errors.New("Can't put key pair")
	}

	return nil
}

func (client *Redis) Pull() (string, error) {
	if client.Client == nil {
		return "", errors.New("Empty redis client")
	}

	res, err := client.Client.LPop(client.Client.Context(), key).Result()

	if err != nil {
		return "", err
	}

	return res, nil
}

func (client *Redis) GetMap() (map[string]string, error) {
	return client.Client.HGetAll(client.Client.Context(), key).Result()
}

func (client *Redis) Close() {
	client.Client.Close()
}
