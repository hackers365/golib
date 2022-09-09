package redis_metrics

import (
	"context"
	"fmt"
	"testing"
	"time"

	//"github.com/gin-gonic/gin"
	//"github.com/penglongli/gin-metrics/ginmetrics"
	redis_v8 "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedis(t *testing.T) {
	conf := &RedisConf{
		Host:    "192.168.208.214",
		Port:    6379,
		Pwd:     "ticket_dev",
		Timeout: 2,
		Db:      2,
	}

	/*metricRouter := gin.Default()
	SetMetrics(metricRouter)
	go metricRouter.Run(":8534")*/

	name := "default"

	_, err := NewRedisInstance(name, conf)
	assert.Equal(t, err, nil, "err must be nil")

	err = GetRedis(name).Set("key", "value", 0).Err()
	assert.Equal(t, err, nil, "err must be nil")

	val, err := GetRedis(name).Get("key").Result()
	assert.Equal(t, err, nil, "err must be nil")
	fmt.Println(val)
}

/*
func SetMetrics(r gin.IRoutes) {
	m := ginmetrics.GetMonitor()
	m.Expose(r)
}
*/

func TestRedisV8(t *testing.T) {
	addr := "192.168.208.214:6379"
	psw := "ticket_dev"
	rdb := redis_v8.NewClient(&redis_v8.Options{
		Addr:         addr,
		Password:     psw,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	rHook := &RedisV8Hook{}
	rdb.AddHook(rHook)

	err := rdb.Set(context.TODO(), "key", "value", 0).Err()
	assert.Equal(t, err, nil, "err must be nil")

	val, err := rdb.Get(context.TODO(), "key").Result()
	assert.Equal(t, err, nil, "err must be nil")
	fmt.Println(val)
}
