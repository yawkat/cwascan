package main

// adapted from https://gist.github.com/thetzel/398c5c504a4616732e78c991e2478e52

import (
	"fmt"
	"log"
	"os"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"

	"database/sql"

	_ "github.com/lib/pq"
)

func onStateChanged(device gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("Scanning for Rolling Proximity Identifiers")
		device.Scan([]gatt.UUID{}, true)
		return
	default:
		device.StopScanning()
	}
}

var db *sql.DB

func onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	if len(a.Services) > 0 {
		if gatt.UUID.Equal(a.Services[0], gatt.MustParseUUID("fd6f")) {
			if len(a.ServiceData) > 0 {
				log.Printf("Received beacon message\n")
				stmt, err := db.Prepare("insert into rpik (rpik, metadata, rssi) values ($1, $2, $3)")
				if err != nil {
					log.Printf("Failed to prepare, err: %s\n", err)
					return
				}
				_, err = stmt.Exec(a.ServiceData[0].Data[:16], a.ServiceData[0].Data[16:], rssi)
				if err != nil {
					log.Printf("Failed to execute, err: %s\n", err)
					return
				}
			}
		}
	}
}

func main() {
	var err error
	db, err = sql.Open("postgres", os.Args[1])
	if err != nil {
		log.Fatalf("Failed to connect to db, err: %s\n", err)
		return
	}

	device, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}
	device.Handle(gatt.PeripheralDiscovered(onPeripheralDiscovered))
	device.Init(onStateChanged)
	select {}
}
