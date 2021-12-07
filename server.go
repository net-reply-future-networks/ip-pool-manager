package main

import (
	"flag"
	"fmt"
	"go-chi-example/handlers"
	"net/http"

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
		Addr:     rServerAddress, //change!!!!!!!!!!!!!!!!
		Password: "",             // no password set
		DB:       0,              // use default DB
	})

	//	creating chi multiplexor (router) for handlers
	r := chi.NewRouter()

	//	setting middlewear to log server actions and compressing JSON data
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5, "application/json"))

	//	Gets all users Key/Value pairs
	r.Get("/allUsers", handlers.GetAllUsers(rdb))
	//	Get only a single user
	r.Get("/getUser", handlers.GetUser(rdb))
	//	Delete a single user
	r.Delete("/deleteUser", handlers.DeleteUser(rdb))
	//	Post a single user
	r.Post("/createUser", handlers.CreateUser(rdb))
	//	Update user value
	r.Put("/updateUser", handlers.UpdateUser(rdb))

	//	update user
	//	change db and justify it
	//	Make todo list

	http.ListenAndServe(serverAddress, r)
}

//	curl "localhost:3000/allUsers"
//	curl "localhost:3000/getUser?key=name1'
//	curl -X DELETE "localhost:3000/deleteUser?key=user2"
//  curl -X POST -H 'content-type: application/json' --data '{"key": "john","vaye": "stuff"}' http://localhost:3000/createUser

//	go run server.go --port 8080 --address 0.0.0.0 --redisPort 6378 --redisAddress 0.0.0.0
