package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"ip-pool-manager/IP"
	"ip-pool-manager/handlers"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
)

var (
	//	flag for custom port number and addresses to run server and redis server
	portNum     = flag.Int("port", 3000, "Please enter port number")
	portAddress = flag.String("address", "localhost", "Please enter a port address")

	rPortNum     = flag.Int("redisPort", 6379, "Please enter port number for redis server")
	rPortAddress = flag.String("redisAddress", "localhost", "Please enter a port address for redis server")
)

func main() {

	flag.Parse()

	serverAddress := fmt.Sprintf("%v:%v", *portAddress, *portNum)
	rServerAddress := fmt.Sprintf("%v:%v", *rPortAddress, *rPortNum)

	fmt.Println(serverAddress)
	fmt.Println(rServerAddress)

	//	creating redis server
	rdb := redis.NewClient(&redis.Options{
		Addr:     rServerAddress, // redis address
		Password: "",             // no password set
		DB:       0,              // use default DB
	})

	addTestingIPs(rdb)

	go checkNotAvailableIPs(rdb)
	//	creating chi multiplexor (router) for handlers
	r := chi.NewRouter()

	//	setting middlewear to log server actions and compressing JSON data
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5, "application/json"))

	//	Get available single IP details from DB. Replaces the availble IP with idential na-IP
	r.Get("/getIP", handlers.GetIP(rdb))
	//	Get all available IP addresses from DB
	r.Get("/allAvailbleIPs", handlers.AllAvailbleIPs(rdb))
	//	Delete an IP "Must include IP key name. Not just IP"
	r.Delete("/deleteIPfromPool", handlers.DeleteIPfromPool(rdb))
	//	Create new IP and store into DB
	r.Post("/addIPtoPool", handlers.AddToIPtoPool(rdb))
	//	Update IP details (Not create new IP)
	r.Put("/createNewIPpool", handlers.CreateNewIPinPool(rdb))

	err := http.ListenAndServe(serverAddress, r)
	if err != nil {
		fmt.Println(err)
	}

}

func addTestingIPs(rdb *redis.Client) {
	IP1 := IP.IPpost{
		IPaddress: "a-185.9.249.220",
		Detail: IP.IPdetails{
			MACaddress: "89-43-5F-60-DC-76",
			LeaseTime:  time.Now(),
			Available:  true,
		},
	}

	IP2 := IP.IPpost{
		IPaddress: "na-102.131.46.22",
		Detail: IP.IPdetails{
			MACaddress: "20-F0-8F-95-CD-83",
			LeaseTime:  time.Now(),
			Available:  false,
		},
	}

	IP3 := IP.IPpost{
		IPaddress: "a-253.14.93.192",
		Detail: IP.IPdetails{
			MACaddress: "C2-A7-D2-35-8C-FD",
			LeaseTime:  time.Now(),
			Available:  true,
		},
	}

	sliceIPs := []IP.IPpost{IP1, IP2, IP3}

	ctx := context.Background()

	// Encodes and stores IP's into DB
	for _, IP := range sliceIPs {
		//	Encode data into glob format to be stored into DB
		BufEnString := encodeIP(IP)
		nameKey := IP.IPaddress

		err1 := rdb.Set(ctx, nameKey, BufEnString, 0).Err()
		if err1 != nil {
			log.Println(err1)
		}
	}

	//	Loop used to iterate other each key in DB
	iter := rdb.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		//	Storing each IP in DB
		foundIP, err := rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			fmt.Println("IP not found. ERR: ", err)
		} else {
			// fmt.Println(foundIP)

			// Gob to Struct
			bufDe := &bytes.Buffer{}

			bufDe.WriteString(foundIP)

			//	Decode returned Gob format into IP struct
			var dataDecode IP.IPpost
			if err := gob.NewDecoder(bufDe).Decode(&dataDecode); err != nil {
				log.Println(err)
			}
			// fmt.Println("data decoded from gob:", dataDecode)

		}

	}

}

//	Encodes IP into glob format
func encodeIP(IP IP.IPpost) string {
	// struct to Gob
	bufEn := &bytes.Buffer{}
	if err := gob.NewEncoder(bufEn).Encode(IP); err != nil {
		panic(err)
	}
	BufEnString := bufEn.String()

	return BufEnString
}

func checkNotAvailableIPs(rdb *redis.Client) {
	for i := 1; i < 10000; i++ {
		t1 := time.Now().Unix()

		ctx := context.Background()
		//	Loop used to iterate other each key that stars with "a-" in DB
		iter := rdb.Scan(ctx, 0, "na-*", 0).Iterator()
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

			//	Making sure that every Go routine create has a 5 second life span
			t2 := dataDecode.Detail.LeaseTime.Add(time.Second * 5).Unix()

			if t1 >= t2 {
				fmt.Println("GO ROUTINE EXPIRED")
				replaceNAip(rdb, dataDecode)

			}
		}

		time.Sleep(5 * time.Second)
	}

}

func replaceNAip(rdb *redis.Client, dataDecode IP.IPpost) {
	returnIP := IP.IPpost{
		IPaddress: strings.Replace(dataDecode.IPaddress, "na", "a", 1),
		Detail: IP.IPdetails{
			MACaddress: dataDecode.Detail.MACaddress,
			LeaseTime:  dataDecode.Detail.LeaseTime,
			Available:  true,
		},
	}
	// Convert IP struct into Gob format to store in DB
	bufEn := &bytes.Buffer{}
	if err := gob.NewEncoder(bufEn).Encode(returnIP); err != nil {
		panic(err)
	}
	returnIPdecode := bufEn.String()

	ctx := context.Background()
	//	Storing user key & value into db
	rdb.Set(ctx, returnIP.IPaddress, returnIPdecode, 0)

	//	If IP doesn't exist throw an err
	if err := rdb.Del(ctx, dataDecode.IPaddress).Err(); err != nil {
		fmt.Println(dataDecode.IPaddress, "Cannot delete original IP: ", err)
	}
	fmt.Println("deleted old na IP")
}

//	go run server.go --port 8080 --address 0.0.0.0 --redisPort 6378 --redisAddress 0.0.0.0
