package mideaac

import (
	"encoding/binary"
	"fmt"
	"time"
)

type PacketBuilder struct {
	Command  []byte
	Security security
	Packet   []byte
}

func NewPacketBuilder(deviceID uint64) PacketBuilder {
	p := PacketBuilder{}
	p.Security = NewSecurity()
	p.Packet = []byte{
		// # 2 bytes - StaicHeader
		0x5a, 0x5a,
		// # 2 bytes - mMessageType
		0x01, 0x11,
		// # 2 bytes - PacketLenght
		0x00, 0x00,
		// # 2 bytes
		0x20, 0x00,
		// # 4 bytes - MessageId
		0x00, 0x00, 0x00, 0x00,
		// # 8 bytes - Date&Time
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// # 6 bytes - mDeviceID
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// # 12 bytes
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	copy(p.Packet[12:20], p.PacketTime())
	// p.Packet[20:28] = device_id.to_bytes(8, 'little')
	binary.LittleEndian.PutUint64(p.Packet[20:28], deviceID)
	return p
}

func (p *PacketBuilder) SetCommand(command BaseCommand) {
	p.Command = command.Finalize()
}

func (p *PacketBuilder) Finalize() []byte {
	p.Packet = append(p.Packet, p.Security.AesEncrypt(p.Command)[:48]...)
	binary.LittleEndian.PutUint16(p.Packet[4:6], uint16(len(p.Packet))+16)
	p.Packet = append(p.Packet, p.Encode32(p.Packet)...)
	return p.Packet
}

func (p PacketBuilder) Encode32(data []byte) []byte {
	b := p.Security.Encode32Data(data)
	return b[:]
}

func (p PacketBuilder) Checksum(data []byte) byte {
	return (^sum(data) + 1) & 0xFF
}

func (p PacketBuilder) PacketTime() []byte {
	now := time.Now()
	t := fmt.Sprintf("%d%02d%02d%02d%02d%06d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000)
	t = t[:16]
	// t := fmt.Format("%Y%m%d%H%m%S%f")[:16]
	var b []byte
	for i := 0; i < len(t); i += 2 {
		d := t[i : i+2]
		b = append([]byte{d[0]}, b...)
	}
	return b
}
