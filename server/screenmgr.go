package server

import (
	"github.com/cooper/screenmgr/device"
	"log"
)

func Run() {
	checkError("Find devices", device.FindDevices())
	checkError("Setup devices", device.SetupDevices())
	checkError("Run HTTP server", runHTTPServer())
}

func reportError(action string, err error) {
	if err != nil {
		log.Println(action, " error: ", err.Error())
	}
}

func checkError(action string, err error) {
	if err != nil {
		log.Fatal(action, " error: ", err.Error())
	}
}
