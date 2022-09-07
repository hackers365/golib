package redis_metrics

import (
  "fmt"
  "sync"
  "time"

  "github.com/penglongli/gin-metrics/ginmetrics"

  "github.com/go-redis/redis"
)

type RedisConf struct {
  Host    string `json:"Host"`
  Port    int    `json:"Port"`
  Pwd     string `json:"pwd"`
  Timeout int    `json:"Timeout"`
  Db      int    `json:"db"`
}

type mredis struct {
  sync.RWMutex
  redisMap map[string]*redis.Client
}

func (m *mredis) Add(name string, instance *redis.Client) {
  m.Lock()
  defer m.Unlock()
  m.redisMap[name] = instance
}

func (m *mredis) Get(name string) *redis.Client {
  m.RLock()
  defer m.RUnlock()
  if redisInstance, ok := m.redisMap[name]; ok {
    return redisInstance
  }
  return nil
}

var mRedis *mredis

func init() {
  mRedis = &mredis{
    redisMap: map[string]*redis.Client{},
  }
}

func GetRedis(name string) *redis.Client {
  return mRedis.Get(name)
}

//init instance
func NewRedisInstance(name string, conf *RedisConf) (*redis.Client, error) {
  redisClient, err := initRedisInstance(conf.Host, conf.Port, conf.Pwd, conf.Db)
  if err != nil {
    return nil, err
  }
  // 注册指标
  err = registerRedisPrometheus()
  if err != nil {
    return nil, err
  }

  redisCollector.Client = new(RedisClient)
  redisCollector.execDurationHistogram = HistogramRedisMetric
  addRedisExecDuration(redisClient)
  mRedis.Add(name, redisClient)

  _, err = redisClient.Ping().Result()
  if err != nil {
    return nil, err
  }

  return redisClient, nil
}

var redisCollector = &RedisCollector{}

func initRedisInstance(host string, port int, passwd string, db int) (*redis.Client, error) {
  addr := fmt.Sprintf("%s:%d", host, port)

  c := redis.NewClient(&redis.Options{
    Addr:     addr,
    Password: passwd,
    DB:       db,
  })
  return c, nil
}

// HistogramRedisMetric 声明prometheus 指标
var (
  HistogramRedisMetric = &ginmetrics.Metric{
    Type:        ginmetrics.Histogram,
    Name:        "redis_operate_duration_milliseconds",
    Description: "an example of gauge type metric",
    Buckets:     []float64{0.1, 0.5, 1, 2, 3, 5, 10, 20, 50, 100},
    Labels:      []string{"cmd"},
  }
)

// RegisterRedisPrometheus 注册 redis prometheus 指标
func registerRedisPrometheus() error {
  return ginmetrics.GetMonitor().AddMetric(HistogramRedisMetric)
}

func addRedisPrometheusLabelValue(cmd string, constTime float64) {
  var label []string
  label = append(label, cmd)
  ginmetrics.GetMonitor().GetMetric("redis_operate_duration_milliseconds").Observe(label, constTime)
}

type RedisClient interface {
  WrapProcess(fn func(oldProcess func(cmd redis.Cmder) error) func(cmd redis.Cmder) error)
  WrapProcessPipeline(old func(old func([]redis.Cmder) error) func([]redis.Cmder) error)
}

type RedisCollector struct {
  Client                *RedisClient
  once                  sync.Once
  execDurationHistogram *ginmetrics.Metric
}

func addRedisExecDuration(client RedisClient) {
  client.WrapProcess(func(oldProcess func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
    return func(cmd redis.Cmder) error {
      start := time.Now()
      err := oldProcess(cmd)
      addRedisPrometheusLabelValue(cmd.Name(), float64(time.Since(start).Nanoseconds())/1e6)
      return err
    }
  })

  client.WrapProcessPipeline(func(oldProcess func([]redis.Cmder) error) func([]redis.Cmder) error {
    return func(cmds []redis.Cmder) error {
      err := oldProcess(cmds)
      return err
    }
  })
}
