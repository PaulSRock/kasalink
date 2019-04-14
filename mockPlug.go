package kasalink

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

// MockPlug is for running unit tests, it'll fake responses as if it's an actual plug (eventually)
type MockPlug struct {
	net.Conn
	ln       net.Listener
	lastSent string
}

// NewMockPlug gives you a new MockPlug with a running TCP Server instance to handle request. As written, a MockPlug
// will only handle a single command and exit.
func NewMockPlug() (mp MockPlug, err error) {

	mp.ln, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return
	}
	mp.Conn, err = net.Dial("tcp", mp.ln.Addr().String())
	if err != nil {
		return
	}
	go func() {
		defer func() {
			err = mp.ln.Close()
			if err != nil {
				log.Println("Error trying to close out new mock plug listener:", err)
			}
		}()
		myConn, err := mp.ln.Accept()
		if err != nil {
			return
		}
		mp.ConnectionHandler(myConn)
	}()
	return
}

// DialMe returns a connection to the MockPlug's internal server so you can send it requests and receive answers (eventually)
func (m *MockPlug) DialMe() (net.Conn, error) {
	return net.Dial("tcp", m.ln.Addr().String())
}

// ConnectionHandler handles the connection when something connects to the MockPlug and sends a command.
// If it doesn't send a supported command (and I've only actually implemented a few), bad things may occur.
func (m *MockPlug) ConnectionHandler(myConnection net.Conn) {

	var bodySize uint32
	err := binary.Read(myConnection, binary.BigEndian, &bodySize)
	if err != nil {
		return
	}
	var buf = make([]byte, bodySize)

	_, err = io.ReadAtLeast(myConnection, buf, int(bodySize))
	if err != nil {
		return
	}
	var indexString string
	var clearBits = decrypt(buf)
	log.Printf("Got the following: %s", clearBits)
	if bytes.Contains(clearBits, []byte(`"context":{"child_ids":["`)) {
		indexString = fmt.Sprintf("{%s", clearBits[bytes.Index(clearBits, []byte(`"]},`))+4:])
	} else {
		indexString = string(clearBits)
	}
	//log.Println("indexString:", indexString)
	var cmdMap = map[string]string{
		getSysInfo:           `{"system": {"get_sysinfo": {"sw_ver":"1.0.6 Build 180627 Rel.081000", "hw_ver":"1.0", "model":"HS300(US)", "deviceId":"8006E92180ADBEA7B3E4820027152BE21ACC7D77", "oemId":"5C9E6254BEBAED63B2B6102966D24C17", "hwId":"34C41AA028022D0CCEA5E678E8547C54", "rssi":-35, "longitude_i":-775702, "latitude_i":391156, "alias":"TP-LINK_Power Strip_14A9", "mic_type":"IOT.SMARTPLUGSWITCH", "feature":"TIM:ENE", "mac":"B0:BE:76:80:14:A9", "updating":0, "led_off":0, "children":[ {"id":"8006E92180ADBEA7B3E4820027152BE21ACC7D7700", "state":1, "alias":"Top Tank Light", "on_time":394931, "next_action":{"type":-1} }, {"id":"8006E92180ADBEA7B3E4820027152BE21ACC7D7701", "state":1, "alias":"Top Tank Heater", "on_time":1200805, "next_action":{"type":-1} }, {"id":"8006E92180ADBEA7B3E4820027152BE21ACC7D7702", "state":1, "alias":"Top Tank Filter", "on_time":50621, "next_action":{"type":-1} }, {"id":"8006E92180ADBEA7B3E4820027152BE21ACC7D7703", "state":1, "alias":"Top Tank Powerhead", "on_time":385855, "next_action":{"type":-1} }, {"id":"8006E92180ADBEA7B3E4820027152BE21ACC7D7704", "state":1, "alias":"Air Pump", "on_time":1200806, "next_action":{"type":-1} }, {"id":"8006E92180ADBEA7B3E4820027152BE21ACC7D7705", "state":1, "alias":"Plug 6", "on_time":1289431, "next_action":{"type":-1} }], "child_num":6, "err_code":0 } }}`,
		turnOffLED:           `{"system":{"error":0}}`,
		turnOnLED:            `{"system":{"error":0}}`,
		getCurrentAndVoltage: `{"emeter":{"get_realtime":{"voltage_mv":121122,"current_ma":34,"power_mw":2079,"total_wh":3376,"err_code":0}}}`,
		turnOn:               `{"system":{"set_relay_state":{"err_code":0}}}`,
		turnOff:              `{"system":{"set_relay_state":{"err_code":0}}}`,
	}
	response, ok := cmdMap[indexString]
	if ok {
		_, err = myConnection.Write(encrypt(response))
		if err != nil {
			return
		}
	} else {
		_, err = myConnection.Write(encrypt(`{"system":{"error":1}}`))
		if err != nil {
			return
		}
	}
	return
}
