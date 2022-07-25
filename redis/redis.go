package redisutil

import "github.com/gomodule/redigo/redis"

type RedisPool struct {
	Pool *redis.Pool
}

func(r *RedisPool) NewPool()  {
	pool := &redis.Pool{
		MaxIdle: 80,
		MaxActive: 0,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			
			if err != nil {
				panic(err)
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



