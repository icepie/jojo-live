package mideaac

type DeviceDetail struct {
	ID           uint64
	Name         string
	ModelNumber  string
	SerialNumber string
	DeviceType   byte
	ActiveStatus bool
	OnlineStatus bool
}
