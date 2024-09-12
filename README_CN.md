# CANLib-Go

[English](README.md)

## 简介
这个一个基于 [gousb](https://github.com/google/gousb) 实现的 CAN 协议通信库，CANLib-Go 直接和 USB 硬件进行最底层的通信，你可以使用它通过 USB - CAN 实现对 CAN 总线的数据接收和发送，同时 CANLib-Go 开放非常底层的接口，你可以自定义很多更底层的内容。

## 支持
系统：支持Windows、Linux、MacOS
固件：理论上支持所有由Linux定义的 gs_usb 的 usb - can 固件

但是目前仅在 Windows 上测试了 [candleLight](https://github.com/candle-usb/candleLight_fw) 固件，其他固件欢迎大家协助测试并在 [issue](https://github.com/Kirizu-Official/CANLib-Go/issues) 中进行反馈。


## 使用
CANLib-Go 是基于 [gousb](https://github.com/google/gousb) 开发的，不过 [gousb](https://github.com/google/gousb) 是通过 CGO 调用 [libusb](https://github.com/libusb/libusb) 实现的，所以你需要先按照 [gousb安装文档](https://github.com/google/gousb?tab=readme-ov-file#installation) 进行环境配置后才能使用这个库。

---

这个库和文档依然还在开发中，欢迎协助完善这个仓库！
