package server

import "net/http"
import "fmt"

func runHTTPServer() error {
	http.HandleFunc("/", httpHandler)
	return http.ListenAndServe(":8080", nil)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	for _, dev := range devices {
		w.Write([]byte(fmt.Sprintf("%+v\n", dev)))
	}
}
