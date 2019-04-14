package kasalink

//{"system":
// {"get_sysinfo":
//   {"sw_ver":"1.0.6 Build 180627 Rel.081000",
//    "hw_ver":"1.0",
//    "model":"HS300(US)",
//    "deviceId":"8006E92180ADBEA7B3E4820027152BE21ACC7D77",
//    "oemId":"5C9E6254BEBAED63B2B6102966D24C17",
//    "hwId":"34C41AA028022D0CCEA5E678E8547C54",
//    "rssi":-35,
//    "longitude_i":-775702,
//    "latitude_i":391156,
//    "alias":"TP-LINK_Power Strip_14A9",
//    "mic_type":"IOT.SMARTPLUGSWITCH",
//    "feature":"TIM:ENE",
//    "mac":"B0:BE:76:80:14:A9",
//    "updating":0,
//    "led_off":0,
//    "children":[
//      {"id":"8006E92180ADBEA7B3E4820027152BE21ACC7D7700",
//       "state":1,
//       "alias":"Top Tank Light",
//       "on_time":394931,
//       "next_action":{"type":-1}
//      },
//      {"id":"8006E92180ADBEA7B3E4820027152BE21ACC7D7701",
//       "state":1,
//       "alias":"Top Tank Heater",
//       "on_time":1200805,
//       "next_action":{"type":-1}
//      },
//      {"id":"8006E92180ADBEA7B3E4820027152BE21ACC7D7702",
//       "state":1,
//       "alias":"Top Tank Filter",
//       "on_time":50621,
//       "next_action":{"type":-1}
//      },
//      {"id":"8006E92180ADBEA7B3E4820027152BE21ACC7D7703",
//       "state":1,
//       "alias":"Top Tank Powerhead",
//       "on_time":385855,
//       "next_action":{"type":-1}
//      },
//      {"id":"8006E92180ADBEA7B3E4820027152BE21ACC7D7704",
//       "state":1,
//       "alias":"Air Pump",
//       "on_time":1200806,
//       "next_action":{"type":-1}
//      },
//      {"id":"8006E92180ADBEA7B3E4820027152BE21ACC7D7705",
//       "state":1,
//       "alias":"Plug 6",
//       "on_time":1289431,
//       "next_action":{"type":-1}
//      }],
//      "child_num":6,
//      "err_code":0
//    }
//  }
//}

type action struct {
	Type int `json:"type"`
}

type childState struct {
	ID         string `json:"id"`
	State      int    `json:"state"`
	Alias      string `json:"alias"`
	OnTime     int    `json:"on_time"`
	NextAction action
}

// SystemInfo is the system information about a TP-Link/Kasa device
type SystemInfo struct {
	SoftwareVersion string       `json:"sw_ver"`
	HardwareVersion string       `json:"hw_ver"`
	Model           string       `json:"model"`
	DeviceID        string       `json:"deviceId"`
	OEMID           string       `json:"oemId"`
	HardwareID      string       `json:"hwId"`
	RSSI            int          `json:"rssi"`
	Longitude       int          `json:"longitude_i"`
	Latitude        int          `json:"latitude_i"`
	Alias           string       `json:"alias"`
	MICType         string       `json:"mic_type"`
	Feature         string       `json:"feature"`
	MAC             string       `json:"mac"`
	Updating        int          `json:"updating"`
	LEDOff          int          `json:"led_off"`
	Children        []childState `json:"children"`
	ChildNum        int          `json:"child_num"`
	ErrCode         int          `json:"err_code"`
}

type systemResponse struct {
	GetSysInfo *SystemInfo       `json:"get_sysinfo,omitempty"`
	SetLED     *thingWithErrCode `json:"set_led_off,omitempty"`
}

// KasaResponse is a wrapper object for JSON responses from KasaPlugs
type KasaResponse struct {
	System      *systemResponse `json:"system,omitempty"`
	EnergyMeter *energyMeter    `json:"emeter,omitempty"`
}

type energyMeter struct {
	Realtime *realtimeEnergyMeter `json:"get_realtime"`
}

type realtimeEnergyMeter struct {
	Voltage    int `json:"voltage_mv"`
	Current    int `json:"current_ma"`
	Power      int `json:"power_mw"`
	TotalWatts int `json:"total_wh"`
	thingWithErrCode
}

type thingWithErrCode struct {
	ErrorCode    int    `json:"err_code,omitempty"`
	ErrorMessage string `json:"err_msg,omitempty"`
}
