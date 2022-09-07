package redis_metrics

import (
	"fmt"
	"testing"
	//"time"

	//"github.com/gin-gonic/gin"
	//"github.com/penglongli/gin-metrics/ginmetrics"
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
