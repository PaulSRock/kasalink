package kasalink

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
)

const (
	plugIP = "10.0.0.4"
	//plugIP  = "10.0.0.5"
	plugPort = 9999
)

var (
	useMock     bool
	debugLogger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
)

func init() {

	flag.BoolVar(&useMock, "useMock", true, "use the MockPlug instead of the real one.")
	flag.Parse()
}

func mockOrNot(kpp **KasaPowerPlug, t *testing.T) {
	var err error
	if useMock {
		*kpp = new(KasaPowerPlug)
		(*kpp).tplinkClient, err = NewMockPlug()
		if err != nil {
			t.Fatal(err)
		}
		_, err = (*kpp).GetSystemInfo()
		if err != nil {
			t.Fatal(err)
		}
		(*kpp).tplinkClient, err = NewMockPlug()
		if err != nil {
			t.Fatal(err)
		}
	} else {
		*kpp, err = NewKasaPowerPlug(fmt.Sprintf("%s:%d", plugIP, plugPort))
		if err != nil {
			t.Fatal(err)
		}
		(*kpp).SetLogger(debugLogger)
		(*kpp).debug = true
	}

}

func TestKasaPowerPlug_DisableLED(t *testing.T) {
	var (
		kpp *KasaPowerPlug
		rw  *KasaResponse
		err error
	)
	mockOrNot(&kpp, t)
	rw, err = kpp.DisableLED()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", rw.System.SetLED)
}

func TestKasaPowerPlug_EnableLED(t *testing.T) {
	var (
		kpp *KasaPowerPlug
		err error
		rw  *KasaResponse
	)
	mockOrNot(&kpp, t)
	rw, err = kpp.EnableLED()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", rw.System.SetLED)
}

func TestKasaTalkToChildrenNothingPassed(t *testing.T) {

	var (
		kpp       *KasaPowerPlug
		jsonBytes []byte
		err       error
		rw        KasaResponse
	)
	mockOrNot(&kpp, t)
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
		kpp *KasaPowerPlug
		err error
		rw  *KasaResponse
	)
	//useMock = false
	mockOrNot(&kpp, t)
	rw, err = kpp.GetRealtimeCurrentAndVoltage()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", *rw.EnergyMeter.Realtime)
	for i := 0; i < 6; i++ {
		mockOrNot(&kpp, t)
		rw, err = kpp.GetRealtimeCurrentAndVoltage(i)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%+v\n", *rw.EnergyMeter.Realtime)
	}
}

func TestKasaPowerPlug_TurnDevice(t *testing.T) {
	var (
		kpp       *KasaPowerPlug
		jsonBytes []byte
		//rw        KasaResponse
		err error
	)
	//useMock = false
	mockOrNot(&kpp, t)
	jsonBytes, err = kpp.TurnDeviceOff()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Turn Off Response: %s", jsonBytes)
	if !bytes.Contains(jsonBytes, []byte(`"err_code":0`)) {
		t.Fatal("Error response from the TPLink Device")
	}

	mockOrNot(&kpp, t)
	jsonBytes, err = kpp.TurnDeviceOn()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Turn On Response: %s", jsonBytes)
	if !bytes.Contains(jsonBytes, []byte(`"err_code":0`)) {
		t.Fatal("Error response from the TPLink Device")
	}
}
