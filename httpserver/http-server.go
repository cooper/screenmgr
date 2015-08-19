package httpserver

import (
	"github.com/cooper/screenmgr/device"
	"html/template"
	"net/http"
)

var templates *template.Template

func Run() error {

	// initialize templates
	templates = template.Must(template.ParseGlob("server/templates/*"))

	// serve static files in devices directory
	// FIXME: this allows access to vnc passwd files
	http.Handle("/devices/", http.StripPrefix("/devices/",
		http.FileServer(http.Dir("devices"))))

	// serve static files in resources directory
	http.Handle("/resources/", http.StripPrefix("/resources/",
		http.FileServer(http.Dir("resources"))))

	// main handler
	http.HandleFunc("/", httpHandler)

	return http.ListenAndServe(":8080", nil)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// use all devices
	devices := device.AllDevices()
	devs := make([]*device.Device, 0, len(devices))
	for _, dev := range devices {
		devs = append(devs, dev)
	}

	// serve template
	page := &DevicePage{devs}
	templates.ExecuteTemplate(w, "device-page.html", page)

}
