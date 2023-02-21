package dataaccess

import (
	"context"
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

// SetRedis Setting Redis Database
func SetRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", setting.S.RedisHost, setting.S.RedisPort),
		Password: setting.S.RedisPassword,
		DB:       setting.S.RedisDatabase,
	})
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second) // Timeout 10s
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	DA = dataAccessRedis{
		client: client,
	}
}

func (da dataAccessRedis) Get(key string) (DataSet, error) {
	return DataSet{}, nil
}

func (da dataAccessRedis) Has(key string) (bool, error) {
	return false, nil
}

func (da dataAccessRedis) Set(key string, value DataSet) error {
	return nil
}

func (da dataAccessRedis) Delete(key string) error {
	return nil
}

func (da dataAccessRedis) ListAll() ([]param, error) {
	p := []param{}
	return p, nil
}

func (da dataAccessRedis) UpdateDate(key string) error {
	d, err := da.Get(key)
	if err != nil {
		return err
	}
	d.LastUpdatedAt = time.Now().Unix()
	return da.Set(key, d)
}

func (da dataAccessRedis) Close() error {
	return da.client.Close()
}
