package server

import (
	"log"
	"net/http"
)

func runHTTPServer() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
