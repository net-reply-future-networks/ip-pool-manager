package handlers

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"ip-pool-manager/IP"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
)

//	Return specified IP's
func GetIP(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		ctx := context.Background()

		//	Storing "key" url param that contains IP key/id
		param := r.URL.Query().Get("key")

		//	Check if URL param is empty or if specified IP is not availble
		if param == "" {
			fmt.Println("Empty URL parameter")
			w.Write([]byte("Empty URL parameter"))
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if strings.HasPrefix(param, "na-") {
			fmt.Println("Must select available IP")
			w.Write([]byte("Must select available IP"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//	Retrieving IP stored in DB
		val, err := rdb.Get(ctx, param).Result()
		if err != nil {
			fmt.Println("Cannot find IP")
			w.Write([]byte("Cannot find IP"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Gob to Struct
		bufDe := &bytes.Buffer{}
		bufDe.WriteString(val)

		//	Decode returned Gob format into IP struct
		var valDecode IP.IPpost
		if err := gob.NewDecoder(bufDe).Decode(&valDecode); err != nil {
			log.Println(err)
		}

		// Converting Struct into JSON byte to return to user
		responseIP, err := json.Marshal(valDecode)
		if err != nil {
			panic(err)
		}

		//	Send back ok status response and specified IP
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseIP))
	}

}

//	curl "localhost:3000/getIP?key=a-185.9.249.220"
