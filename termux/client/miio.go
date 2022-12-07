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
	getProps, err := MiLight.Send("get_prop", []interface{}{"name",
		"lan_ctrl",
		"save_state",
		"delayoff",
		"music_on",
		"power",
		"bright",
		"color_mode",
		"rgb",
		"hue",
		"sat",
		"ct",
		"flowing",
		"flow_params",
		"active_mode",
		"nl_br",
		"bg_power",
		"bg_bright",
		"bg_lmode",
		"bg_rgb",
		"bg_hue",
		"bg_sat",
		"bg_ct",
		"bg_flowing",
		"bg_flow_params",
	})
	if err != nil {
		panic(err)
	}

	log.Println(string(getProps))
}
