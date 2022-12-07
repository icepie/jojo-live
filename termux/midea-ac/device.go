package mideaac

import (
	"encoding/binary"
	"encoding/hex"
	"log"
	"net"
)

const version = "0.2.1"

func convertDeviceIdHex(deviceID uint32) []byte {
	b := make([]byte, 6)
	binary.LittleEndian.PutUint32(b, deviceID)
	return b
}

func converDeviceIDInt(deviceID []byte) uint32 {
	return binary.LittleEndian.Uint32(deviceID)
}

type Device struct {
	lanService               Lan
	ip                       net.IP
	id                       uint64
	port                     string
	keepLastKnownOnlineState bool
	deviceType               byte
	updating                 bool
	deferUpdate              bool
	halfTempStep             bool
	support                  bool
	online                   bool
	active                   bool
	protocolVersion          uint
	token                    []byte
	key                      []byte
	name                     string
	modelNumber              string
	serialNumber             string
}

func NewDevice(deviceIP string, deviceID uint64, devicePort string) Device {
	return Device{
		lanService:               newLan(deviceIP, deviceID, devicePort),
		ip:                       net.IP(deviceIP),
		id:                       deviceID,
		port:                     devicePort,
		keepLastKnownOnlineState: false,
		deviceType:               0xAC,
		updating:                 false,
		deferUpdate:              false,
		halfTempStep:             false,
		support:                  false,
		online:                   true,
		active:                   true,
		protocolVersion:          2,
	}
}

func (device *Device) AuthenticateV3(key string, token string) bool {
	var err error
	device.protocolVersion = 3
	if device.token, err = hex.DecodeString(token); err != nil {
		return false
	}
	if device.key, err = hex.DecodeString(key); err != nil {
		return false
	}
	return device.authenticate()
}

func (device *Device) authenticate() bool {
	return device.lanService.Authenticate(device.token, device.key)
}

func (device *Device) SetDeviceDetail(deviceDetail DeviceDetail) {
	device.id = deviceDetail.ID
	device.name = deviceDetail.Name
	device.modelNumber = deviceDetail.ModelNumber
	device.serialNumber = deviceDetail.SerialNumber
	device.deviceType = deviceDetail.DeviceType
	device.active = deviceDetail.ActiveStatus
	device.online = deviceDetail.OnlineStatus
}

type FanSpeed byte

const (
	Auto   FanSpeed = 102
	Full   FanSpeed = 100
	High   FanSpeed = 80
	Medium FanSpeed = 60
	Low    FanSpeed = 40
	Silent FanSpeed = 20
)

type OperationalMode byte

const (
	auto OperationalMode = iota + 1
	cool
	dry
	heat
	fanOnly
)

type SwingMode byte

const (
	Off        = 0x0
	Vertical   = 0xC
	Horizontal = 0x3
	Both       = 0xF
)

type AirConditioningDevice struct {
	Device
	fanSpeed           FanSpeed
	operationalMode    OperationalMode
	swingMode          SwingMode
	promtTone          bool
	powerState         bool
	targetTemperature  float64
	ecoMode            bool
	turboMode          bool
	fahrenheitUnit     bool
	onTimer            interface{}
	offTimer           interface{}
	onlive             bool
	active             bool
	indoorTemperature  float64
	outdoorTemperature float64
}

func NewAirConditioningDevice(d Device) AirConditioningDevice {
	return AirConditioningDevice{
		Device:          d,
		fanSpeed:        Auto,
		operationalMode: auto,
		swingMode:       Off,

		promtTone:          false,
		powerState:         false,
		targetTemperature:  17.0,
		ecoMode:            false,
		turboMode:          false,
		fahrenheitUnit:     false,
		onTimer:            nil,
		offTimer:           nil,
		onlive:             true,
		active:             true,
		indoorTemperature:  0.0,
		outdoorTemperature: 0.0,
	}
}

func (d *AirConditioningDevice) Refresh() {
	cmd := NewBaseCommand(d.deviceType)
	d.sendCmd(cmd)
}

func (d *AirConditioningDevice) sendCmd(cmd BaseCommand) {
	pktBuilder := NewPacketBuilder(d.id)
	pktBuilder.SetCommand(cmd)
	data := pktBuilder.Finalize()
	//TODO LOG debug
	// sendTime := time.Now()
	var responses [][]byte
	if d.protocolVersion == 3 {
		responses = d.lanService.ApplianceTransparentSend8370(data, MSGTYPE_ENCRYPTED_REQUEST)
	} else {
		responses = d.lanService.ApplianceTransparentSend(data)
	}
	// requestTime = time.Now().Sub(sendTime)
	// TODO logger
	if len(responses) == 0 {
		d.active = false
		d.support = false
	}
	for _, response := range responses {
		d.processResponse(response)
	}
}

func (d *AirConditioningDevice) processResponse(data []byte) {
	if len(data) > 0 {
		d.online = true
		d.active = true
		if string(data) == "ERROR" {
			d.support = false
			log.Printf("Got ERROR from %s %d", d.ip, d.id)
			return
		}
		response := NewApplianceResponse(data)
		d.deferUpdate = false
		d.support = true

		if data[0xA] == 0xC0 {
			d.Update(response)
		}
		if data[0xA] == 0xA1 || data[0xA] == 0xA0 {
			// TODO ??
			return
		}
	} else if !d.keepLastKnownOnlineState {
		d.online = false
	}
}

func (d *AirConditioningDevice) Apply() {
	d.updating = true

	cmd := NewSetCommand(d.deviceType)
	cmd.SetPromtTone(d.promtTone)
	cmd.SetPowerState(d.powerState)
	cmd.SetTargetTemperature(d.targetTemperature)
	cmd.SetOperationalMode(byte(d.operationalMode))
	cmd.SetFanSpeed(byte(d.fanSpeed))
	cmd.SetSwingMode(byte(d.swingMode))
	cmd.SetEcoMode(d.ecoMode)
	cmd.SetTurboMode(d.turboMode)
	cmd.SetFahrenheit(d.fahrenheitUnit)
	d.sendCmd(cmd.BaseCommand)

	d.updating = false
	d.deferUpdate = false
}

func (d *AirConditioningDevice) Update(res ApplianceResponse) {
	d.powerState = res.PowerState()
	d.targetTemperature = res.TargetTemperature()
	d.operationalMode = OperationalMode(res.OperationalMode())
	d.fanSpeed = FanSpeed(res.FanSpeed())
	d.swingMode = SwingMode(res.SwingMode())
	d.ecoMode = res.EcoMode()
	d.turboMode = res.TurboMode()
	if res.IndoorTemperature() != 0xFF {
		d.indoorTemperature = res.IndoorTemperature()
	}
	if res.OutdoorTemperature() != 0xFF {
		d.outdoorTemperature = res.OutdoorTemperature()
	}
	d.onTimer = res.OnTimer()
	d.offTimer = res.OffTimer()
}

func (d *AirConditioningDevice) UpdateSpecial(res ApplianceResponse) {
	if res.IndoorTemperature() != 0xFF {
		d.indoorTemperature = res.IndoorTemperature()
	}
	if res.OutdoorTemperature() != 0xFF {
		d.outdoorTemperature = res.OutdoorTemperature()
	}
}

func (d AirConditioningDevice) PromtTone() bool {
	return d.promtTone
}
func (d *AirConditioningDevice) SetPromtTone(feedback bool) {
	if d.updating {
		d.deferUpdate = true
	}
	d.promtTone = feedback
}

func (d AirConditioningDevice) PowerState() bool {
	return d.powerState
}
func (d *AirConditioningDevice) SetPowerState(state bool) {
	if d.updating {
		d.deferUpdate = true
	}
	d.powerState = state
}

func (d AirConditioningDevice) TargetTemperature() float64 {
	return d.targetTemperature
}
func (d *AirConditioningDevice) SetTargetTemperature(t float64) {
	if d.updating {
		d.deferUpdate = true
	}
	d.targetTemperature = t
}

func (d AirConditioningDevice) OperationalMode() OperationalMode {
	return d.operationalMode
}
func (d *AirConditioningDevice) SetOperationalMode(mode OperationalMode) {
	if d.updating {
		d.deferUpdate = true
	}
	d.operationalMode = mode
}

func (d AirConditioningDevice) FanSpeed() FanSpeed {
	return d.fanSpeed
}
func (d *AirConditioningDevice) SetFanSpeed(speed FanSpeed) {
	if d.updating {
		d.deferUpdate = true
	}
	d.fanSpeed = speed
}

func (d AirConditioningDevice) SwingMode() SwingMode {
	return d.swingMode
}
func (d *AirConditioningDevice) SetSwingMode(mode SwingMode) {
	if d.updating {
		d.deferUpdate = true
	}
	d.swingMode = mode
}

func (d AirConditioningDevice) TurboMode() bool {
	return d.turboMode
}
func (d *AirConditioningDevice) SetTurboMode(mode bool) {
	if d.updating {
		d.deferUpdate = true
	}
	d.turboMode = mode
}

func (d AirConditioningDevice) IndoorTemperature() float64 {
	return d.indoorTemperature
}

func (d AirConditioningDevice) Outdoortemperature() float64 {
	return d.outdoorTemperature
}

func (d AirConditioningDevice) OnTimer() interface{} {
	return d.onTimer
}

func (d AirConditioningDevice) OffTimer() interface{} {
	return d.offTimer
}
