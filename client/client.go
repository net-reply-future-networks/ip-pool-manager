package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"ip-pool-manager/ip"
	"log"
	"net/http"
	"time"
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

// Flag for user request selection
var requestOption = flag.Int("request", 1, "Enter number for request type. 1)Get all IPs | 2)Get single IP | 3)Delete IP | 4)Post IP | 5)Put IP ")

func main() {
	// Enables line logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	switch *requestOption {
	case 1:
		fmt.Println("Select all IPs")
		// Get request to return all availble IPs
		resp, err := http.Get("http://localhost:3000/allAvailbleIPs")
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
	case 2:
		fmt.Println("Select single IP")

		// Get request to return specific IP
		resp, err := http.Get("http://localhost:3000/getIP?key=a-185.9.249.220")
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
	case 3:
		fmt.Println("Delete single IP")

		// Delete request to delete specified IP
		req, err := http.NewRequest("DELETE", "http://localhost:3000/deleteIPfromPool?key=a-102.131.46.22", nil)
		if err != nil {
			fmt.Println("err1")
			panic(err)
		}

		// Reads resp from request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("err2")
			panic(err)
		}
		defer resp.Body.Close()

		//	reading response body and printing to console
		scanner := bufio.NewScanner(resp.Body)
		for i := 0; scanner.Scan() && i < 5; i++ {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("err3")
			panic(err)
		}
	case 4:
		fmt.Println("Post new IP a-111.11.11.111")

		// Creating new IP with dummy data
		data := ip.IPpost{
			IPaddress: "a-111.11.11.111",
			Detail: ip.IPdetails{
				MACaddress: "A1-A2-A3-A4-A5-A6",
				LeaseTime:  time.Now(),
				Available:  true,
			},
		}

		// Converts struct data to JSON byte data
		payloadBytes, err := json.Marshal(data)
		if err != nil {
			fmt.Println("err1")
			panic(err)
		}
		// Convert byte data to a type of io reader. Needed to be passed in request
		body := bytes.NewReader(payloadBytes)

		// POST request to add new dummy IP (converted to a byte io.reader)
		req, err := http.NewRequest("POST", "http://localhost:3000/addIPtoPool", body)
		if err != nil {
			fmt.Println("err2")
			panic(err)
		}

		//	Sends the req and returns a response
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("err3")
			panic(err)
		}
		defer resp.Body.Close()
	case 5:
		// Creating dummy putIP meant for changing existing IPs
		data := putIPpost{
			TargetIP:  "a-253.14.93.192",
			IPaddress: "a-111.11.11.111",
			Detail: putIPdetails{
				MACaddress: "A1-A2-A3-A4-A5-A6",
				LeaseTime:  time.Now(),
				Available:  true,
			},
		}

		fmt.Println("PUT new IP a-111.11.11.111")

		// Converts struct data to JSON byte data
		payloadBytes, err := json.Marshal(data)
		if err != nil {
			fmt.Println("err1")
			panic(err)
		}
		// Convert byte data to a type of io reader. Needed to be passed in request
		body := bytes.NewReader(payloadBytes)

		// PUT request to add new dummy putIP (converted to a byte io.reader)
		req, err := http.NewRequest("PUT", "http://localhost:3000/createNewIPpool", body)
		if err != nil {
			fmt.Println("err2")
			panic(err)
		}

		//	Sends the req and returns a response
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("err3")
			panic(err)
		}
		defer resp.Body.Close()
	}
}
