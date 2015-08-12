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
    }

    return nil
}
