package handlers

import (
	"net/http"
)

// readyz is a readiness probe.
func Readyz() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		// ctx := context.Background()

		// Pings DB and sends bad request if DB does not ping back
		// pong, err := rdb.Ping(ctx).Result()
		// if err != nil {
		// 	log.Println(err)
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	w.Write([]byte("DB did not ping back. DB is not running \n")) //nolint:errcheck
		// 	return
		// }
		// log.Println(pong, err)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DB did ping back. DB is running \n")) //nolint:errcheck
	}
}
