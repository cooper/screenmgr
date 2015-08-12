package server

import "io/ioutil"
import "fmt"

var devices map[string]device

func init() {
	devices = make(map[string]device)
}

func findDevices() error {
	files, err := ioutil.ReadDir("devices")

	// error
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		fmt.Printf("found file: %+v\n", fileInfo)
		dev := deviceWithID(fileInfo.Name())
		checkError("Read JSON", dev.readInfo())
		fmt.Printf("device: %+v\n", dev)
		fmt.Printf("info: %+v\n", dev.info)
	}

	return nil
}
