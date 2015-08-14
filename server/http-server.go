package server

import (
	"net/http"
	"fmt"
)

func runHTTPServer() error {

	// server static files in devices directory
	http.Handle("/devices/", http.StripPrefix("/devices/",
		http.FileServer(http.Dir("devices"))))

	// main handler
	http.HandleFunc("/", httpHandler)

	return http.ListenAndServe(":8080", nil)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	for _, dev := range devices {
		w.Write([]byte(fmt.Sprintf("%+v\n", dev)))
	}
}
