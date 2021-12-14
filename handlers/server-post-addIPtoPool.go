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

	"github.com/go-redis/redis/v8"
)

func AddToIPtoPool(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		// Creating a empty user post called "u"
		var u IP.IPpost

		// Decodes response JSON into a userPostIP object and catches any errors
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			fmt.Println("ERR: ", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Cannot decode request")) //nolint:errcheck
			return
		}

		// Checking if user IP value is correct lengh
		if len(u.IPaddress) != 15 {
			w.WriteHeader(http.StatusBadGateway)
			resp := "IP is not correct length " + u.IPaddress
			w.Write([]byte(resp)) //nolint:errcheck
			return
		}

		ctx := context.Background()

		// These print are for debug purposes
		log.Printf("Values of new IP. IP address : %v. Value: %v", u, u.IPaddress)

		encodedU := encodeIP(u)
		// Storing user key & value into db
		rdb.Set(ctx, u.IPaddress, encodedU, 0)

		userResponse := u.IPaddress + "IP has been added to DB"
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userResponse)) //nolint:errcheck
	}

}

// Encodes IP into glob format
func encodeIP(IP IP.IPpost) string {
	// struct to Gob
	bufEn := &bytes.Buffer{}
	if err := gob.NewEncoder(bufEn).Encode(IP); err != nil {
		fmt.Println(err)
	}
	BufEnString := bufEn.String()

	return BufEnString
}

// curl -X POST -H 'content-type: application/json' --data '{"ip":"a-222.2.222.222","detail":{"MACaddress":"89-43-5F-60-DC-76","leaseTime":"2021-12-13T11:11:52.106975Z","available":true}}' http://localhost:3000/addIPtoPool
