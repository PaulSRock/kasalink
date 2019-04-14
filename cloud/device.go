package cloud

// TPLinkDevice contains the device properties given to us by the TPLink Cloud
type TPLinkDevice struct {
	AppServerURL          string `json:"appServerUrl"`
	IsSameRegion          bool   `json:"isSameRegion"`
	DeviceMAC             string `json:"deviceMac"`
	Status                int    `json:"status"`
	HardWareID            string `json:"hwId"`
	DeviceID              string `json:"deviceId"`
	OEMID                 string `json:"oemId"`
	FirmwareVersion       string `json:"fwVer"`
	DeviceType            string `json:"deviceType"`
	Alias                 string `json:"alias"`
	FirmwareID            string `json:"fwId"`
	DeviceName            string `json:"deviceName"`
	DeviceHardwareVersion string `json:"deviceHwVer"`
	Role                  int    `json:"role"`
	DeviceModel           string `json:"deviceModel"`
}
