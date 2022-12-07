package client

import (
	"encoding/json"
	"log"

	"github.com/icepie/miio.go"
)

type MiLightStatus struct {
	Id     int           `json:"id"`
	Result []interface{} `json:"result"`
}

func GetMiLightStatus() (status MiLightStatus, err error) {

	MiLight := miio.New("192.168.2.214").SetToken("1b89c51e6d9d95a36300238b77170a98").SetDid("")

	// for get siid and piid, see https://home.miot-spec.com/
	// get properties
	getProps, err := MiLight.Send("get_prop", []interface{}{"power",
		"bright",
		"color_mode",
		"ct",
	})
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(getProps, &status)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func SetMiLightPower(power bool) (err error) {

	MiLight := miio.New("192.168.2.214").SetToken("1b89c51e6d9d95a36300238b77170a98").SetDid("")

	if power {
		// set properties
		_, err = MiLight.Send("set_power", []interface{}{"on", "smooth", 500})
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		// set properties
		_, err = MiLight.Send("set_power", []interface{}{"off", "smooth", 500})
		if err != nil {
			log.Println(err)
			return
		}
	}
	return
}
