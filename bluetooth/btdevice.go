package bluetooth

import (
	"time"
	"fmt"
	"errors"

	"golang.org/x/net/context"
	"github.com/currantlabs/ble"
)

type BTDevice struct {
	address ble.Addr
	client  ble.Client
	profile *ble.Profile
}

var (
	Timeout = 10 * time.Second

	ErrNotConnected = errors.New("BTDevice is not connected.")
	ErrUUID         = errors.New("BTDevice error while searching for UUID ")
)

func NewDevice(address string) BTDevice {
	var device BTDevice
	device.address = ble.NewAddr(address)

	return device
}

func (d *BTDevice) connect() error {
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), Timeout))

	client, err := ble.Dial(ctx, d.address)
	if err != nil {
		return err
	}
	d.client = client

	go func() {
		<-d.client.Disconnected()
		fmt.Printf("\n%s disconnected\n", client.Address().String())
		d.client = nil
		d.profile = nil
	}()

	return nil
}

func (d *BTDevice) discover() error {
	if d.client == nil {
		return ErrNotConnected
	}

	p, err := d.client.DiscoverProfile(true)
	if err != nil {
		return err
	}

	d.profile = p
	return nil
}

func (d *BTDevice) Connect() error {
	var err error

	if d.Connected() {
		return nil
	}

	fmt.Println("Connecting")
	err = d.connect()
	if err != nil {
		// todo: double check connected
		if err.Error() != "can't dial: ACL Connection Already Exists" {
			return err
		}
	}
	fmt.Println("Connected")

	err = d.discover()
	if err != nil {
		return err
	}
	fmt.Println("Discovered")

	return nil
}

func (d *BTDevice) Disconnect() error {
	if d.client == nil {
		return ErrNotConnected
	}

	defer func() {
		d.client = nil
		d.profile = nil
	}()

	d.client.CancelConnection()
	time.Sleep(5 * time.Second)

	return nil
}

func (d *BTDevice) Read(uuid_string string) error {
	uuid, err := ble.Parse(uuid_string)
	if err != nil {
		return err
	}

	u := d.profile.Find(ble.NewCharacteristic(uuid))
	if u == nil {
		return ErrUUID
	}

	b, err := d.client.ReadCharacteristic(u.(*ble.Characteristic))
	if err != nil {
		return err
	}
	fmt.Printf("    Value         %x | %q\n", b, b)

	return nil
}

func (d *BTDevice) Write(uuid_string string, value []byte) error {
	uuid, err := ble.Parse(uuid_string)
	if err != nil {
		return err
	}

	u := d.profile.Find(ble.NewCharacteristic(uuid))
	if u == nil {
		return ErrUUID
	}

	err = d.client.WriteCharacteristic(u.(*ble.Characteristic), value, true)

	if err != nil {
		return err
	}

	return nil
}

func (d *BTDevice) Connected() bool {
	if d.client != nil && d.profile != nil {
		return true
	}
	return false
}
