package server

import "log"

func Run() {
	checkError("Find devices", findDevices())
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
