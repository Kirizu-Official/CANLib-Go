package canlib

import (
	"context"
	"errors"
	"github.com/google/gousb"
	"math"
	"sync"
	"time"
)

type CanUSB struct {
	CanControl Control

	//context
	ctx       context.Context
	ctxCancel context.CancelFunc

	// statistics
	readNum    int
	writeNum   int
	readBytes  int
	writeBytes int

	// read stream
	ReadSteam      *gousb.ReadStream
	readData       chan []byte
	readDataCtx    context.Context
	readDataCancel context.CancelFunc
	readDataLock   sync.Mutex
	readCallBack   func(*GsHostFrame)

	// write
	WriteSteam *gousb.WriteStream
	WriteLock  sync.Mutex
}

func New(ctx context.Context, d *gousb.Device) (*CanUSB, error) {
	d.ControlTimeout = time.Millisecond * 100
	defaultInterface, err := d.Config(1)
	if err != nil {
		return nil, err
	}
	endpoint, err := defaultInterface.Interface(0, 0)
	if err != nil {
		return nil, err
	}
	inEndpoint, err := endpoint.InEndpoint(1)
	if err != nil {
		return nil, err
	}
	outEndpoint, err := endpoint.OutEndpoint(2)
	if err != nil {
		return nil, err
	}
	can := &CanUSB{}
	can.ctx, can.ctxCancel = context.WithCancel(ctx)
	can.CanControl = Control{
		Device:      d,
		InEndpoint:  inEndpoint,
		OutEndpoint: outEndpoint,
	}

	err = can.CanControl.GetBreqDeviceConfig()
	if err != nil {
		return nil, err
	}
	_, err = can.CanControl.GetGsUsbBreqBtConst()
	if err != nil {
		return nil, err
	}
	return can, nil
}
func (c *CanUSB) SetBitrate(bitrate int) error {

	var info GsUsbBreqBittimingInfo
	info.PropSeg = 1
	info.Sjw = 1

	if c.CanControl.FClk == 48000000 {
		info.PhaseSeg1 = 12
		info.PhaseSeg2 = 2

		switch bitrate {
		case 10000:
			info.Brp = 300
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 20000:
			info.Brp = 150
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 50000:
			info.Brp = 60
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 83333:
			info.Brp = 36
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 100000:
			info.Brp = 30
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 125000:
			info.Brp = 24
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 250000:
			info.Brp = 12
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 500000:
			info.Brp = 6
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 800000:
			info.Brp = 4
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 1000000:
			info.Brp = 3
			return c.CanControl.SetGsUsbBreqBittiming(info)
		default:
			return errors.New("bitrate can not support")
		}
	}

	if c.CanControl.FClk == 80000000 {
		info.PhaseSeg1 = 12
		info.PhaseSeg2 = 2
		switch bitrate {
		case 10000:
			info.PhaseSeg1 = 12
			info.PhaseSeg2 = 2
			info.Brp = 500
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 20000:
			info.PhaseSeg1 = 12
			info.PhaseSeg2 = 2
			info.Brp = 250
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 50000:
			info.PhaseSeg1 = 12
			info.PhaseSeg2 = 2
			info.Brp = 100
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 83333:
			info.PhaseSeg1 = 12
			info.PhaseSeg2 = 2
			info.Brp = 60
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 100000:
			info.PhaseSeg1 = 12
			info.PhaseSeg2 = 2
			info.Brp = 50
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 125000:
			info.PhaseSeg1 = 12
			info.PhaseSeg2 = 2
			info.Brp = 40
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 250000:
			info.PhaseSeg1 = 12
			info.PhaseSeg2 = 2
			info.Brp = 20
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 500000:
			info.PhaseSeg1 = 12
			info.PhaseSeg2 = 2
			info.Brp = 10
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 800000:
			info.PhaseSeg1 = 7
			info.PhaseSeg2 = 1
			info.Brp = 10
			return c.CanControl.SetGsUsbBreqBittiming(info)
		case 1000000:
			info.PhaseSeg1 = 12
			info.PhaseSeg2 = 2
			info.Brp = 5
			return c.CanControl.SetGsUsbBreqBittiming(info)
		default:
			return errors.New("bitrate can not support")
		}
	}
	return errors.New("must GetBreqDeviceConfig at first")
}
func (c *CanUSB) InitDevice(bitrate int, deviceFlag GsCanModeFlags, readDataCallBack func(data *GsHostFrame)) error {
	if !c.WriteLock.TryLock() {
		return errors.New("is already init")
	}
	err := c.CanControl.SetDeviceMode(GsCanModeReset, GsCanModeFlags{})
	if err != nil {
		return err
	}
	err = c.SetBitrate(bitrate)
	if err != nil {
		return err
	}
	err = c.CanControl.SetDeviceMode(GsCanModeStart, deviceFlag)
	if err != nil {
		return err
	}

	err = c.startReadSteam(readDataCallBack)
	if err != nil {
		return err
	}

	return nil
}
func (c *CanUSB) startReadSteam(readCallBack func(data *GsHostFrame)) error {
	if !c.CanControl.bitrateSet {
		return errors.New("canusb is not init")
	}
	var err error
	c.ReadSteam, err = c.CanControl.InEndpoint.NewStream(20, 10)
	if err != nil {
		return err
	}
	c.WriteSteam, err = c.CanControl.OutEndpoint.NewStream(20, 10)
	if err != nil {
		return err
	}

	c.readCallBack = readCallBack
	c.readData = make(chan []byte)

	go c.readProcess()
	go func() {

		time.Sleep(time.Second * 1)
		//情况缓冲区数据
		for {
			select {
			case <-c.readData:
				continue
			default:
				goto out
			}
		}
	out:
		c.newBusRead()
		c.WriteLock.Unlock()
	}()

	return nil
}
func (c *CanUSB) newBusRead() {
	c.readDataCtx, c.readDataCancel = context.WithCancel(c.ctx)
	go c.canBusReadData()
}

func (c *CanUSB) readProcess() {
	buf := make([]byte, 24)
	for {

		read, err := c.ReadSteam.ReadContext(c.ctx, buf)
		if err != nil {
			c.ctxCancel()
			return
		} else {
			c.readData <- buf[:read]
			c.readNum++
		}
	}
}

func (c *CanUSB) canBusReadData() {
	c.readDataLock.Lock()
	select {
	case <-c.readDataCtx.Done():
	case data := <-c.readData:
	processData:
		go c.readCallBack(UnpackFrame(data))
		go c.canBusReadData()
		if len(c.readData) > 0 {
			data = <-c.readData
			goto processData
		}
	}
	c.readDataLock.Unlock()
}

func (c *CanUSB) WriteAndReadSimpleData(canID uint32, data [8]byte, timeout time.Duration) (respID uint32, respData []byte, err error) {
	var read *GsHostFrame
	read, err = c.WriteData(GsHostFrame{
		EchoID:   math.MaxUint32,
		CanID:    canID,
		CanDlc:   1,
		Channel:  0,
		Flags:    0,
		Reserved: 0,
		Data:     data,
	}, timeout, true)
	if err != nil {
		return 0, nil, err
	}
	return read.CanID, read.Data[:], nil
}

func (c *CanUSB) WriteData(data GsHostFrame, timeout time.Duration, read bool) (*GsHostFrame, error) {
	c.WriteLock.Lock()
	defer c.WriteLock.Unlock()
	c.readDataCancel()

	c.readDataLock.Lock()
	defer c.readDataLock.Unlock()

	defer c.newBusRead()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := c.WriteSteam.WriteContext(ctx, PackNewFrame(data))
	if err != nil {
		return nil, err
	}

	//loopBack
	select {
	case <-ctx.Done():
		return nil, errors.New("read loopback data timeout")
	case <-c.readData:
		break
	}

	if read {
		//data
		select {
		case <-ctx.Done():
			return nil, errors.New("read response data timeout")
		case readData := <-c.readData:
			return UnpackFrame(readData), nil
		}
	}

	return nil, nil
}
