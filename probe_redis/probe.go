package probe_redis

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	REDIS_KEY            = "ProbeRedisKey"
	REDIS_EXPIRATION     = time.Hour * 1
	REDIS_DATA_GEN_COUNT = 100
)

var redisOpts = redis.Options{
	Network:     "tcp",
	Addr:        "localhost:9100",
	Password:    "",
	DB:          0,
	PoolSize:    50,
	DialTimeout: time.Second * 5,
}

type ContractorInfo struct {
	ContractorName string
	ContractorInns []string
}

type ProbeRedis struct {
	name string
}

func (probe ProbeRedis) Name() string {
	return probe.name
}

type ProdeRedisTestStruct struct {
	Index int `json:"index"`
}

func (probe ProbeRedis) Present(ctx *gin.Context) (any, error) {
	redisCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	client := redis.NewClient(&redisOpts)
	defer client.Close()

	localHKey, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	setData := make([]ProdeRedisTestStruct, REDIS_DATA_GEN_COUNT)
	for i := range REDIS_DATA_GEN_COUNT {
		setData[i] = ProdeRedisTestStruct{Index: i}
	}

	valueJson, err := json.Marshal(setData)
	if err != nil {
		return nil, err
	}
	respSet := client.HSet(redisCtx, REDIS_KEY, localHKey.String(), valueJson)
	if err := respSet.Err(); err != nil {
		return nil, err
	}

	result, err := client.HGet(redisCtx, REDIS_KEY, localHKey.String()).Result()
	if err != nil {
		return nil, err
	}

	ret := make([]ProdeRedisTestStruct, REDIS_DATA_GEN_COUNT)
	json.Unmarshal([]byte(result), &ret)

	if err := client.HDel(redisCtx, REDIS_KEY, localHKey.String()).Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func Realization() ProbeRedis {
	return ProbeRedis{"redis"}
}
