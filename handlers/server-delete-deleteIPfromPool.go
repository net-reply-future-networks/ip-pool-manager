package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
)

//	Deletes the specified IP from the IP pool
func DeleteIPfromPool(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		//	Storing "key" url param that contains IP key/id
		param := r.URL.Query().Get("key")

		//	Delete specified IP key from db
		if param != "" {
			//	If IP doesn't exist throw an err
			if err := rdb.Del(ctx, param).Err(); err != nil {
				http.Error(w, http.StatusText(404), 404)
				fmt.Println(param, " Error: ", err)
			} else {
				//	Send back ok status response and "User Deleted" message
				responseMsg := param + " IP Deleted"
				w.Header().Set("content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(responseMsg))
			}
		}
	}
}

// curl -X DELETE "localhost:3000/deleteIPfromPool?key=na-102.131.46.22"
