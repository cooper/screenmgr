package main

import (
	_ "github.com/cooper/screenmgr/agent"
	"github.com/cooper/screenmgr/device"
	"github.com/cooper/screenmgr/httpserver"
	_ "github.com/cooper/screenmgr/vnc"
	"log"
)

func main() {
	checkError("Find devices", device.FindDevices())
	checkError("Setup devices", device.SetupDevices())
	checkError("Run HTTP server", httpserver.Run())
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
