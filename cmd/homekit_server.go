package main

import (
	"time"
	"os"
	"os/signal"

	"github.com/brutella/hc"
	"github.com/brutella/hc/log"

	"github.com/cpp/homegoing/homekit"
)

func main() {
	log.Debug.Enable()

	h := homekit.NewHomeKit()
	h.NewILEDBulb("LED01", "00:E5:C1:B1:00:AF")

	hc.OnTermination(func() {
		for _, light := range h.Lights {
			light.Stop()
		}

		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
	})

	var signal_channel chan os.Signal
	signal_channel = make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)

	<-signal_channel
}
