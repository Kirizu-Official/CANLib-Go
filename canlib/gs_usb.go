package canlib

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/google/gousb"
	"unsafe"
)

const (
	GsUsbBreqHostFormat = iota
	GsUsbBreqBittiming
	GsUsbBreqMode
	GS_USB_BREQ_BERR
	GsUsbBreqBtConst
	GsUsbBreqDeviceConfig
	GsUsbBreqTimestamp
	GsUsbBreqIdentify
	GsUsbBreqGetUserID
	GsUsbBreqSetUserID
	GsUsbBreqDataBittiming
	GsUsbBreqBtConstExt
	GsUsbBreqSetTermination
	GsUsbBreqGetTermination
	GsUsbBreqGetState
)

const GsCanModeReset = 0
const GsCanModeStart = 1

type Control struct {
	Device      *gousb.Device
	Version     ControlVersion
	FClk        int
	InEndpoint  *gousb.InEndpoint
	OutEndpoint *gousb.OutEndpoint
	bitrateSet  bool
}
type ControlVersion struct {
	HwVersion uint32
	SwVersion uint32
}

type GsCanModeFlags struct {
	ListenOnly          bool
	LoopBack            bool
	TripleSample        bool //CandleLight 没有实现
	OneShot             bool //开启Oneshot禁止报文重传
	HwTimeStamp         bool
	PadPktsToMaxPktSize bool
	FD                  bool
	BerrReporting       bool //CandleLight 没有实现
}

func (c *Control) SetDeviceMode(mode uint32, flags GsCanModeFlags) error {
	flagData := uint32(0)
	if flags.ListenOnly {
		flagData = flagData | 1
	}
	if flags.LoopBack {
		flagData = flagData | (1 << 1)
	}

	if flags.TripleSample {
		flagData = flagData | (1 << 2)
	}
	if flags.OneShot {
		flagData = flagData | (1 << 3)
	}
	if flags.HwTimeStamp {
		flagData = flagData | (1 << 4)
	}
	if flags.PadPktsToMaxPktSize {
		flagData = flagData | (1 << 7)
	}
	if flags.FD {
		flagData = flagData | (1 << 8)
	}
	if flags.BerrReporting {
		flagData = flagData | (1 << 12)
	}

	data := bytes.NewBuffer(nil)
	err := binary.Write(data, binary.LittleEndian, mode)
	if err != nil {
		return err
	}
	err = binary.Write(data, binary.LittleEndian, flagData)
	if err != nil {
		return err
	}

	bytesResult := data.Bytes()

	_, err = c.Device.Control(gousb.ControlVendor, GsUsbBreqMode, 0, 0, bytesResult)
	//fmt.Println("set mod", control, err, bytesResult)
	//c.Device.Close()
	return err
}

func (c *Control) GetBreqDeviceConfig() error {
	data := make([]byte, 12)
	control, err := c.Device.Control(gousb.ControlVendor|gousb.ControlIn, GsUsbBreqDeviceConfig, 0, 0, data)
	if err != nil {
		return err
	}
	if control < 12 {
		return errors.New("can not get device config")
	}
	//fmt.Println("get breq device config", control, err, data)
	//fmt.Println(data[4:7], data[8:11])
	c.Version.SwVersion = binary.LittleEndian.Uint32(data[4:8])
	c.Version.HwVersion = binary.LittleEndian.Uint32(data[8:12])
	return nil
}

type GsUsbBreqBtConstInfo struct {
	Feature uint32
	FclkCan uint32
	Btc     GsUsbBreqBtConstInfoCanBittimingConst
}

type GsUsbBreqBtConstInfoCanBittimingConst struct {
	Tseg1Min uint32
	Tseg1Max uint32
	Tseg2Min uint32
	Tseg2Max uint32
	SjwMax   uint32
	BrpMin   uint32
	BrpMax   uint32
	BrpInc   uint32
}

func (c *Control) GetGsUsbBreqBtConst() (*GsUsbBreqBtConstInfo, error) {
	data := make([]byte, 40)
	control, err := c.Device.Control(gousb.ControlVendor|gousb.ControlIn, GsUsbBreqBtConst, 0, 0, data)
	if err != nil {
		return nil, err
	}
	if control < 40 {
		return nil, errors.New("can not get device config")
	}
	info := *(**GsUsbBreqBtConstInfo)(unsafe.Pointer(&data))
	c.FClk = int(info.FclkCan)
	return info, nil
}

func (c *Control) GetGsUsbBreqTimestamp() (uint32, error) {
	data := make([]byte, 4)
	control, err := c.Device.Control(gousb.ControlVendor|gousb.ControlIn, GsUsbBreqTimestamp, 0, 0, data)
	if err != nil {
		return 0, err
	}
	if control < 4 {
		return 0, errors.New("can not get device config")
	}
	return binary.LittleEndian.Uint32(data), nil
}

type GsUsbBreqBtConstExtInfo struct {
	Feature uint32
	FclkCan uint32
	Btc     GsUsbBreqBtConstInfoCanBittimingConst
	DBtc    GsUsbBreqBtConstInfoCanBittimingConst
}

func (c *Control) GetGsUsbBreqBtConstExt() (*GsUsbBreqBtConstExtInfo, error) {
	data := make([]byte, 72)
	control, err := c.Device.Control(gousb.ControlVendor|gousb.ControlIn, GsUsbBreqBtConstExt, 0, 0, data)
	if err != nil {
		return nil, err
	}
	if control < 4 {
		return nil, errors.New("can not get device config")
	}
	info := *(**GsUsbBreqBtConstExtInfo)(unsafe.Pointer(&data))
	return info, nil
}

func (c *Control) GetGsUsbBreqGetTermination() (uint32, error) {
	data := make([]byte, 4)
	control, err := c.Device.Control(gousb.ControlVendor|gousb.ControlIn, GsUsbBreqGetTermination, 0, 0, data)
	if err != nil {
		return 0, err
	}
	if control < 4 {
		return 0, errors.New("can not get device config")
	}
	return binary.LittleEndian.Uint32(data), nil
}

type GsUsbBreqBittimingInfo struct {
	PropSeg   uint32
	PhaseSeg1 uint32
	PhaseSeg2 uint32
	Sjw       uint32
	Brp       uint32
}

func (c *Control) SetGsUsbBreqBittiming(info GsUsbBreqBittimingInfo) error {
	data := bytes.NewBuffer(nil)
	err := binary.Write(data, binary.LittleEndian, info.PropSeg)
	if err != nil {
		return err
	}
	err = binary.Write(data, binary.LittleEndian, info.PhaseSeg1)
	if err != nil {
		return err
	}
	err = binary.Write(data, binary.LittleEndian, info.PhaseSeg2)
	if err != nil {
		return err
	}
	err = binary.Write(data, binary.LittleEndian, info.Sjw)
	if err != nil {
		return err
	}
	err = binary.Write(data, binary.LittleEndian, info.Brp)
	if err != nil {
		return err
	}
	bytesResult := data.Bytes()
	_, err = c.Device.Control(gousb.ControlVendor, GsUsbBreqBittiming, 0, 0, bytesResult)
	c.bitrateSet = true
	return err
}

type GsUsbBreqDeviceMode struct {
	Mode  uint32
	Flags uint32
}

func (c *Control) SetGsUsbBreqMode(mode uint32, flags uint32) error {
	data := bytes.NewBuffer(nil)
	err := binary.Write(data, binary.LittleEndian, mode)
	if err != nil {
		return err
	}
	err = binary.Write(data, binary.LittleEndian, flags)
	if err != nil {
		return err
	}
	bytesResult := data.Bytes()
	_, err = c.Device.Control(gousb.ControlVendor, GsUsbBreqMode, 0, 0, bytesResult)
	return err
}

func (c *Control) SetGsUsbBreqIdentify(mode uint32) error {
	data := bytes.NewBuffer(nil)
	err := binary.Write(data, binary.LittleEndian, mode)
	if err != nil {
		return err
	}
	bytesResult := data.Bytes()
	_, err = c.Device.Control(gousb.ControlVendor, GsUsbBreqIdentify, 0, 0, bytesResult)
	return err
}

func (c *Control) SetGsUsbBreqDataBittiming(mode GsUsbBreqBittimingInfo) error {
	data := bytes.NewBuffer(nil)
	err := binary.Write(data, binary.LittleEndian, mode.PropSeg)
	if err != nil {
		return err
	}
	err = binary.Write(data, binary.LittleEndian, mode.PhaseSeg1)
	if err != nil {
		return err
	}
	err = binary.Write(data, binary.LittleEndian, mode.PhaseSeg2)
	if err != nil {
		return err
	}
	err = binary.Write(data, binary.LittleEndian, mode.Sjw)
	if err != nil {
		return err
	}
	err = binary.Write(data, binary.LittleEndian, mode.Brp)
	if err != nil {
		return err
	}
	bytesResult := data.Bytes()
	_, err = c.Device.Control(gousb.ControlVendor, GsUsbBreqDataBittiming, 0, 0, bytesResult)
	return err
}

func (c *Control) SetGsUsbBreqSetTermination(state uint32) error {
	data := bytes.NewBuffer(nil)
	err := binary.Write(data, binary.LittleEndian, state)
	if err != nil {
		return err
	}
	bytesResult := data.Bytes()
	_, err = c.Device.Control(gousb.ControlVendor, GsUsbBreqSetTermination, 0, 0, bytesResult)
	return err
}
