package redisutil

import (
	"os"

	"github.com/gomodule/redigo/redis"
)

type RedisPool struct {
	Pool *redis.Pool
}

func(r *RedisPool) NewPool()  {
	pool := &redis.Pool{
		MaxIdle: 80,
		MaxActive: 0,
		Dial: func() (redis.Conn, error) {
			redisHost := os.Getenv("REDIS_HOST")

			if redisHost == "" {
				redisHost = "redis"
			}

			c, err := redis.Dial("tcp", redisHost + ":6379")
			
			if err != nil {
				panic(err)
			}

			redisPwd := os.Getenv("REDIS_PASSWORD")

			if redisPwd != "" {
				_, authErr := c.Do("AUTH", redisPwd)

				if authErr != nil {
				panic(authErr)
			}
			}

			return c, err

		},
	}

	r.Pool = pool
}

func(r *RedisPool) GetPool() *redis.Pool {
	redisPool := &r.Pool

	return *redisPool
	
}



