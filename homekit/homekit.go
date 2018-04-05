package homekit

import (
	"fmt"
	"log"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux"
)

type HomeKit struct {
	pin    string
	Lights []*ILedBulb
}

func NewHomeKit() HomeKit {
	var h HomeKit

	// default device pin
	h.pin = "12344321"

	// Bluetooth Stick
	d, err := linux.NewDevice()
	if err != nil {
		log.Fatalf("Can't create linux device: %s", err)
	}

	ble.SetDefaultDevice(d)
	m := map[bool]string{true: "started", false: "error"}
	fmt.Printf("USB Device %s\n", m[d != nil])

	return h
}
