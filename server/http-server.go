package server

import (
	"net/http"
	"html/template"
)

var templates *template.Template

func runHTTPServer() error {

	// initialize templates
	templates = template.Must(template.ParseGlob("server/templates/*"))

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

	// find online devices
	devs := make([]*Device, 0, len(devices))
	for _, dev := range devices {
		if dev.Online {
			devs = append(devs, dev)
		}
	}

	// serve template
	page := &DevicePage{devs}
	checkError("template", templates.ExecuteTemplate(w, "device-page.html", page))

}
