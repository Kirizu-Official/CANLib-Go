# CANLib-Go

[简体中文](README_CN.md)


## What this?
This is a CAN protocol implementation library based on [gousb](https://github.com/google/gousb).

You can access can data in Go via USB, including sending and receiving data, and this library provides the lowest-level operations for easy expansion


## Support
System: Windows、Linux、MacOS

Firmware: All standard [gs_usb](https://github.com/torvalds/linux/blob/master/drivers/net/can/usb/gs_usb.c) protocols defined by Linux, However, only tested on [candleLight](https://github.com/candle-usb/candleLight_fw) firmware. 


## Usage
CANLib-Go is based on the [gousb](https://github.com/google/gousb) library, while the gousb is implemented by calling [libusb](https://github.com/libusb/libusb) through CGO, You must install the runtime and compilation environment as required by [gousb installation document](https://github.com/google/gousb?tab=readme-ov-file#installation) at first.

---

This document is still being improved and edited.

Join us to help improve this repository!

To be continued...
