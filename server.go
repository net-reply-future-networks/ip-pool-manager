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

	addTestingIP(rdb)

	//	creating chi multiplexor (router) for handlers
	r := chi.NewRouter()

	//	setting middlewear to log server actions and compressing JSON data
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5, "application/json"))

	//	Get all availble IP addresses
	r.Get("/getIPpoolInfo", handlers.GetIPpoolInfo(rdb))
	//	Delete a single user
	r.Delete("/deleteIPfromPool", handlers.DeleteIPfromPool(rdb))
	//	Post a single user
	r.Post("/addIPtoPool", handlers.AddToIPtoPool(rdb))
	//	Update user value
	r.Put("/createNewIPpool", handlers.CreateNewIPinPool(rdb))

	http.ListenAndServe(serverAddress, r)

}

func addTestingIP(rdb *redis.Client) {
	IP1 := IP.IPpost{
		IPaddress: "185.9.249.220",
		Detail: IP.IPdetails{
			MACaddress: "89-43-5F-60-DC-76",
			LeaseTime:  time.Now(),
			Available:  true,
		},
	}

	IP2 := IP.IPpost{
		IPaddress: "102.131.46.22",
		Detail: IP.IPdetails{
			MACaddress: "20-F0-8F-95-CD-83",
			LeaseTime:  time.Now(),
			Available:  false,
		},
	}

	IP3 := IP.IPpost{
		IPaddress: "253.14.93.192",
		Detail: IP.IPdetails{
			MACaddress: "C2-A7-D2-35-8C-FD",
			LeaseTime:  time.Now(),
			Available:  false,
		},
	}

	sliceIPs := []IP.IPpost{IP1, IP2, IP3}

	ctx := context.Background()

	// Encodes and stores IP's into DB
	for _, IP := range sliceIPs {
		//	Encode data into glob format to be stored into DB
		BufEnString := encodeIP(IP)
		nameKey := "a" + IP.IPaddress

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

//	Return all IP's as a string
// func getTestingIP(rdb *redis.Client) {

// 	ctx := context.Background()

// 	//	Loop used to iterate other each key in DB
// 	iter := rdb.Scan(ctx, 0, "*", 0).Iterator()

// 	// fmt.Println("--------------------------------")
// 	// fmt.Println("Getting all users...")

// 	for iter.Next(ctx) {
// 		//	Storing each user value (key value)
// 		foundIP, err := rdb.Get(ctx, iter.Val()).Result()
// 		if err != nil {
// 			fmt.Println(err)
// 		} else {
// 			fmt.Println("found IP: ", foundIP)
// 		}

// 	}

// 	foundUser := getIP(rdb, "1")
// 	fmt.Println("Found user --> ", foundUser)

// 	// //	Convert JSON byte data to IP struct
// 	// structIPs := IPsJsonToStruct(jsonIPs)
// }

//	Convert struct to JSON byte data
// func structIPtoJson(IP IP.IPpost) []byte {
// 	// Struct to JSON
// 	var jsonIPmarshal []byte

// 	jsonIPmarshal, err := json.Marshal(IP)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	return jsonIPmarshal
// }

// return IP from DB as a string (unmarshaled)
// func getIP(rdb *redis.Client, Key string) string {
// 	ctx := context.Background()

// 	val, err := rdb.Get(ctx, Key).Result()
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	var dataUnmarshal IP.IPpost
// 	err = json.Unmarshal(val, &dataUnmarshal)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return val
// }

//	Convert a slice of data into JSON byte data
// func IPsJsonToStruct(jsonIPs []byte) []IP.IPpost {
// 	IPstruct := make([]IP.IPpost, 0)
// 	json.Unmarshal(jsonIPs, &IPstruct)
// 	println("-----IP addresses stored-----")
// 	for _, user := range IPstruct {
// 		fmt.Println(user)
// 	}
// 	return IPstruct
// }

//	curl "localhost:3000/allUsers"
//	curl "localhost:3000/getUser?key=name1'
//	curl -X DELETE "localhost:3000/deleteUser?key=user2"
//  curl -X POST -H 'content-type: application/json' --data '{"key": "john","vaye": "stuff"}' http://localhost:3000/createUser

//	go run server.go --port 8080 --address 0.0.0.0 --redisPort 6378 --redisAddress 0.0.0.0
