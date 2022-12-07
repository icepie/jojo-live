package client

import (
	"log"

	"github.com/icepie/miio.go"
)

var (
	MiLight *miio.Client
)

func init() {
	MiLight = miio.New("192.168.2.214").SetToken("1b89c51e6d9d95a36300238b77170a98").SetDid("")
}

func GetMiLightStatus() {

	// for get siid and piid, see https://home.miot-spec.com/
	// get properties
	getProps, err := MiLight.GetProps(miio.PropParam{
		Siid: 2,
		Piid: 1,
	})
	if err != nil {
		panic(err)
	}

	log.Println(string(getProps))
}
