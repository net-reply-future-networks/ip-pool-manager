package handlers

import (
	"net/http"

	"github.com/go-redis/redis/v8"
)

func AddToIPtoPool(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
