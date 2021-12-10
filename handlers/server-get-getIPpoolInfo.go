package handlers

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"ip-pool-manager/IP"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func GetIPpoolInfo(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allAvailbleIPs := getAllIPs(rdb)

		fmt.Println(allAvailbleIPs)
		// w.Header().Set("content-type", "application/json")
		// w.WriteHeader(http.StatusOK)
		// //	Response is a sinlge string containg all users and values (key value pairs)
		// w.Write([]byte(allAvailbleIPs))

	}

}

func getAllIPs(rdb *redis.Client) []IP.IPpost {
	fmt.Println("GET ALL USERSSS!!!!!!")
	ctx := context.Background()
	allIPs := []IP.IPpost{}

	//	Loop used to iterate other each key in DB
	iter := rdb.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		//	Storing each IP in DB
		foundIP, err := rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			fmt.Println("IP not found. ERR: ", err)
			continue
		}

		// Gob to Struct
		bufDe := &bytes.Buffer{}

		bufDe.WriteString(foundIP)

		//	Decode returned Gob format into IP struct
		var dataDecode IP.IPpost
		if err := gob.NewDecoder(bufDe).Decode(&dataDecode); err != nil {
			log.Println(err)
			continue
		}

		allIPs = append(allIPs, dataDecode)
	}
	allAvailbleIPs := findAvailbleIP(allIPs)
	return allAvailbleIPs
}

func findAvailbleIP(allIPs []IP.IPpost) []IP.IPpost {
	allAvailbleIPs := []IP.IPpost{}
	for _, IP := range allIPs {
		if IP.Detail.Available {
			allAvailbleIPs = append(allAvailbleIPs, IP)
		} else {
			fmt.Println(IP.IPaddress, " is not availble right now")
		}
	}

	fmt.Println("ALL AVAILABLE IP'S --> ", allAvailbleIPs)
	return allAvailbleIPs
}
