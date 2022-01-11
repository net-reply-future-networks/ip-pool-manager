package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
)

// readyz is a readiness probe.
func Readyz(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		ctx := context.Background()

		// Pings DB and sends bad request if DB does not ping back
		pong, err := rdb.Ping(ctx).Result()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("DB did not ping back. DB is not running \n")) //nolint:errcheck
			return
		}
		fmt.Println(pong, err)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DB did ping back. DB is running \n")) //nolint:errcheck
	}
}
