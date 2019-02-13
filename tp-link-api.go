package kasalink

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// KasaGetSystemInfo is the JSON used to query device information from a Kasa API device
const KasaGetSystemInfo = `{"system":{"get_sysinfo":{}}}`

// KasaReboot is the JSON used to issue a reboot command to a Kasa API device
const KasaReboot = `{"system":{"reboot":{"delay":1}}}`

// KasaFactorReset is the JSON used to issue a factory reset command to a Kasa API device
const KasaFactorReset = `{"system":{"reset":{"delay":1}}}`

// KasaTurnDeviceOn is the JSON used to issue a power on command to every socket on a Kasa enabled device
// It does not turn on the Kasa Device itself.
const KasaTurnDeviceOn = `{"system":{"set_relay_state":{"state":1}}}`

// KasaTurnPlugOn is the JSON used to issue a power on command to individual sockets on a Kasa enabled device
// It does not turn on the Kasa Device itself.
func KasaTurnPlugOn(deviceID string, children ...int) string {
	var (
		sb          strings.Builder
		err         error
		finalString string
	)
	if _, err = sb.WriteString(`{"context":{"child_ids":[`); err != nil {
		log.Fatal(err)
	}
	for _, child := range children {
		if _, err = sb.WriteString(fmt.Sprintf(`"%s%02d",`, deviceID, child)); err != nil {
			log.Fatal(err)
		}
	}
	if _, err = sb.WriteString(`]},"system":{"set_relay_state":{"state":1}}}`); err != nil {
		log.Fatal(err)
	}
	finalString = sb.String()
	return trimJSONArray(finalString)
}

// KasaTurnDeviceOff is the JSON used to issue a power off command to a Kasa Enabled switch or socket
// It does not turn off the Kasa Device itself.
const KasaTurnDeviceOff = `{"system":{"set_relay_state":{"state":0}}}`

// KasaTurnPlugOff is the JSON used to issue a power on command to a Kasa Enabled switch or socket
// It does not turn on the Kasa Device itself.
func KasaTurnPlugOff(deviceID string, children ...int) string {
	var (
		sb          strings.Builder
		err         error
		finalString string
	)

	if _, err = sb.WriteString(`{"context":{"child_ids":[`); err != nil {
		log.Fatal(err)
	}
	for _, child := range children {
		if _, err = sb.WriteString(fmt.Sprintf(`"%s%02d",`, deviceID, child)); err != nil {
			log.Fatal(err)
		}
	}
	if _, err = sb.WriteString(`]},"system":{"set_relay_state":{"state":0}}}`); err != nil {
		log.Fatal(err)
	}
	finalString = sb.String()
	return trimJSONArray(finalString)
}

// KasaDisableLED is the JSON used to turn off the LED indicator for a Kasa enabled switch or socket
const KasaDisableLED = `{"system":{"set_led_off":{"off":1}}}`

// KasaEnableLED is the JSON used to turn on the LED indicator for a Kasa enabled switch or socket
const KasaEnableLED = `{"system":{"set_led_off":{"off":0}}}`

//KasaSetDeviceAliasString takes a string to assign as the device alias and returns the JSON required to do so
func KasaSetDeviceAliasString(alias string) string {
	return fmt.Sprintf("{\"system\":{\"set_dev_alias\":{\"alias\":\"%s\"}}}", alias)
}

// KasaSetDeviceMACString takes a string and returns the JSON required to do so
// If the input string fails to validate as a MAC address this function returns an empty string and an error
func KasaSetDeviceMACString(newMAC string) (formatedJSON string, err error) {
	if _, err = net.ParseMAC(newMAC); err != nil {
		return "", err
	}
	formatedJSON = fmt.Sprintf("{\"system\":{\"set_mac_addr\":{\"mac\":\"%s\"}}}", newMAC)
	return formatedJSON, nil
}

//KasaSetDeviceID takes a string and returns the JSON required to set the new Device ID
func KasaSetDeviceID(newDeviceID string) string {
	return fmt.Sprintf("{\"system\":{\"set_device_id\":{\"deviceId\":\"%s\"}}}", newDeviceID)
}

// KasaSetHardwareID takes a string and returns the JSON required to set the new Hardware ID
func KasaSetHardwareID(newHardwareID string) string {
	return fmt.Sprintf("{\"system\":{\"set_hw_id\":{\"hwId\":\"%s\"}}}", newHardwareID)
}

//KasaSetLongLat returns the JSON required to set the location of a device
func KasaSetLongLat(long, lat float64) string {
	return fmt.Sprintf("{\"system\":{\"set_dev_location\":{\"longitude\":%f,\"latitude\":%f}}}", long, lat)
}

//KasaBootloaderCheck is the JSON to perform a uBoot bootloader check
const KasaBootloaderCheck = `{"system":{"test_check_uboot":}}`

//KasaGetDeviceIcon is the JSON to get the device icon
const KasaGetDeviceIcon = `{"system":{"get_dev_icon":}}`

//KasaSetDeviceIcon returns the JSON to set the devce icon
func KasaSetDeviceIcon(s1, s2 string) string {
	return fmt.Sprintf("{\"system\":{\"set_dev_icon\":{\"icon\":\"%s\",\"hash\":\"%s\"}}}", s1, s2)
}

//Set Test Mode (command only accepted coming from IP 192.168.1.100)
//const KasaSetTestMode  = "{"system":{"set_test_mode":{"enable":1}}}"

//KasaDownloadFirmaware returns the JSON to download firmware from a given URL
func KasaDownloadFirmaware(url string) string {
	return fmt.Sprintf("{\"system\":{\"download_firmware\":{\"url\":\"%s\"}}}", url)
}

//KasaGetDownloadState is the JSON to get current download state
const KasaGetDownloadState = `{"system":{"get_download_state":{}}}`

//KasaFlashDownloadedFirmware is the JSON to initiate a flash for the currently downloaded firmware
const KasaFlashDownloadedFirmware = `{"system":{"flash_firmware":{}}}`

//KasaCheckNewConfig is the JSON to check the current configuration of the device
const KasaCheckNewConfig = `{"system":{"check_new_config":}}`

//WLAN Commands

//KasaScanForAccessPoints is the JSON to tell the device to scan for list of available wireless access points
const KasaScanForAccessPoints = `{"netif":{"get_scaninfo":{"refresh":1}}}`

//KasaConnectToAccessPoint Connect to AP with given SSID and Password
func KasaConnectToAccessPoint(ssid, passwd string) string {
	return fmt.Sprintf("{\"netif\":{\"set_stainfo\":{\"ssid\":\"%s\",\"password\":\"%s\",\"key_type\":3}}}", ssid, passwd)
}

//Cloud Commands

//KasaGetCloudInfo is the JSON to retrieve the current cloud configuration (Server, Username, Connection Status)
const KasaGetCloudInfo = `{"cnCloud":{"get_info":}}`

//KasaGetFirmwareFromCloud is the JSON to retrieve a list of firmware from the cloud server
const KasaGetFirmwareFromCloud = `{"cnCloud":{"get_intl_fw_list":{}}}`

//KasaSetServerURL returns the JSON required to set a new server URL
func KasaSetServerURL(newServer string) string {
	return fmt.Sprintf("{\"cnCloud\":{\"set_server_url\":{\"server\":\"%s\"}}}", newServer)
}

//KasaSetDefaultServerURL is the JSON to set the default server URL (devs.tplinkcloud.com)
const KasaSetDefaultServerURL = `{"cnCloud":{"set_server_url":{"server":"devs.tplinkcloud.com"}}}`

//KasaConnectWithUserPass returns the JSON required to connect to the TP-Link Cloud service with a username & password
func KasaConnectWithUserPass(user, pass string) string {
	return fmt.Sprintf("{\"cnCloud\":{\"bind\":{\"username\":\"%s\", \"password\":\"%s\"}}}", user, pass)
}

//KasaUnregisterFromCloud is the JSON to unregister the device from a TP-Link Cloud Account
const KasaUnregisterFromCloud = `{"cnCloud":{"unbind":}}`

//Time Commands

//KasaGetDeviceTime is the JSON to retrieve the current device time
const KasaGetDeviceTime = `{"time":{"get_time":}}`

//KasaGetDeviceTimezone is the JSON to get the current device timezone
const KasaGetDeviceTimezone = `{"time":{"get_timezone":}}`

//KasaSetDeviceTimeZone returns the JSON to set the time, date and time zone
func KasaSetDeviceTimeZone(year, month, day, hour, min, sec, tzindex int) string {
	return fmt.Sprintf("{\"time\":{\"set_timezone\":{\"year\":%d,\"month\":%d,\"mday\":%d,\"hour\":%d,\"min\":%d,\"sec\":%d,\"index\":%d}}}",
		year, month, day, hour, min, sec, tzindex)
}

//EMeter Energy Usage Statistics Commands
//(for TP-Link HS110)

//KasaGetRealtimeCurrentAndVoltage is the JSON to get realtime current and voltage readings
const KasaGetRealtimeCurrentAndVoltage = `{"emeter":{"get_realtime":{}}}`

//KasaGetVGainAndIGain is the JSON to get EMeter VGain and IGain settings
const KasaGetVGainAndIGain = `{"emeter":{"get_vgain_igain":{}}}`

//KasaSetVGainAndIGain returns the JSON to set EMeter VGain and Igain values
func KasaSetVGainAndIGain(newVGain, newIGain int) string {
	return fmt.Sprintf("{\"emeter\":{\"set_vgain_igain\":{\"vgain\":%d,\"igain\":%d}}}", newVGain, newIGain)
}

//KasaStartEMeterCalibration returns the JSON to start EMeter calibration
func KasaStartEMeterCalibration(vTarget, iTarget int) string {
	return fmt.Sprintf("{\"emeter\":{\"start_calibration\":{\"vtarget\":%d,\"itarget\":%d}}}", vTarget, iTarget)
}

//KasaGetDailyStatsForMonthYear returns the JSON to get daily statistic for a given month
func KasaGetDailyStatsForMonthYear(month, year int) (formattedJSON string, err error) {
	if month > 11 || month < 0 {
		return "", fmt.Errorf("%d is an invalid value for month [0-11]", month)
	}
	if year < 0 || year > time.Now().Year() {
		return "", fmt.Errorf("%d is an invalid value for year", year)
	}
	if year >= time.Now().Year() && month > int(time.Now().Month()) {
		return "", fmt.Errorf("%d/%d appear to be a month/year in the future", month, year)
	}
	formattedJSON = fmt.Sprintf("{\"emeter\":{\"get_daystat\":{\"month\":%d,\"year\":%d}}}", month, year)
	return formattedJSON, nil
}

//KasaGetMonthlyStatsForYear returns the JSON required to get monthly statistic for given year
func KasaGetMonthlyStatsForYear(year int) (formattedJSON string, err error) {
	if year > time.Now().Year() {
		return "", fmt.Errorf("%d appears to be in the future", year)
	}
	formattedJSON = fmt.Sprintf("{\"emeter\":{\"\"get_monthstat\":{\"year\":%d}}}", year)
	return formattedJSON, nil
}

//KasaEraseEMeterStats is the JSON to erase all EMeter statistics
const KasaEraseEMeterStats = `{"emeter":{"erase_emeter_stat":}}`

//Schedule Commands
//(action to perform regularly on given weekdays)

//KasaGetNexedScheduledAction is the JSON to get the next scheduled action
const KasaGetNexedScheduledAction = `{"schedule":{"get_next_action":}}`

//KasaGetScheduleRulesList is the JSON to get the schedule rules list
const KasaGetScheduleRulesList = `{"schedule":{"get_rules":}}`

//KasaAddScheduleRule returns the JSON required to add a new schedule rule
func KasaAddScheduleRule() string {
	return fmt.Sprintf("{\"schedule\":{\"add_rule\":{\"stime_opt\":0,\"wday\":[1,0,0,1,1,0,0],\"smin\":1014,\"enable\":1,\"repeat\":1,\"etime_opt\":-1,\"name\":\"lights on\",\"eact\":-1,\"month\":0,\"sact\":1,\"year\":0,\"longitude\":0,\"day\":0,\"force\":0,\"latitude\":0,\"emin\":0},\"set_overall_enable\":{\"enable\":1}}}")
}

//KasaEditScheduleRule returns the JSON required to edit a schedule rule with the given ID
func KasaEditScheduleRule() string {
	return fmt.Sprintf("{\"schedule\":{\"edit_rule\":{\"stime_opt\":0,\"wday\":[1,0,0,1,1,0,0],\"smin\":1014,\"enable\":1,\"repeat\":1,\"etime_opt\":-1,\"id\":\"4B44932DFC09780B554A740BC1798CBC\",\"name\":\"lights on\",\"eact\":-1,\"month\":0,\"sact\":1,\"year\":0,\"longitude\":0,\"day\":0,\"force\":0,\"latitude\":0,\"emin\":0}}}")
}

//KasaDeleteScheduleRule returns the JSON to delete a schedule rule with the given ID
func KasaDeleteScheduleRule(id string) string {
	return fmt.Sprintf("{\"schedule\":{\"delete_rule\":{\"id\":\"%s\"}}}", id)
}

//KasaDeleteAllScheduleRules is the JSON to delete all schedule rules and erase statistics
const KasaDeleteAllScheduleRules = `{"schedule":{"delete_all_rules":,"erase_runtime_stat":}}`

//Countdown Rule Commands
//(action to perform after number of seconds)

//KasaGetCountdownRule is the JSON toge the existing countdown rule
const KasaGetCountdownRule = `{"count_down":{"get_rules":}}`

//KasaAddNewCountdownRule is the JSON to add a new countdown rule
func KasaAddNewCountdownRule(enable, delay, act int, name string) string {
	return fmt.Sprintf("{\"count_down\":{\"add_rule\":{\"enable\":%d,\"delay\":%d,\"act\":%d,\"name\":\"%s\"}}}",
		enable, delay, act, name)
}

//KasaEditCountdownRule returns the JSON to edit a countdown rule with the given ID
func KasaEditCountdownRule(enable, delay, act int, name, id string) string {
	return fmt.Sprintf("{\"count_down\":{\"edit_rule\":{\"enable\":%d,\"id\":\"%s\",\"delay\":%d,\"act\":%d,\"name\":\"%s\"}}}",
		enable, id, delay, act, name)
}

//KasaDeleteCountdownRule returns the JSON to delete a countdown rule with the given ID
func KasaDeleteCountdownRule(id string) string {
	return fmt.Sprintf("{\"count_down\":{\"delete_rule\":{\"id\":\"%s\"}}}", id)
}

//KasaDeleteAllCountdownRules is the JSON to delete all countdown rules
const KasaDeleteAllCountdownRules = `{"count_down":{"delete_all_rules":}}`

//Anti-Theft Rule Commands (aka Away Mode)
//(period of time during which device will be randomly turned on and off to deter thieves)

//KasaGetAntiTheftRules is the JSON to retrieve the existing anti-theft rule set
const KasaGetAntiTheftRules = `{"anti_theft":{"get_rules":}}`

//KasaAddAntiTheftRule returns the JSON reuqired to add a new anti-theft rule
func KasaAddAntiTheftRule() string {
	return fmt.Sprintf("{\"anti_theft\":{\"add_rule\":{\"stime_opt\":0,\"wday\":[0,0,0,1,0,1,0],\"smin\":987,\"enable\":1,\"frequency\":5,\"repeat\":1,\"etime_opt\":0,\"duration\":2,\"name\":\"test\",\"lastfor\":1,\"month\":0,\"year\":0,\"longitude\":0,\"day\":0,\"latitude\":0,\"force\":0,\"emin\":1047},\"set_overall_enable\":1}}")
}

//KasaEditAntiTheftRule returns the JSON required to edit an anti-theft rule
func KasaEditAntiTheftRule(id string) string {
	return fmt.Sprintf("{\"anti_theft\":{\"edit_rule\":{\"stime_opt\":0,\"wday\":[0,0,0,1,0,1,0],\"smin\":987,\"enable\":1,\"frequency\":5,\"repeat\":1,\"etime_opt\":0,\"id\":\"%s\",\"duration\":2,\"name\":\"test\",\"lastfor\":1,\"month\":0,\"year\":0,\"longitude\":0,\"day\":0,\"latitude\":0,\"force\":0,\"emin\":1047},\"set_overall_enable\":1}}", id)
}

//KasaDeleteAntiTheftRule returns the JSON required to delete an anti-theft rule with given ID
func KasaDeleteAntiTheftRule(id string) string {
	return fmt.Sprintf("{\"anti_theft\":{\"delete_rule\":{\"id\":\"%s\"}}}", id)
}

//KasaDeleteAllAntiTheftRules is the JSON to delete all the anti-theft rules
const KasaDeleteAllAntiTheftRules = `{"anti_theft":{"delete_all_rules":}}`

func trimJSONArray(s string) string {
	return strings.Replace(s, `,]`, `]`, 1)
}
