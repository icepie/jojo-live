package mideaac

import (
	"math"
	"time"
)

type BaseCommand struct {
	data []byte
}

func NewBaseCommand(deviceType byte) BaseCommand {
	return BaseCommand{
		data: []byte{
			// 0 header
			0xaa,
			// 1 command lenght: N+10
			0x20,
			// 2 device type
			deviceType,
			//0xac,
			// 3 Frame SYN CheckSum
			0x00,
			// 4-5 Reserved
			0x00, 0x00,
			// 6 Message ID
			0x00,
			// 7 Frame Protocol Version
			0x00,
			// 8 Device Protocol Version
			0x00,
			// 9 Messgae Type: request is 0x03; setting is 0x02
			0x03,

			// Byte0 - Data request/response type: 0x41 - check status; 0x40 - Set up
			0x41,
			// Byte1
			0x81,
			// Byte2 - operational_mode
			0x00,
			// Byte3
			0xff,
			// Byte4
			0x03,
			// Byte5
			0xff,
			// Byte6
			0x00,
			// Byte7 - Room Temperature Request: 0x02 - indoor_temperature, 0x03 - outdoor_temperature
			// when set, this is swing_mode
			0x02,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			// Message ID
			byte(time.Now().Second()),
		},
	}
}

func (b BaseCommand) Checksum(data []byte) byte {
	return (^sum(data) + 1) & 0xFF
}

func (b *BaseCommand) Finalize() []byte {
	b.data = append(b.data, Calculate(b.data[10:]))
	b.data = append(b.data, b.Checksum(b.data[1:]))
	return b.data
}

type SetCommand struct {
	BaseCommand
}

func NewSetCommand(deviceType byte) SetCommand {
	b := NewBaseCommand(deviceType)
	b.data[0x01] = 0x23
	b.data[0x09] = 0x02
	b.data[0x0a] = 0x40
	b.data[0x0b] = 0x40
	return SetCommand{
		BaseCommand: b,
	}
}

func (s SetCommand) PromtTone() byte {
	return s.data[0x0B] & 0x42
}

func (s *SetCommand) SetPromtTone(feedbackEnabled bool) {
	s.data[0x0B] &= ^byte(0x42)
	if feedbackEnabled {
		s.data[0x0B] |= 0x42
	} else {
		s.data[0x0B] |= 0
	}
}

func (s SetCommand) PowerState() byte {
	return s.data[0x0B] & 0x01
}

func (s *SetCommand) SetPowerState(state bool) {
	s.data[0x0B] &= ^byte(0x01)
	if state {
		s.data[0x0B] |= 0x01
	} else {
		s.data[0x0B] |= 0
	}
}

func (s SetCommand) TargetTemperature() byte {
	return s.data[0x0C] & 0x1f
}

func (s *SetCommand) SetTargetTemperature(temp float64) {
	s.data[0x0C] &= ^byte(0x0f)
	// TODO ???
	s.data[0x0C] |= byte(temp) & 0xF
	s.SetTemperatureDot5(int(math.Round(temp*2))%2 != 0)
}

func (s SetCommand) OperationalMode() byte {
	return (s.data[0x0C] & 0xE0) >> 5
}

func (s *SetCommand) SetOperationalMode(mode byte) {
	s.data[0x0C] &= ^byte(0xE0)
	s.data[0x0C] |= (mode << 5) & 0xE0
}

func (s SetCommand) FanSpeed() byte {
	return s.data[0x0D]
}

func (s *SetCommand) SetFanSpeed(speed byte) {
	s.data[0x0D] = speed
}

func (s SetCommand) EcoMode() bool {
	return s.data[0x13] > 0
}

func (s *SetCommand) SetEcoMode(mode bool) {
	if mode {
		s.data[0x13] = 0xFF
	} else {
		s.data[0x13] = 0
	}
}

func (s SetCommand) SwingMode() byte {
	return s.data[0x11]
}

func (s *SetCommand) SetSwingMode(mode byte) {
	s.data[0x11] = 0x30
	s.data[0x11] |= mode & 0x3F
}

func (s SetCommand) TurboMode() bool {
	return s.data[0x14] > 0
}

func (s *SetCommand) SetTurboMode(b bool) {
	if b {
		s.data[0x14] |= 0x02
	} else {
		s.data[0x14] &= ^byte(0x02)
	}
}

func (s SetCommand) ScreenDisplay() bool {
	return s.data[0x14]&0x10 > 0
}

func (s *SetCommand) SetScreenDisplay(b bool) {
	if b {
		s.data[0x14] |= 0x10
	} else {
		s.data[0x14] &= ^byte(0x10)
	}
}

func (s SetCommand) TemperatureDot5() bool {
	return s.data[0x0C]*0x10 > 0
}

func (s *SetCommand) SetTemperatureDot5(b bool) {
	if b {
		s.data[0x0C] |= 0x10
	} else {
		s.data[0x0C] &= ^byte(0x10)
	}
}

func (s SetCommand) Fahrenheit() bool {
	return s.data[0x14]&0x04 > 0
}

func (s *SetCommand) SetFahrenheit(b bool) {
	if b {
		s.data[0x14] |= 0x04
	} else {
		s.data[0x14] &= ^byte(0x04)
	}
}

type ApplianceResponse struct {
	Data []byte
}

func NewApplianceResponse(data []byte) ApplianceResponse {
	return ApplianceResponse{
		Data: data[0xA:],
	}
}

func (a ApplianceResponse) PowerState() bool {
	return (a.Data[0x01] & 0x1) > 0
}

func (a ApplianceResponse) ImodeResume() bool {
	return a.Data[0x01]&0x10 > 0
}

func (a ApplianceResponse) ApplianceError() bool {
	return a.Data[0x01]&0x80 > 0
}

// Byte 0x02
func (a ApplianceResponse) TargetTemperature() float64 {
	if a.Data[0x02]&0x10 > 0 {
		return float64(a.Data[0x02]&0xF) + 16.0 + 0.5
	} else {
		return float64(a.Data[0x02]&0xF) + 16.0
	}
}

func (a ApplianceResponse) OperationalMode() byte {
	return (a.Data[0x02] & 0xE0) >> 5
}

func (a ApplianceResponse) FanSpeed() byte {
	return a.Data[0x03] & 0x7F
}

type Timer struct {
	status  bool
	hour    byte
	minutes byte
}

func (a ApplianceResponse) OnTimer() Timer {
	onTimerValue := a.Data[0x04]
	onTimerMinutes := a.Data[0x06]
	return Timer{
		status:  ((onTimerValue & 0x80) >> 7) > 0,
		hour:    (onTimerValue & 0x7c) >> 2,
		minutes: (onTimerValue & 0x3) | ((onTimerMinutes & 0xF0) >> 4),
	}
}

func (a ApplianceResponse) OffTimer() Timer {
	offTimerValue := a.Data[0x05]
	offTimerMinutes := a.Data[0x06]
	return Timer{
		status:  ((offTimerValue & 0x80) >> 7) > 0,
		hour:    (offTimerValue & 0x7c) >> 2,
		minutes: (offTimerValue & 0x3) | (offTimerMinutes & 0xF0),
	}
}

func (a ApplianceResponse) SwingMode() byte {
	return a.Data[0x07] & 0x0F
}

func (a ApplianceResponse) EcoMode() bool {
	return a.Data[0x09]&0x10 > 0
}

func (a ApplianceResponse) TurboMode() bool {
	return a.Data[0x0A]&0x02 > 0
}

func (a ApplianceResponse) IndoorTemperature() float64 {
	var indoorTempInteger int
	var indoorTempDecimal float64
	if a.Data[0] == 0xC0 {
		if (int(a.Data[11]-50)/2) < -19 || (int(a.Data[11]-50)/2) > 50 {
			return 0xFF
		} else {
			indoorTempInteger = int(a.Data[11]-50) / 2
		}
		indoorTemperatureDot := getBits4(a.Data, 15, 0, 3)
		indoorTempDecimal = float64(indoorTemperatureDot) * 0.1
		if a.Data[11] > 49 {
			return float64(indoorTempInteger) + indoorTempDecimal
		} else {
			return float64(indoorTempInteger) - indoorTempDecimal
		}
	}
	if a.Data[0] == 0xA0 || a.Data[0] == 0xA1 {
		if a.Data[0] == 0xA0 {
			if (a.Data[1]>>2)-4 == 0 {
				indoorTempInteger = -1
			} else {
				indoorTempInteger = int(a.Data[1]>>2) + 12
			}
			if (a.Data[1]>>1)&0x01 == 1 {
				indoorTempDecimal = 0
			}
		}
		if a.Data[0] == 0xA1 {
			if (int(a.Data[13]-50)/2) < -19 || (int(a.Data[13]-15)/2) > 50 {
				return 0xFF
			} else {
				indoorTempInteger = (int(a.Data[13]-50) / 2)
			}
			indoorTempDecimal = float64(a.Data[18]&0x0F) * 0.1
		}
		if a.Data[13] > 49 {
			return float64(indoorTempInteger) + indoorTempDecimal
		} else {
			return float64(indoorTempInteger) - indoorTempDecimal
		}
	}
	return 0xFF
}

func (a ApplianceResponse) OutdoorTemperature() float64 {
	return float64(a.Data[0x0c]-50) / 2
}

func getBits(pByte, pIndex byte) byte {
	return (pByte >> pIndex) & 0x01
}

func getBits4(pBytes []byte, pIndex, pStartIndex, pEndIndex byte) byte {
	var startIndex byte
	var endIndex byte
	if pStartIndex > pEndIndex {
		startIndex = pEndIndex
		endIndex = pStartIndex
	} else {
		startIndex = pStartIndex
		endIndex = pEndIndex
	}
	tempVal := 0x00
	for i := startIndex; i <= endIndex; i++ {
		tempVal = tempVal | int(getBits(pBytes[pIndex], i))<<(i-startIndex)
	}
	return byte(tempVal)
}

func sum(array []byte) byte {
	var result byte
	for _, v := range array {
		result += v
	}
	return result
}
