package kasalink

import (
	"fmt"
	"net"
)

type unsafe struct {
}

// FactoryReset is the JSON used to issue a factory reset command to a Kasa API device
func (unsafe) FactoryReset() string {
	return `{"system":{"reset":{"delay":1}}}`
}

// SetDeviceMACString takes a string and returns the JSON required to do so
// If the input string fails to validate as a MAC address this function returns an empty string and an error
func (unsafe) SetDeviceMACString(newMAC string) (formatedJSON string, err error) {
	if _, err = net.ParseMAC(newMAC); err != nil {
		return "", err
	}
	formatedJSON = fmt.Sprintf("{\"system\":{\"set_mac_addr\":{\"mac\":\"%s\"}}}", newMAC)
	return formatedJSON, nil
}

// SetDeviceID takes a string and returns the JSON required to set the new Device ID
func (unsafe) SetDeviceID(newDeviceID string) string {
	return fmt.Sprintf("{\"system\":{\"set_device_id\":{\"deviceId\":\"%s\"}}}", newDeviceID)
}

// SetHardwareID takes a string and returns the JSON required to set the new Hardware ID
func (unsafe) SetHardwareID(newHardwareID string) string {
	return fmt.Sprintf("{\"system\":{\"set_hw_id\":{\"hwId\":\"%s\"}}}", newHardwareID)
}

//KasaBootloaderCheck is the JSON to perform a uBoot bootloader check
func (unsafe) BootloaderCheck() string {
	return `{"system":{"test_check_uboot":}}`
}

// Set Test Mode (command only accepted coming from IP 192.168.1.100)
func (unsafe) SetTestMode() string {
	return `{"system":{"set_test_mode":{"enable":1}}}`
}

// DownloadFirmaware returns the JSON to download firmware from a given URL
func (unsafe) DownloadFirmaware(url string) string {
	return fmt.Sprintf("{\"system\":{\"download_firmware\":{\"url\":\"%s\"}}}", url)
}

// GetDownloadState is the JSON to get current download state
func (unsafe) GetDownloadState() string {
	return `{"system":{"get_download_state":{}}}`
}

// FlashDownloadedFirmware is the JSON to initiate a flash for the currently downloaded firmware
func (unsafe) FlashDownloadedFirmware() string {
	return `{"system":{"flash_firmware":{}}}`
}

// CheckNewConfig is the JSON to check the current configuration of the device
func (unsafe) CheckNewConfig() string {
	return `{"system":{"check_new_config":}}`
}
