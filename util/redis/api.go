package redis

import (
	"github.com/go-redis/redis"
	"log"
)

var RedisClient *redis.ClusterClient

func InitRedis(clusters []string, poolSize, minIdleConns int, password string) {
	RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        clusters,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
		Password:     password})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Fatal("init redis pool failure!", err)
	}
}

//func InitRedis(clusters []string, poolSize, dialTimeout, idleTimeout, poolTimeout, readTimeout, writeTimeout int, password string) {
//	RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
//		Addrs:        clusters,
//		PoolSize:     poolSize,
//		DialTimeout:  time.Duration(dialTimeout) * time.Second,
//		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
//		PoolTimeout:  time.Duration(poolTimeout) * time.Second,
//		ReadTimeout:  time.Duration(readTimeout) * time.Second,
//		WriteTimeout: time.Duration(writeTimeout) * time.Second,
//		Password:     password})
//
//	_, err := RedisClient.Ping().Result()
//	if err != nil {
//		log.Fatal("init redis pool failure!", err)
//	}
//}

func Member(score int64, data interface{}) redis.Z {
	return redis.Z{float64(score), data}
}
