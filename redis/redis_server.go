package redis

import (
	"fmt"
)

import (
	"github.com/go-redis/redis"
)

type RedisConf struct {
	Host    string `json:"Host"`
	Port    int    `json:"Port"`
	Pwd     string `json:"pwd"`
	Timeout int    `json:"Timeout"`
	Db      int    `json:"db"`
}

type MRedis struct {
	redisClient *redis.Client
}

var instance = initRedis()

//init
func initRedis() *MRedis {
	m := &MRedis{}
	return m
}

//get
func GetRedis() *redis.Client {
	return instance.redisClient
}

//init instance
func NewRedisInstance(conf *RedisConf) (*MRedis, error) {
	redisClient, err := initRedisInstance(conf.Host, conf.Port, conf.Pwd, conf.Db)
	if err != nil {
		return nil, err
	}

	instance.redisClient = redisClient

	return instance, nil
}

func initRedisInstance(host string, port int, passwd string, db int) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})

	return c, nil
}
