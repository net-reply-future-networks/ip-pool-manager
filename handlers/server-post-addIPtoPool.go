package handlers

import (
	"net/http"

	"github.com/go-redis/redis/v8"
)

func AddToIPtoPool(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Temporary text
		w.Header().Set("content-type", "application/json")
	}
}
