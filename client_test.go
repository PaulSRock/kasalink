package kasalink

import (
	"log"
	"testing"
)

func TestGetDeviceInfo(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var response, err = TalkToPlug(KasaGetSystemInfo)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n", response)
}

func TestTurningOffAPlug(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var response, err = TalkToPlug(KasaTurnPlugOff("8006E92180ADBEA7B3E4820027152BE21ACC7D77", 2, 3, 5))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n", response)
}

func TestTurningOnAPlug(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var response, err = TalkToPlug(KasaTurnPlugOn("8006E92180ADBEA7B3E4820027152BE21ACC7D77", 5))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n", response)
}
