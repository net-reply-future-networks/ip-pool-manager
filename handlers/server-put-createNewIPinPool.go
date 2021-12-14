package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"ip-pool-manager/IP"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type putIPpost struct {
	TargetIP  string       `json:"targetIp"`
	IPaddress string       `json:"ip"`
	Detail    putIPdetails `json:"detail"`
}

type putIPdetails struct {
	MACaddress string    `json:"MACaddress"`
	LeaseTime  time.Time `json:"leaseTime"`
	Available  bool      `json:"available"`
}

func CreateNewIPinPool(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		//	Creating a empty user post called "uIP"
		var uPut putIPpost

		//	Decodes response JSON into a userPostIP object and catches any errors
		if err := json.NewDecoder(r.Body).Decode(&uPut); err != nil {
			fmt.Println("ERR: ", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Cannot decode request"))

			return
		}

		//	Checking if user IP value is correct lengh
		if len(uPut.IPaddress) != 15 && len(uPut.IPaddress) != 16 {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Println(len(uPut.IPaddress))
			resp := uPut.IPaddress + " IP is not correct length . May need to contain a- or na-"
			w.Write([]byte(resp))

		}

		ctx := context.Background()
		_, err := rdb.Get(ctx, uPut.TargetIP).Result()
		switch {
		case err == redis.Nil:
			fmt.Println("key does not exist")
			fmt.Println(uPut.TargetIP)
			return
		case err != nil:
			fmt.Println("Get failed", err)
			return
		}

		tempIPpost := IP.IPpost{
			IPaddress: uPut.IPaddress,
			Detail: IP.IPdetails{
				MACaddress: uPut.Detail.MACaddress,
				LeaseTime:  uPut.Detail.LeaseTime,
				Available:  uPut.Detail.Available,
			},
		}

		newIPencoded := encodeIP(tempIPpost)

		rdb.Rename(ctx, uPut.TargetIP, uPut.IPaddress)
		//	Storing user key & value into db
		rdb.Set(ctx, uPut.IPaddress, newIPencoded, 0)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("IP has changed"))

	}
}

// curl -X PUT -H "Content-Type: application/json" -d '{"targetIp":"a-185.9.249.220","ip":"na-185.9.249.220","detail":{"MACaddress":"11-11-11-11-11-","leaseTime":"2021-12-13T11:11:52.106975Z","available":true}}' http://localhost:3000/createNewIPpool
