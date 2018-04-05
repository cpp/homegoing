package homekit

import (
	"fmt"
	"time"
	"encoding/hex"

	"github.com/brutella/hc"
	"github.com/brutella/hc/log"
	"github.com/brutella/hc/accessory"

	"github.com/cpp/homegoing/bluetooth"
)

const (
	BrightnessMax int = 255
)

type ILedBulb struct {
	name       string
	acc        *accessory.Lightbulb
	transport  hc.Transport
	bt         bluetooth.BTDevice
	Brightness int
}

func (h *HomeKit) NewILEDBulb(name string, address string) ILedBulb {
	var bulb ILedBulb
	h.Lights = append(h.Lights, &bulb)

	bulb.Brightness = 255
	bulb.name = name

	bulb.bt = bluetooth.NewDevice(address)

	info := accessory.Info{
		Name:         bulb.name,
		Manufacturer: "BEKEN",
		Model:        "iLedBulb",
	}

	bulb.acc = accessory.NewLightbulb(info)
	bulb.acc.Lightbulb.On.SetValue(true)

	config := hc.Config{Pin: h.pin}
	var err error
	bulb.transport, err = hc.NewIPTransport(config, bulb.acc.Accessory)
	if err != nil {
		log.Info.Panic(err)
	}

	hc.OnTermination(func() {
		<-bulb.transport.Stop()
	})

	go func() {
		bulb.transport.Start()
	}()

	bulb.acc.OnIdentify(func() {
		timeout := 1 * time.Second

		for i := 0; i < 4; i++ {
			//ToggleLight(light)
			//bulb.Toggle()
			time.Sleep(timeout)
		}
	})

	bulb.acc.Lightbulb.On.OnValueRemoteUpdate(func(power bool) {
		log.Debug.Printf("Changed State for iLedBulb %s to %t", bulb.name, power)
		bulb.Toggle(power)
	})

	bulb.acc.Lightbulb.Brightness.OnValueRemoteUpdate(func(value int) {
		log.Debug.Printf("Changed Brightness for %s to %d", bulb.name, value)

		// Set brightness on bulb
		err := bulb.SetBrightness(value)
		if err != nil {
			log.Debug.Printf(err.Error())
		}
	})

	return bulb
}
func (b *ILedBulb) Toggle(status bool) error {
	var err error

	err = b.bt.Connect()
	if err != nil {
		return err
	}

	if status {
		b.Brightness = 255
		err = b.bt.Write("ee03", []byte("\x01\x00\x01\x00\x01\x00\x01\x00\x01\xff"))
	} else {
		b.Brightness = 0
		err = b.bt.Write("ee03", []byte("\x01\x00\x01\x00\x01\x00\x01\x00\x01\x00"))
	}
	if err != nil {
		return err
	}
	fmt.Printf("Toggled %s to %t", b.name, status)

	return nil
}

// todo: generate hex command from value
func (b *ILedBulb) SetBrightness(value int) error {
	percent := float64(value) / 100
	fmt.Printf("percent: %f\n", percent)
	if percent > 1 {
		value = BrightnessMax
	} else {
		value = int(float64(BrightnessMax) * percent)
		fmt.Printf("value: %d\n", value)
	}

	//pre := "\x00\x00\x01\x00\x00\x00\x00\x00\x01"
	h := fmt.Sprintf("%02x", value)
	fmt.Printf("h: %q\n", value)
	command := intStr2HexStr("000001000000000001" + h)
	fmt.Printf("command: %s\n", command)

	err := b.bt.Connect()
	if err != nil {
		if err != nil {
			return err
		}
	}

	err = b.bt.Write("ee03", []byte(fmt.Sprintf("%q", command)))
	if err != nil {
		fmt.Println("command error")
		return err
	}
	fmt.Println("Wrote")

	err = b.bt.Read("ee01")
	if err != nil {
		return err
	}

	return nil
}

func (b *ILedBulb) Stop() {
	b.bt.Disconnect()
	b.transport.Stop()
}

func intStr2HexStr(str string) string {
	src := []byte(str)

	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		log.Debug.Fatal(err)
	}

	return fmt.Sprintf("%q", dst[:n])
}
