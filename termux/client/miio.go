package client

import (
	"github.com/icepie/miio.go"
)

var (
	MiLight *miio.Client
)

func init() {
	MiLight = miio.New("192.168.1.12").SetToken("").SetDid("")
}
