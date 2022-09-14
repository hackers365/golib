package redis_metrics

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type RedisV8Hook struct {
}

func (r *RedisV8Hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	nCtx := context.WithValue(ctx, "ts", time.Now())
	return nCtx, nil
}

func (r *RedisV8Hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if v := ctx.Value("ts"); v != nil {
		start := v.(time.Time)
		addRedisPrometheusLabelValue(cmd.Name(), float64(time.Since(start).Nanoseconds())/1e6)
	}
	return nil
}

func (r *RedisV8Hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (r *RedisV8Hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
