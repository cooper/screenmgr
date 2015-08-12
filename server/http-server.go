package server

import "net/http"

func runHTTPServer() error {
	return http.ListenAndServe(":8080", nil)
}
