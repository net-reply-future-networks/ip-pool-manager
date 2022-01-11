package handlers

import (
	"log"
	"net/http"
)

func Healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DEBUG: Response is:  v%",  http.StatusOK)
		w.WriteHeader(http.StatusOK)
	}
}
