package client

import (
	"github.com/icepie/miio.go"
)

var (
	MiLight *miio.Client
)

func init() {
	MiLight = miio.New("192.168.2.214").SetToken("1b89c51e6d9d95a36300238b77170a98").SetDid("")
}
