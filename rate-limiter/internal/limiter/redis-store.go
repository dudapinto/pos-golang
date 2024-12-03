package limiter

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr, password string) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // Seleciona o banco de dados
	})

	return &RedisStore{client: rdb}
}

func (rs *RedisStore) Get(key string) (*RateLimit, error) {
	ctx := context.Background()
	val, err := rs.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	limit, err := strconv.Atoi(val["limit"])
	if err != nil {
		return nil, err
	}

	blockDuration, err := strconv.ParseInt(val["block_duration"], 10, 64)
	if err != nil {
		return nil, err
	}

	lastRequest, err := strconv.ParseInt(val["last_request"], 10, 64)
	if err != nil {
		return nil, err
	}

	requests, err := strconv.Atoi(val["requests"])
	if err != nil {
		return nil, err
	}

	rl := &RateLimit{
		LimitPerSecond: limit,
		BlockDuration:  time.Duration(blockDuration) * time.Second,
		LastRequest:    time.Unix(lastRequest, 0),
		Requests:       requests,
	}

	return rl, nil
}

func (rs *RedisStore) Set(key string, rl *RateLimit) error {
	// Convertemos para segundos (para armazenar no Redis)
	blockDurationSeconds := int64(rl.BlockDuration.Seconds())

	// Criamos um hash no Redis para armazenar as informações do rate limit
	data := map[string]interface{}{
		"limit":          rl.LimitPerSecond,
		"block_duration": blockDurationSeconds,
		"last_request":   rl.LastRequest.Unix(),
		"requests":       rl.Requests,
	}

	err := rs.client.HSet(context.Background(), key, data).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rs *RedisStore) Increment(key string) error {
	// Incrementa o campo "requests" do hash
	err := rs.client.HIncrBy(context.Background(), key, "requests", 1).Err()
	if err != nil {
		return err
	}
	return nil
}
