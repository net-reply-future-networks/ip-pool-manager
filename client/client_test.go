package main

import (
	"log"
	"testing"
)

func TestAllAvailbleIPs(t *testing.T) {
	requestOption = new(int)
	*requestOption = 1
	respMsg, respStatus := requestSelection(requestOption)

	log.Println("Resp Status code", respStatus)
	log.Println("Resp len ", len(respMsg))

	if len(respMsg) < 2 {
		t.Errorf("Returned JSON is empty (Char less than 2)")
	}

	if respStatus != "200 OK" {
		t.Errorf("Return status is not 200 OK")
	}
}

func TestGetIP(t *testing.T) {
	requestOption = new(int)
	*requestOption = 2
	respMsg, respStatus := requestSelection(requestOption)

	log.Println("Resp Status code", respStatus)
	log.Println("Resp len ", len(respMsg))

	if respStatus != "200 OK" {
		t.Errorf("Return status is not 200 OK")
	}

	if len(respMsg) < 2 {
		t.Errorf("Returned JSON is empty (Char less than 2)")
	}
}

func TestDeleteIPfromPool(t *testing.T) {
	requestOption = new(int)
	*requestOption = 3
	respMsg, respStatus := requestSelection(requestOption)

	log.Println("Resp Status code", respStatus)
	log.Println("Resp len ", len(respMsg))

	if respStatus != "200 OK" {
		t.Errorf("Return status is not 200 OK")
	}

	if respMsg != "a-102.131.46.22 IP deleted " {
		t.Errorf("Returned response message is incorrect.")
		t.Errorf(respMsg)
	}
}

func TestAddIPtoPool(t *testing.T) {
	requestOption = new(int)
	*requestOption = 4
	respMsg, respStatus := requestSelection(requestOption)

	log.Println("Resp Status code", respStatus)
	log.Println("Resp len ", len(respMsg))

	if respStatus != "200 OK" {
		t.Errorf("Return status is not 200 OK")
	}

	if respMsg != "New IP posted" {
		t.Errorf("Returned response message is incorrect.")
		t.Errorf(respMsg)
	}
}

func TestCreateNewIPpool(t *testing.T) {
	requestOption = new(int)
	*requestOption = 5
	respMsg, respStatus := requestSelection(requestOption)

	log.Println("Resp Status code", respStatus)
	log.Println("Resp len ", len(respMsg))

	if respStatus != "200 OK" {
		t.Errorf("Return status is not 200 OK")
	}

	if respMsg != "IP address a-253.14.93.192 changed to a-111.11.11.111" {
		t.Errorf("Returned response message is incorrect.")
		t.Errorf(respMsg)
	}
}
