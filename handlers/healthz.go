package handlers

import (
	"fmt"
	"net/http"
)

func Healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-----------------------------")
		fmt.Println("Response is", http.StatusOK)
		fmt.Println("-----------------------------")
		w.WriteHeader(http.StatusOK)
	}
}
