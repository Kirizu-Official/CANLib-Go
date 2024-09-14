package canlib

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

//struct gs_host_frame {
//u32 echo_id;
//__le32 can_id;
//
//u8 can_dlc;
//u8 channel;
//u8 flags;
//u8 reserved;
//
//union {
//DECLARE_FLEX_ARRAY(struct classic_can, classic_can);
//DECLARE_FLEX_ARRAY(struct classic_can_ts, classic_can_ts);
//DECLARE_FLEX_ARRAY(struct classic_can_quirk, classic_can_quirk);
//DECLARE_FLEX_ARRAY(struct canfd, canfd);
//DECLARE_FLEX_ARRAY(struct canfd_ts, canfd_ts);
//DECLARE_FLEX_ARRAY(struct canfd_quirk, canfd_quirk);
//};
//} __packed;

type GsHostFrame struct {
	EchoID   uint32
	CanID    uint32
	CanDlc   byte
	Channel  byte
	Flags    byte
	Reserved byte

	// classic_can
	Data [8]byte
}

// #define GS_CAN_FLAG_OVERFLOW						(1<<0)
// #define GS_CAN_FLAG_FD							(1<<1) /* is a CAN-FD frame */
// #define GS_CAN_FLAG_BRS							(1<<2) /* bit rate switch (for CAN-FD frames) */
// #define GS_CAN_FLAG_ESI							(1<<3) /* error state indicator (for CAN-FD frames) */

type GsHostFrameFlags struct {
	OverFlow bool
	FD       bool
	BRS      bool
	ESI      bool
}

func PackFrameFlag(flags GsHostFrameFlags) byte {
	flag := byte(0)
	if flags.OverFlow {
		flag = flag | 1<<0
	}
	if flags.FD {
		flag = flag | 1<<1
	}
	if flags.BRS {
		flag = flag | 1<<2
	}
	if flags.ESI {
		flag = flag | 1<<3
	}
	return flag
}
func UnpackFrameFlag(flags byte) GsHostFrameFlags {
	flag := GsHostFrameFlags{}
	if flags&0b1 > 0 {
		flag.OverFlow = true
	}
	if flags&0b10 > 0 {
		flag.FD = true
	}
	if flags&0b100 > 0 {
		flag.BRS = true
	}
	if flags&0b1000 > 0 {
		flag.ESI = true
	}
	return flag
}

func PackNewFrame(frame *GsHostFrame) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, frame.EchoID)
	binary.Write(buf, binary.LittleEndian, frame.CanID)
	buf.WriteByte(frame.CanDlc)
	buf.WriteByte(frame.Channel)
	buf.WriteByte(frame.Flags)
	buf.WriteByte(frame.Reserved)
	buf.Write(frame.Data[:])
	return buf.Bytes()
}

func UnpackFrame(data []byte) *GsHostFrame {
	if len(data) > 20 {
		data = data[:20]
	}
	if len(data) < 20 {
		data = append(data, make([]byte, 20-len(data))...)
	}

	return *(**GsHostFrame)(unsafe.Pointer(&data))
}
