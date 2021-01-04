package dataaccess

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/potproject/ikuradon-salmon/setting"
)

type dataAccessRedis struct {
	DataAccess
	client *redis.Client
}

var ctx = context.Background()
var timeout = time.Second * 10 // Timeout 10s

// SetRedis Setting Redis Database
func SetRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", setting.S.RedisHost, setting.S.RedisPort),
		Password: setting.S.RedisPassword,
		DB:       setting.S.RedisDatabase,
	})
	ctxt, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := client.Ping(ctxt).Result()
	if err != nil {
		log.Fatal(err)
	}
	DA = dataAccessRedis{
		client: client,
	}
}

func (da dataAccessRedis) Get(key string) (DataSet, error) {
	var ds DataSet
	ctxt, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	data, err := da.client.Get(ctxt, "key").Result()
	if err != nil {
		return ds, err
	}
	err = json.Unmarshal([]byte(data), &ds)
	if err != nil {
		return ds, err
	}
	return ds, nil
}

func (da dataAccessRedis) Has(key string) (bool, error) {
	ctxt, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	res, err := da.client.Exists(ctxt, key).Result()
	return res != 0, err
}

func (da dataAccessRedis) Set(key string, value DataSet) error {
	j, jsonErr := json.Marshal(value)
	if jsonErr != nil {
		return jsonErr
	}
	ctxt, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	redisErr := da.client.Set(ctxt, key, j, 0).Err()
	if redisErr != nil {
		return redisErr
	}
	return nil
}

func (da dataAccessRedis) Delete(key string) error {
	ctxt, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := da.client.Del(ctxt, key).Result()
	return err
}

// TODO
func (da dataAccessRedis) ListAll() ([]param, error) {
	p := []param{}
	return p, nil
}

func (da dataAccessRedis) Close() error {
	return da.client.Close()
}
