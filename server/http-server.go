package server

import (
	"fmt"
	"net/http"
)

func runHTTPServer() error {

	// server static files in devices directory
	// FIXME: this allows access to vnc passwd files
	http.Handle("/devices/", http.StripPrefix("/devices/",
		http.FileServer(http.Dir("devices"))))

	// main handler
	http.HandleFunc("/", httpHandler)

	return http.ListenAndServe(":8080", nil)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	for _, dev := range devices {
		w.Write([]byte(fmt.Sprintf("<pre>%+v</pre><br>", dev)))
		w.Write([]byte(fmt.Sprintf(`<img src="devices/%s/screenshots/%s" />`, dev.deviceID, dev.lastScreenshot)))
	}
}
