package main

import (
	"bufio"
	"fmt"
	"net/http"
)

func main() {
	//	URL port number needs to be changed manually
	resp, err := http.Get("http://localhost:3000/allUsers")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	//	reading response body and printing to console
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

// 	List of endpoints API can use
//	resp, err := http.Get("http://localhost:3000/allUsers")
//	resp, err := http.Get("http://localhost:3000/getUser?key=key1")

//	curl -X DELETE "localhost:3000/deleteUser?key=user2"
//  curl -X POST -H 'content-type: application/json' --data '{"key": "john","vaye": "stuff"}' http://localhost:3000/createUser
