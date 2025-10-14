package limitter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimitter struct {
	Client *redis.Client
	Window time.Duration
	Limit  int16
}

func NewRedisLimitter(client *redis.Client, wdw time.Duration, lim int16) *RedisLimitter {
	return &RedisLimitter{
		Client: client,
		Window: wdw,
		Limit:  lim,
	}
}

func (r *RedisLimitter) AllowPost(ctx context.Context, key string) (bool, error) {

	currentTime := time.Now().UnixMicro()
	//seconds := 300
	//allowedWindow := currentTime - int64(seconds)
	min := currentTime - r.Window.Nanoseconds()

	//expiry := 30*time.Second

	redisKey := fmt.Sprintf("ratelimit:%s", key)

	pipe := r.Client.TxPipeline()

	//atomic transaction
	pipe.ZRemRangeByScore(ctx, redisKey, "0", fmt.Sprintf("%d", min))

	countCmd := pipe.ZAdd(ctx, redisKey, redis.Z{Score: float64(currentTime), Member: float64(currentTime)})

	pipe.ZCard(ctx, redisKey)

	pipe.Expire(ctx, redisKey, r.Window*2)

	//exceute all the commands in the pipe all at one

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("issue executing the commands due to: %s", err)
		return false, err
	}

	count, err := countCmd.Result()
	if err != nil {
		return false, fmt.Errorf("failed to count commands due to: ", err)
	}

	if int(count) > int(r.Limit) {
		return false, nil // deny the request
	}

	return true, nil // request is allowed
}
