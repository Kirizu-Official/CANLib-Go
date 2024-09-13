package main

import (
	"context"
	"fmt"
	"github.com/Kirizu-Official/CANLib-Go/canlib"
	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"time"
)

func main() {

	usb := gousb.NewContext()
	defer usb.Close()
	//Note: ctx.OpenDevices will return an error, but this error is meaningless, even if there is an available device will return an error,
	//PLEASE CHECK the length of the device to determine whether there is an available device
	devices, err := usb.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		if desc.Vendor == 0x1d50 && desc.Product == 0x606f {
			fmt.Printf("Find Device: %03d.%03d %s:%s %s\n", desc.Bus, desc.Address, desc.Vendor, desc.Product, usbid.Describe(desc))
			return true
		} else {
			return false
		}
	})

	if len(devices) == 0 {
		fmt.Println("No Device Found")
		panic(err)
	}

	device := devices[0]
	//Use CANLib
	ctx, cancle := context.WithCancel(context.Background())
	//If you no longer use it or the program exits, you should cancel the context to release the goroutine in CANLib
	defer cancle()
	defer device.Close()
	defer usb.Close()
	can, err := canlib.New(ctx, device)
	if err != nil {
		panic(err)
	}

	canFlag := canlib.GsCanModeFlags{
		ListenOnly:          false,
		LoopBack:            false,
		TripleSample:        false,
		OneShot:             false,
		HwTimeStamp:         false,
		PadPktsToMaxPktSize: false,
		FD:                  false,
		BerrReporting:       false,
	}
	err = can.InitDevice(1000000, canFlag, readCallBack)
	if err != nil {
		panic(err)
	}

	//time.Sleep(time.Second * 2)
	respCanID, respData, err := can.WriteAndReadSimpleData(0x3f0, [8]byte{0x00}, time.Millisecond*500)
	if err != nil {
		panic(err)
	}
	fmt.Printf("respCanID: %02X, respData: %02X\n", respCanID, respData)

	// OutPut:
	/**

	$ go run ./example/basicuse/main.go

	Find Device: 001.063 1d50:606f Unknown (OpenMoko, Inc.)
	respCanID: 3F1, respData: 20CA06CC88A76201

	*/

}

func readCallBack(data *canlib.GsHostFrame) {
	fmt.Println("read call back", data.Data, data.CanID)
}
