package kasalink

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

const (
	plugIP = "10.0.0.4"
	//plugIP  = "10.0.0.5"
	plugPort = 9999
)

func TestKasaPowerPlug_DisableLED(t *testing.T) {
	var (
		kpp *KasaPowerPlug
		rw  *KasaResponse
		err error
	)
	kpp, err = NewKasaPowerPlug(fmt.Sprintf("%s:%d", plugIP, plugPort))
	if err != nil {
		t.Fatal(err)
	}
	rw, err = kpp.DisableLED()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", rw.System.SetLED)
}

func TestKasaPowerPlug_EnableLED(t *testing.T) {
	var (
		kpp KasaPowerPlug
		err error
		rw  *KasaResponse
	)
	kpp.plugNetworkLocation = fmt.Sprintf("%s:%d", plugIP, plugPort)
	rw, err = kpp.EnableLED()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", rw.System.SetLED)
}

func TestKasaTalkToChildrenNothingPassed(t *testing.T) {

	var (
		kpp       KasaPowerPlug
		jsonBytes []byte
		err       error
		rw        KasaResponse
	)
	kpp.plugNetworkLocation = fmt.Sprintf("%s:%d", plugIP, plugPort)
	jsonBytes, err = kpp.querySystemInfo()
	if err != nil {
		t.Fatal(err)
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	t.Logf("%s\n", jsonBytes)

	err = json.Unmarshal(jsonBytes, &rw)
	if err != nil {
		log.Printf("Err: %s", err)
		t.Fatal(err)
	}
	t.Logf("%+v\n", rw.System.GetSysInfo)
}

func TestKasaTalkToChildrenWithChildrenPassed(t *testing.T) {
	var (
		kpp KasaPowerPlug
		err error
		rw  *KasaResponse
	)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	kpp.plugNetworkLocation = "10.0.0.4:9999"
	//{"context":{"child_ids":["8006E92180ADBEA7B3E4820027152BE21ACC7D7700"]},"system":{"set_relay_state":{"state":1}}}
	//{"context":{"child_ids":["8006E92180ADBEA7B3E4820027152BE21ACC7D7700"]},"system":{"set_relay_state":{"state":1}}}
	kpp.deviceID = "8006E92180ADBEA7B3E4820027152BE21ACC7D77"
	rw, err = kpp.GetRealtimeCurrentAndVoltage()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", *rw.EnergyMeter.Realtime)
	for i := 0; i < 6; i++ {
		rw, err = kpp.GetRealtimeCurrentAndVoltage(i)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%+v\n", *rw.EnergyMeter.Realtime)
	}
}

func TestKasaPowerPlug_TurnDevice(t *testing.T) {
	var (
		kpp       KasaPowerPlug
		jsonBytes []byte
		//rw        KasaResponse
		err error
	)
	jsonBytes, err = kpp.TurnDeviceOff()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Turn Off Response: %s", jsonBytes)
	jsonBytes, err = kpp.TurnDeviceOn()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Turn On Response: %s", jsonBytes)

}
