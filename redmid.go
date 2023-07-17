package redmid

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redmid interface {
	Middleware(next http.Handler) http.Handler
}

type redmid struct {
	redis  *redis.Client
	option *RedmidOption
}

func NewRedmid(redis *redis.Client, option *RedmidOption) Redmid {

	if option == nil {
		option = &RedmidOption{
			Expired: time.Second * 10,
		}
	}

	return &redmid{
		redis:  redis,
		option: option,
	}
}

func (redmid *redmid) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		newWriter := NewResponseWriter(w)

		// get the uri path
		path := r.URL.Path

		// get query params
		query := r.URL.Query()

		redisKey := path

		// loop through query params
		for key, value := range query {
			redisKey += key + ":" + value[0]
		}

		// check if key exists in redis
		val, err := redmid.redis.Get(ctx, redisKey).Result()
		if err != nil {
			// key does not exist, continue to next handler
			next.ServeHTTP(newWriter, r)

			// get response body from next handler
			data := newWriter.Data()

			status := newWriter.Code()

			redisData := Model{
				Status: status,
				Data:   data,
			}

			jsonData, err := json.Marshal(redisData)
			if err != nil {
				panic(err)
			}

			// write response body to redis
			redmid.redis.Set(ctx, redisKey, jsonData, redmid.option.Expired)

			return
		}

		// unmarshal response body from redis
		var redisData Model
		err = json.Unmarshal([]byte(val), &redisData)
		if err != nil {
			panic(err)
		}

		if redisData.Status != 0 {
			w.WriteHeader(redisData.Status)
		}

		w.Write([]byte(redisData.Data))
	})
}
