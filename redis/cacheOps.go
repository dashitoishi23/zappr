package redisutil

import (
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
	"github.com/gomodule/redigo/redis"
)

func DecodeCacheResponse[T any](client redis.Conn, cacheKey string) (T, error) {
	var cacheResp T

	cachedResponse, err := client.Do("GET", cacheKey)

	if err != nil {
		return cacheResp, err
	}

	if cachedResponse != nil {
		decodedResponse, err := util.JsonDecoder[T](cachedResponse.([]byte))

		return decodedResponse, err
	}

	return cacheResp, nil
}

func SetCache[T any](client redis.Conn, cacheKey string, dbResponse T) error {
	stringifiedResponse, err := util.StringifyJson(dbResponse)

	if err != nil {
		return err
	}

	_, setError := client.Do("SET", cacheKey, stringifiedResponse)

	return setError
}

func DeleteMultipleKeys(client redis.Conn, keys []string) error {
	for _, key := range keys {
		_, err := client.Do("DEL", key)

		if err != nil {
			return err
		}
	}

	return nil
}