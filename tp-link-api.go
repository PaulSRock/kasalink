package kasalink

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	getSysInfo                         = `{"system":{"get_sysinfo":{}}}`
	reboot                             = `{"system":{"reboot":{"delay":1}}}`
	turnOn                             = `{"system":{"set_relay_state":{"state":1}}}`
	turnOff                            = `{"system":{"set_relay_state":{"state":0}}}`
	turnOffLED                         = `{"system":{"set_led_off":{"off":1}}}`
	turnOnLED                          = `{"system":{"set_led_off":{"off":0}}}`
	getDeviceIcon                      = `{"system":{"get_dev_icon":}}`
	getCloudInfo                       = `{"cnCloud":{"get_info":}}`
	getFirmwareList                    = `{"cnCloud":{"get_intl_fw_list":{}}}`
	setDefaultCloudURL                 = `{"cnCloud":{"set_server_url":{"server":"devs.tplinkcloud.com"}}}`
	unbindDeviceFromCloud              = `{"cnCloud":{"unbind":}}`
	getDeviceTime                      = `{"time":{"get_time":}}`
	getDeviceTimeZone                  = `{"time":{"get_timezone":}}`
	getCurrentAndVoltage               = `{"emeter":{"get_realtime":{}}}`
	getVandIGain                       = `{"emeter":{"get_vgain_igain":{}}}`
	scanForAccessPoints                = `{"netif":{"get_scaninfo":{"refresh":1}}}`
	setDeviceAliasFormatString         = "{\"system\":{\"set_dev_alias\":{\"alias\":\"%s\"}}}"
	latLongFormatString                = "{\"system\":{\"set_dev_location\":{\"longitude\":%f,\"latitude\":%f}}}"
	setDeviceIconFormatString          = "{\"system\":{\"set_dev_icon\":{\"icon\":\"%s\",\"hash\":\"%s\"}}}"
	eraseEnergyMeterStats              = `{"emeter":{"erase_emeter_stat":}}`
	connecToAccessPointFormatString    = "{\"netif\":{\"set_stainfo\":{\"ssid\":\"%s\",\"password\":\"%s\",\"key_type\":3}}}"
	setCloudURLFormatString            = "{\"cnCloud\":{\"set_server_url\":{\"server\":\"%s\"}}}"
	cloudConnectFormatString           = "{\"cnCloud\":{\"bind\":{\"username\":\"%s\", \"password\":\"%s\"}}}"
	setDeviceTimeFormatString          = "{\"time\":{\"set_timezone\":{\"year\":%d,\"month\":%d,\"mday\":%d,\"hour\":%d,\"min\":%d,\"sec\":%d,\"index\":%d}}}"
	setVandIGainFormatString           = "{\"emeter\":{\"set_vgain_igain\":{\"vgain\":%d,\"igain\":%d}}}"
	startEMeterCalibrationFormatString = "{\"emeter\":{\"start_calibration\":{\"vtarget\":%d,\"itarget\":%d}}}"
	getDailyEnergyStatsFormatString    = "{\"emeter\":{\"get_daystat\":{\"month\":%d,\"year\":%d}}}"
	getMonthlyEnergyStatsFormatString  = "{\"emeter\":{\"\"get_monthstat\":{\"year\":%d}}}"
)

// GetSystemInfo is the is the Struct that contains info about the Kasa Device
func (kpp *KasaPowerPlug) GetSystemInfo() (*SystemInfo, error) {
	if kpp.SysInfo != nil {
		return kpp.SysInfo, nil
	}
	b, err := kpp.querySystemInfo()
	if err != nil {
		return nil, err
	}
	si := &KasaResponse{}
	err = json.Unmarshal(b, si)
	return si.System.GetSysInfo, err
}

func (kpp *KasaPowerPlug) querySystemInfo(children ...int) ([]byte, error) {
	if children != nil {
		return kpp.tellChild(getSysInfo, children...)
	}
	return kpp.talkToPlug(getSysInfo)
}

// Reboot is the JSON used to issue a reboot command to a Kasa API device
func (kpp *KasaPowerPlug) Reboot(children ...int) ([]byte, error) {
	if children != nil {
		return kpp.tellChild(reboot, children...)
	}
	return kpp.talkToPlug(reboot)
}

// TurnDeviceOn is the JSON used to issue a power on command to every socket on a Kasa enabled device
// It does not turn on the Kasa Device itself.
func (kpp *KasaPowerPlug) TurnDeviceOn(children ...int) ([]byte, error) {
	if children != nil {
		return kpp.tellChild(turnOn, children...)
	}
	return kpp.talkToPlug(turnOn)
}

// TurnDeviceOff is the JSON used to issue a power off command to a Kasa Enabled switch or socket
// It does not turn off the Kasa Device itself.
func (kpp *KasaPowerPlug) TurnDeviceOff(children ...int) ([]byte, error) {
	if children != nil {
		return kpp.tellChild(turnOff, children...)
	}
	return kpp.talkToPlug(turnOff)
}

// DisableLED is the JSON used to turn off the LED indicator for a Kasa enabled switch or socket
func (kpp *KasaPowerPlug) DisableLED() (wrapper *KasaResponse, err error) {
	var jsonBytes []byte
	jsonBytes, err = kpp.talkToPlug(turnOffLED)
	if err != nil {
		return
	}
	log.Printf("%s", jsonBytes)
	wrapper = &KasaResponse{}
	err = json.Unmarshal(jsonBytes, wrapper)
	return
}

// EnableLED is the JSON used to turn on the LED indicator for a Kasa enabled switch or socket
func (kpp *KasaPowerPlug) EnableLED() (wrapper *KasaResponse, err error) {
	var jsonBytes []byte
	jsonBytes, err = kpp.talkToPlug(turnOnLED)
	if err != nil {
		return
	}
	log.Printf("%s", jsonBytes)
	wrapper = &KasaResponse{}
	err = json.Unmarshal(jsonBytes, wrapper)
	return
}

// SetDeviceAliasString takes a string to assign as the device alias
func (kpp *KasaPowerPlug) SetDeviceAliasString(alias string, children ...int) ([]byte, error) {
	if children != nil {
		return kpp.tellChild(fmt.Sprintf(setDeviceAliasFormatString, alias), children...)
	}
	return kpp.talkToPlug(fmt.Sprintf(setDeviceAliasFormatString, alias))
}

// SetLongLat returns the JSON required to set the location of a device
func (kpp *KasaPowerPlug) SetLongLat(long, lat float64) ([]byte, error) {
	return kpp.talkToPlug(fmt.Sprintf(latLongFormatString, long, lat))
}

// GetDeviceIcon is the JSON to get the device icon
func (kpp *KasaPowerPlug) GetDeviceIcon(children ...int) ([]byte, error) {
	if children != nil {
		return kpp.tellChild(getDeviceIcon)
	}
	return kpp.talkToPlug(getDeviceIcon)
}

// SetDeviceIcon returns the JSON to set the devce icon
func (kpp *KasaPowerPlug) SetDeviceIcon(s1, s2 string, children ...int) ([]byte, error) {
	if children != nil {
		return kpp.tellChild(fmt.Sprintf(setDeviceIconFormatString, s1, s2), children...)
	}
	return kpp.talkToPlug(fmt.Sprintf(setDeviceIconFormatString, s1, s2))
}

//WLAN Commands

// ScanForAccessPoints is the JSON to tell the device to scan for list of available wireless access points
func (kpp *KasaPowerPlug) ScanForAccessPoints(children ...int) ([]byte, error) {
	return kpp.talkToPlug(scanForAccessPoints)
}

// ConnectToAccessPoint Connect to AP with given SSID and Password
func (kpp *KasaPowerPlug) ConnectToAccessPoint(ssid, passwd string) ([]byte, error) {
	return kpp.talkToPlug(fmt.Sprintf(connecToAccessPointFormatString, ssid, passwd))
}

//Cloud Configuration Commands

// GetCloudInfo is the JSON to retrieve the current cloud configuration (Server, Username, Connection Status)
func (kpp *KasaPowerPlug) GetCloudInfo() ([]byte, error) {
	return kpp.talkToPlug(getCloudInfo)
}

// GetFirmwareFromCloud is the JSON to retrieve a list of firmware from the cloud server
func (kpp *KasaPowerPlug) GetFirmwareFromCloud() ([]byte, error) {
	return kpp.talkToPlug(getFirmwareList)
}

// SetServerURL returns the JSON required to set a new server URL
func (kpp *KasaPowerPlug) SetServerURL(newServer string) ([]byte, error) {
	return kpp.talkToPlug(fmt.Sprintf(setCloudURLFormatString, newServer))
}

// SetDefaultServerURL is the JSON to set the default server URL (devs.tplinkcloud.com)
func (kpp *KasaPowerPlug) SetDefaultServerURL() ([]byte, error) {
	return kpp.talkToPlug(setDefaultCloudURL)
}

// ConnectWithUserPass returns the JSON required to connect to the TP-Link Cloud service with a username & password
func (kpp *KasaPowerPlug) ConnectWithUserPass(user, pass string) ([]byte, error) {
	return kpp.talkToPlug(fmt.Sprintf(cloudConnectFormatString, user, pass))
}

// UnregisterFromCloud is the JSON to unregister the device from a TP-Link Cloud Account
func (kpp *KasaPowerPlug) UnregisterFromCloud() ([]byte, error) {
	return kpp.talkToPlug(unbindDeviceFromCloud)
}

//Time Commands

// GetDeviceTime is the JSON to retrieve the current device time
func (kpp *KasaPowerPlug) GetDeviceTime() ([]byte, error) {
	return kpp.talkToPlug(getDeviceTime)
}

// GetDeviceTimezone is the JSON to get the current device timezone
func (kpp *KasaPowerPlug) GetDeviceTimezone() ([]byte, error) {
	return kpp.talkToPlug(getDeviceTimeZone)
}

// SetDeviceTimeZone returns the JSON to set the time, date and time zone
func (kpp *KasaPowerPlug) SetDeviceTimeZone(t *time.Time) ([]byte, error) {
	var _, offset = t.Zone()
	return kpp.talkToPlug(fmt.Sprintf(setDeviceTimeFormatString,
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), offset))
}

//EMeter Energy Usage Statistics Commands
//(for TP-Link HS110)

// GetRealtimeCurrentAndVoltage is the JSON to get realtime current and voltage readings
func (kpp *KasaPowerPlug) GetRealtimeCurrentAndVoltage(children ...int) (response *KasaResponse, err error) {
	var (
		jsonBytes []byte
	)
	if children != nil {
		jsonBytes, err = kpp.tellChild(getCurrentAndVoltage, children...)
	} else {
		jsonBytes, err = kpp.talkToPlug(getCurrentAndVoltage)
	}
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonBytes, &response)
	if err != nil {
		return
	}
	if response.EnergyMeter.Realtime.ErrorCode != 0 {
		return nil, fmt.Errorf(response.EnergyMeter.Realtime.ErrorMessage)
	}
	return
}

// GetVGainAndIGain is the JSON to get EMeter VGain and IGain settings
func (kpp *KasaPowerPlug) GetVGainAndIGain(children ...int) ([]byte, error) {
	if children != nil {
		return kpp.tellChild(getVandIGain, children...)
	}
	return kpp.talkToPlug(getVandIGain)
}

// SetVGainAndIGain returns the JSON to set EMeter VGain and Igain values
func (kpp *KasaPowerPlug) SetVGainAndIGain(newVGain, newIGain int, children ...int) ([]byte, error) {
	if children != nil {
		return kpp.tellChild(fmt.Sprintf(setVandIGainFormatString, newVGain, newIGain), children...)
	}
	return kpp.talkToPlug(fmt.Sprintf(setVandIGainFormatString, newVGain, newIGain))
}

// StartEMeterCalibration returns the JSON to start EMeter calibration
func (kpp *KasaPowerPlug) StartEMeterCalibration(vTarget, iTarget int, children ...int) ([]byte, error) {
	if children != nil {
		return kpp.tellChild(fmt.Sprintf(startEMeterCalibrationFormatString, vTarget, iTarget), children...)
	}
	return kpp.talkToPlug(fmt.Sprintf(startEMeterCalibrationFormatString, vTarget, iTarget))
}

// GetDailyStatsForMonthYear returns the JSON to get daily statistic for a given month
func (kpp *KasaPowerPlug) GetDailyStatsForMonthYear(month, year int, children ...int) ([]byte, error) {
	var jsonCmd string
	if month > 11 || month < 0 {
		return nil, fmt.Errorf("%d is an invalid value for month [0-11]", month)
	}
	if year < 0 || year > time.Now().Year() {
		return nil, fmt.Errorf("%d is an invalid value for year", year)
	}
	if year >= time.Now().Year() && month > int(time.Now().Month()) {
		return nil, fmt.Errorf("%d/%d appear to be a month/year in the future", month, year)
	}
	jsonCmd = fmt.Sprintf(getDailyEnergyStatsFormatString, month, year)
	if children != nil {
		return kpp.tellChild(jsonCmd, children...)
	}
	return kpp.talkToPlug(jsonCmd)
}

// GetMonthlyStatsForYear returns the JSON required to get monthly statistic for given year
func (kpp *KasaPowerPlug) GetMonthlyStatsForYear(year int, children ...int) ([]byte, error) {
	var jsonCmd string
	if year > time.Now().Year() {
		return nil, fmt.Errorf("%d appears to be in the future", year)
	}
	jsonCmd = fmt.Sprintf(getMonthlyEnergyStatsFormatString, year)
	if children != nil {
		return kpp.tellChild(jsonCmd, children...)
	}
	return kpp.talkToPlug(jsonCmd)
}

// EraseEMeterStats is the JSON to erase all EMeter statistics
func (kpp *KasaPowerPlug) EraseEMeterStats(children ...int) ([]byte, error) {
	if children != nil {
		return kpp.tellChild(eraseEnergyMeterStats, children...)
	}
	return kpp.talkToPlug(eraseEnergyMeterStats)
}

//Schedule Commands
//(action to perform regularly on given weekdays)

// GetNexedScheduledAction is the JSON to get the next scheduled action
func (kpp *KasaPowerPlug) GetNexedScheduledAction(children ...int) ([]byte, error) {
	return kpp.talkToPlug(`{"schedule":{"get_next_action":}}`)
}

// GetScheduleRulesList is the JSON to get the schedule rules list
func (kpp *KasaPowerPlug) GetScheduleRulesList(children ...int) ([]byte, error) {
	return kpp.talkToPlug(`{"schedule":{"get_rules":}}`)
}

// AddScheduleRule returns the JSON required to add a new schedule rule
func (kpp *KasaPowerPlug) AddScheduleRule(children ...int) ([]byte, error) {
	return kpp.talkToPlug(fmt.Sprintf("{\"schedule\":{\"add_rule\":{\"stime_opt\":0,\"wday\":[1,0,0,1,1,0,0],\"smin\":1014,\"enable\":1,\"repeat\":1,\"etime_opt\":-1,\"name\":\"lights on\",\"eact\":-1,\"month\":0,\"sact\":1,\"year\":0,\"longitude\":0,\"day\":0,\"force\":0,\"latitude\":0,\"emin\":0},\"set_overall_enable\":{\"enable\":1}}}"))
}

// EditScheduleRule returns the JSON required to edit a schedule rule with the given ID
func (kpp *KasaPowerPlug) EditScheduleRule(children ...int) ([]byte, error) {
	return kpp.talkToPlug(fmt.Sprintf("{\"schedule\":{\"edit_rule\":{\"stime_opt\":0,\"wday\":[1,0,0,1,1,0,0],\"smin\":1014,\"enable\":1,\"repeat\":1,\"etime_opt\":-1,\"id\":\"4B44932DFC09780B554A740BC1798CBC\",\"name\":\"lights on\",\"eact\":-1,\"month\":0,\"sact\":1,\"year\":0,\"longitude\":0,\"day\":0,\"force\":0,\"latitude\":0,\"emin\":0}}}"))
}

// DeleteScheduleRule returns the JSON to delete a schedule rule with the given ID
func (kpp *KasaPowerPlug) DeleteScheduleRule(id string) ([]byte, error) {
	return kpp.talkToPlug(fmt.Sprintf("{\"schedule\":{\"delete_rule\":{\"id\":\"%s\"}}}", id))
}

// DeleteAllScheduleRules is the JSON to delete all schedule rules and erase statistics
func (kpp *KasaPowerPlug) DeleteAllScheduleRules(children ...int) ([]byte, error) {
	return kpp.talkToPlug(`{"schedule":{"delete_all_rules":,"erase_runtime_stat":}}`)
}

//Countdown Rule Commands
//(action to perform after number of seconds)

// GetCountdownRule is the JSON toge the existing countdown rule
func (kpp *KasaPowerPlug) GetCountdownRule(children ...int) ([]byte, error) {
	return kpp.talkToPlug(`{"count_down":{"get_rules":}}`)
}

// AddNewCountdownRule is the JSON to add a new countdown rule
func (kpp *KasaPowerPlug) AddNewCountdownRule(enable, delay, act int, name string) ([]byte, error) {
	return kpp.talkToPlug(fmt.Sprintf("{\"count_down\":{\"add_rule\":{\"enable\":%d,\"delay\":%d,\"act\":%d,\"name\":\"%s\"}}}",
		enable, delay, act, name))
}

// EditCountdownRule returns the JSON to edit a countdown rule with the given ID
func (kpp *KasaPowerPlug) EditCountdownRule(enable, delay, act int, name, id string) ([]byte, error) {
	return kpp.talkToPlug(fmt.Sprintf("{\"count_down\":{\"edit_rule\":{\"enable\":%d,\"id\":\"%s\",\"delay\":%d,\"act\":%d,\"name\":\"%s\"}}}",
		enable, id, delay, act, name))
}

// DeleteCountdownRule returns the JSON to delete a countdown rule with the given ID
func (kpp *KasaPowerPlug) DeleteCountdownRule(id string) ([]byte, error) {
	return kpp.talkToPlug(fmt.Sprintf("{\"count_down\":{\"delete_rule\":{\"id\":\"%s\"}}}", id))
}

// DeleteAllCountdownRules is the JSON to delete all countdown rules
func (kpp *KasaPowerPlug) DeleteAllCountdownRules(children ...int) ([]byte, error) {
	return kpp.talkToPlug(`{"count_down":{"delete_all_rules":}}`)
}

//Anti-Theft Rule Commands (aka Away Mode)
//(period of time during which device will be randomly turned on and off to deter thieves)

// GetAntiTheftRules is the JSON to retrieve the existing anti-theft rule set
func (kpp *KasaPowerPlug) GetAntiTheftRules(children ...int) ([]byte, error) {
	return kpp.talkToPlug(`{"anti_theft":{"get_rules":}}`)
}

// AddAntiTheftRule returns the JSON reuqired to add a new anti-theft rule
func (kpp *KasaPowerPlug) AddAntiTheftRule(children ...int) ([]byte, error) {
	return kpp.talkToPlug(fmt.Sprintf("{\"anti_theft\":{\"add_rule\":{\"stime_opt\":0,\"wday\":[0,0,0,1,0,1,0],\"smin\":987,\"enable\":1,\"frequency\":5,\"repeat\":1,\"etime_opt\":0,\"duration\":2,\"name\":\"test\",\"lastfor\":1,\"month\":0,\"year\":0,\"longitude\":0,\"day\":0,\"latitude\":0,\"force\":0,\"emin\":1047},\"set_overall_enable\":1}}"))
}

// EditAntiTheftRule returns the JSON required to edit an anti-theft rule
func (kpp *KasaPowerPlug) EditAntiTheftRule(id string) ([]byte, error) {
	return kpp.talkToPlug(fmt.Sprintf("{\"anti_theft\":{\"edit_rule\":{\"stime_opt\":0,\"wday\":[0,0,0,1,0,1,0],\"smin\":987,\"enable\":1,\"frequency\":5,\"repeat\":1,\"etime_opt\":0,\"id\":\"%s\",\"duration\":2,\"name\":\"test\",\"lastfor\":1,\"month\":0,\"year\":0,\"longitude\":0,\"day\":0,\"latitude\":0,\"force\":0,\"emin\":1047},\"set_overall_enable\":1}}", id))
}

// DeleteAntiTheftRule returns the JSON required to delete an anti-theft rule with given ID
func (kpp *KasaPowerPlug) DeleteAntiTheftRule(id string) ([]byte, error) {
	return kpp.talkToPlug(fmt.Sprintf("{\"anti_theft\":{\"delete_rule\":{\"id\":\"%s\"}}}", id))
}

// DeleteAllAntiTheftRules is the JSON to delete all the anti-theft rules
func (kpp *KasaPowerPlug) DeleteAllAntiTheftRules(children ...int) ([]byte, error) {
	return kpp.talkToPlug(`{"anti_theft":{"delete_all_rules":}}`)
}

func trimJSONArray(s string) string {
	return strings.Replace(s, `,]`, `]`, 1)
}
