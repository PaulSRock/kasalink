package kasalink

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

// KasaPowerPlug is the struct that holds info about and methods for talking to a Kasa Power Plug or Power Strip
type KasaPowerPlug struct {
	plugNetworkLocation string
	Unsafe              unsafe
	deviceID            string
	tplinkClient        net.Conn
	timeout             time.Duration
	SysInfo             *SystemInfo
	log                 *log.Logger
	debug               bool
}

// NewKasaPowerPlug gives you a new KasaPowerPlug struct that's already gotten it's system info, or an error
// telling you why that didn't work
func NewKasaPowerPlug(plugAddress string) (kpp *KasaPowerPlug, err error) {
	kpp = &KasaPowerPlug{
		plugNetworkLocation: plugAddress,
		timeout:             5 * time.Second,
	}
	kpp.SysInfo, err = kpp.GetSystemInfo()
	if err != nil {
		return nil, err
	}
	kpp.deviceID = kpp.SysInfo.DeviceID
	return
}

// SetLogger lets you give KasaPowerPlug a place to send logs. If not set, logs will disappear into the void.
func (kpp *KasaPowerPlug) SetLogger(newLogger *log.Logger) {
	kpp.log = newLogger
}

// TalkToPlug sends a command to the plug and returns a response json and error error
func (kpp *KasaPowerPlug) talkToPlug(KasaCommand string) (response []byte, err error) {
	var (
		bitsToSend []byte
		bitsWeRead []byte
	)

	if kpp.tplinkClient == nil {
		if kpp.timeout == 0 {
			kpp.timeout = time.Duration(10) * time.Second
		}
		if kpp.tplinkClient, err = net.DialTimeout("tcp", kpp.plugNetworkLocation, kpp.timeout); err != nil {
			return
		}
	}

	// so you would think that a defer, close here would be good, but it can cause a race condition when there are
	// multiple successive calls to the power plug - the if state above can eval as true, then the defer closes
	// the tcp connection, and we end with a "use of closed network connection" error
	//defer kpp.closer()

	if err = kpp.tplinkClient.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return nil, err
	}

	bitsToSend = encrypt(KasaCommand)
	if _, err = kpp.tplinkClient.Write(bitsToSend); err != nil {
		return
	}

	if bitsWeRead, err = kpp.readResponse(); err != nil {
		return
	}

	bitsWeRead = decrypt(bitsWeRead)
	if kpp.debug && kpp.log != nil {
		kpp.log.Printf("Received:\n%s\n", bitsWeRead)
	}
	return bitsWeRead, nil
}

func (kpp *KasaPowerPlug) readResponse() ([]byte, error) {
	var (
		bodySize uint32
		err      error
	)
	err = binary.Read(kpp.tplinkClient, binary.BigEndian, &bodySize)
	if err != nil {
		return nil, err
	}
	var buf = make([]byte, bodySize)

	_, err = io.ReadAtLeast(kpp.tplinkClient, buf, int(bodySize))
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// tellChild is the JSON used to issue a command to individual sockets on a Kasa enabled device
func (kpp *KasaPowerPlug) tellChild(cmd string, children ...int) ([]byte, error) {
	var (
		sb  strings.Builder
		err error
	)

	if _, err = sb.WriteString(`{"context":{"child_ids":[`); err != nil {
		return nil, err
	}
	for _, child := range children {
		if _, err = sb.WriteString(fmt.Sprintf(`"%s%02d",`, kpp.deviceID, child)); err != nil {
			return nil, err
		}
	}
	if _, err = sb.WriteString(fmt.Sprintf(`]},%s}`, cmd[1:len(cmd)-1])); err != nil {
		return nil, err
	}
	//log.Printf("Child Call: %s\n", sb.String())
	//log.Printf("Child Call Trimmed: %s\n", trimJSONArray(sb.String()))
	return kpp.talkToPlug(trimJSONArray(sb.String()))
}

// Close tells the client to close any active connection it might have to the power strip/plug
func (kpp *KasaPowerPlug) Close() error {
	if kpp.tplinkClient != nil {
		return kpp.tplinkClient.Close()
	}
	// the net.Conn object is nil, so nothing to close, return nil
	return nil
}

func (kpp *KasaPowerPlug) closer() {
	var err error
	if err = kpp.Close(); err != nil {
		if kpp.log != nil {
			kpp.log.Println("Error closing down tcp client: ", err)
		}
	}
}
