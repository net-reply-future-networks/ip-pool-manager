package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"ip-pool-manager/handlers"
	"ip-pool-manager/ip"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
)

// flags for custom port number and addresses to run server and redis server
var (
	serverPort    = flag.Int("port", 3000, "port number")
	serverAddress = flag.String("address", "localhost", "port address")

	redisPort    = flag.Int("redis-port", 6379, "port number of redis server")
	redisAddress = flag.String("redis-address", "localhost", "port address of redis server")

	serviceName = flag.String("name", "ip-pool-manager", "name of service")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // this enables line logging

	flag.Parse()

	serverAddress := fmt.Sprintf("%v:%v", *serverAddress, *serverPort)
	rServerAddress := fmt.Sprintf("%v:%v", *redisAddress, *redisPort)

	myFigure := figure.NewFigure(*serviceName, "", true)
	myFigure.Print()

	log.Printf("INFO: Server address: %v\n", serverAddress)
	log.Printf("INFO: Redis address: %v\n", rServerAddress)

	// creating redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     rServerAddress, // redis address
		Password: "",             // no password set
		DB:       0,              // use default DB
	})

	addTestingIPs(rdb)

	go checkNotAvailableIPs(rdb)

	// creating chi multiplexer (router) for handlers
	r := chi.NewRouter()

	// setting middleware to log server actions and compressing JSON data
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5, "application/json"))

	// Get available single IP details from DB. Replaces the available IP with identical na-IP
	r.Get("/getIP", handlers.GetIP(rdb))
	// Get all available IP addresses from DB
	r.Get("/allAvailbleIPs", handlers.AllAvailbleIPs(rdb))
	// Delete an IP "Must include IP key name. Not just IP"
	r.Delete("/deleteIPfromPool", handlers.DeleteIPfromPool(rdb))
	// Create new IP and store into DB
	r.Post("/addIPtoPool", handlers.AddToIPtoPool(rdb))
	// Update IP details (Not create new IP)
	r.Put("/createNewIPpool", handlers.CreateNewIPinPool(rdb))

	srv := &http.Server{Addr: serverAddress, Handler: r}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		log.Println("INFO: Starting server")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed { // TODO: Will fail wrapped errors using != use errors.Is
				log.Fatal(err)
			}
		}
	}()

	<-stop

	log.Println("INFO: Received shutdown signal, shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func addTestingIPs(rdb *redis.Client) {
	IP1 := ip.IPpost{
		IPaddress: "a-185.9.249.220",
		Detail: ip.IPdetails{
			MACaddress: "89-43-5F-60-DC-76",
			LeaseTime:  time.Now(),
			Available:  true,
		},
	}

	IP2 := ip.IPpost{
		IPaddress: "na-102.131.46.22",
		Detail: ip.IPdetails{
			MACaddress: "20-F0-8F-95-CD-83",
			LeaseTime:  time.Now(),
			Available:  false,
		},
	}

	IP3 := ip.IPpost{
		IPaddress: "a-253.14.93.192",
		Detail: ip.IPdetails{
			MACaddress: "C2-A7-D2-35-8C-FD",
			LeaseTime:  time.Now(),
			Available:  true,
		},
	}

	sliceIPs := []ip.IPpost{IP1, IP2, IP3}

	ctx := context.Background()

	// Encodes and stores IP's into DB
	for _, IP := range sliceIPs {
		log.Println("DEBUG: Adding sample data to Database: v%", IP)
		//	Encode data into glob format to be stored into DB
		BufEnString, err := encodeIP(IP)
		if err != nil {
			log.Println("ERROR: Could not encode IP", err)
			continue
		}
		nameKey := IP.IPaddress

		err = rdb.Set(ctx, nameKey, BufEnString, 0).Err()
		if err != nil {
			log.Printf("ERROR: Could not set Key-Value in Redis: %v\n", err)
			continue
		}
	}
}

// Encodes IP into glob format
func encodeIP(ip ip.IPpost) (string, error) {
	// struct to Gob
	bufEn := &bytes.Buffer{}
	if err := gob.NewEncoder(bufEn).Encode(ip); err != nil {
		return "", err
	}
	BufEnString := bufEn.String()

	return BufEnString, nil
}

func checkNotAvailableIPs(rdb *redis.Client) {
	log.Println("INFO: Goroutine started in the background checking for expired leases")

	for {
		t1 := time.Now().Unix()
		log.Println("INFO: Checking for expired leases in order to free IP for reallocation")
		ctx := context.Background()

		// Loop used to iterate other each key that stars with "na-" in DB
		iter := rdb.Scan(ctx, 0, "na-*", 0).Iterator()
		for iter.Next(ctx) {
			// Storing each IP in DB
			foundIP, err := rdb.Get(ctx, iter.Val()).Result()
			if err != nil {
				log.Printf("ERROR: IP not found: %v\n", err)
				continue
			}

			// Gob to Struct
			bufDe := &bytes.Buffer{}
			bufDe.WriteString(foundIP)

			// Decode returned Gob format into IP struct
			var dataDecode ip.IPpost
			if err := gob.NewDecoder(bufDe).Decode(&dataDecode); err != nil {
				log.Printf("ERROR: Could not decode gob data: %v\n", err)
				continue
			}

			// Making sure that every Go routine create has a 5-second life span
			t2 := dataDecode.Detail.LeaseTime.Add(time.Second * 5).Unix() // TODO: Make 5 seconds a flag
			if t1 >= t2 {
				log.Printf("INFO: Lease expired for IP: %v , MAC: %v \n", dataDecode.IPaddress, dataDecode.Detail.MACaddress)
				replaceNAip(rdb, dataDecode)
			}
		}

		log.Println("INFO: Sleeping before checking again for expired leases")
		time.Sleep(5 * time.Second) // TODO: Make this a flag
	}
}

func replaceNAip(rdb *redis.Client, dataDecode ip.IPpost) {
	returnIP := ip.IPpost{
		IPaddress: strings.Replace(dataDecode.IPaddress, "na", "a", 1),
		Detail: ip.IPdetails{
			MACaddress: dataDecode.Detail.MACaddress,
			LeaseTime:  dataDecode.Detail.LeaseTime,
			Available:  true,
		},
	}
	// Convert IP struct into Gob format to store in DB
	bufEn := &bytes.Buffer{}
	if err := gob.NewEncoder(bufEn).Encode(returnIP); err != nil {
		log.Println(err)
		return // TODO: Return error
	}
	returnIPdecode := bufEn.String()

	ctx := context.Background()
	// Storing user key & value into db
	rdb.Set(ctx, returnIP.IPaddress, returnIPdecode, 0)

	// If IP doesn't exist throw an error
	if err := rdb.Del(ctx, dataDecode.IPaddress).Err(); err != nil {
		log.Println(dataDecode.IPaddress, "Cannot delete original IP: ", err)
	}
	log.Printf("INFO: IP is set free for reallocation: %v\n", returnIP.IPaddress)
}

// List of available curl requests
// curl "localhost:3000/allAvailbleIPs"
// curl "localhost:3000/getIP?key=a-185.9.249.220"
// curl -X DELETE "localhost:3000/deleteIPfromPool?key=na-102.131.46.22"
// curl -X POST -H 'content-type: application/json' --data '{"ip":"a-222.2.222.222","detail":{"MACaddress":"89-43-5F-60-DC-76","leaseTime":"2021-12-13T11:11:52.106975Z","available":true}}' http://localhost:3000/addIPtoPool
// curl -X PUT -H "Content-Type: application/json" -d '{"targetIp":"a-185.9.249.220","ip":"na-185.9.249.220","detail":{"MACaddress":"11-11-11-11-11-","leaseTime":"2021-12-13T11:11:52.106975Z","available":true}}' http://localhost:3000/createNewIPpool
