package reefpihal

import (
	"github.com/PaulSRock/kasalink"
	"sync"
	"time"

	"github.com/reef-pi/hal"
)

type hs300ChildInfo struct {
	lastUpdate time.Time
	kasalink.KasaResponse
}

// HS300 is the main struct to control a Kasa HS300 Smart Plug via the Reef-Pi Output Plug interface.
type HS300 struct {
	kpp       *kasalink.KasaPowerPlug
	childInfo []hs300ChildInfo

	sync.RWMutex
}

// ADCChannels implements the ADCChannels call in the hal.ADCDriver interface
func (h *HS300) ADCChannels() []hal.ADCChannel {
	var adcs []hal.ADCChannel
	adcs = make([]hal.ADCChannel, h.kpp.SysInfo.ChildNum*4)
	for i := 0; i < h.kpp.SysInfo.ChildNum*4; i++ {
		adcs[i] = &HS300ADCChannel{
			hs300: h,
			id:    i % h.kpp.SysInfo.ChildNum,
			name:  h.kpp.SysInfo.Children[i].Alias,
			pm:    powerMetric(i % 4),
		}
	}
	return adcs
}

// ADCChannel implements ADCChannel call in the hal.ADCDriver interface
func (h *HS300) ADCChannel(i int) (hal.ADCChannel, error) {
	return &HS300ADCChannel{
		hs300: h,
		id:    i % h.kpp.SysInfo.ChildNum,
		name:  h.kpp.SysInfo.Children[i].Alias,
		pm:    powerMetric(i % 4),
	}, nil
}

// Close returns nil, this is a remote service call
func (h *HS300) Close() error {
	return h.kpp.Close()
}

// Metadata returns the description for this driver and the fact that we take input, and report output.
func (h *HS300) Metadata() hal.Metadata {
	return hal.Metadata{
		Capabilities: []hal.Capability{hal.Output, hal.Input},
		Name:         h.kpp.SysInfo.Alias,
		Description:  "Driver for a TP-Link/Kasa power Strip",
	}
}

// OutputPins gives you back an array to all of the available plugs on the Kasa device
func (h *HS300) OutputPins() []hal.OutputPin {
	var ops []hal.OutputPin
	ops = make([]hal.OutputPin, h.kpp.SysInfo.ChildNum)
	for i := 0; i < h.kpp.SysInfo.ChildNum; i++ {
		ops[i] = &HS300OutputPin{childID: i, hs300: h}
	}
	return ops
}

// OutputPin Does something, not totally sure what yet.
func (h *HS300) OutputPin(id int) (hal.OutputPin, error) {
	return &HS300OutputPin{hs300: h, childID: id}, nil
}

// NewHS300 takes an IP address (as a string) and builds a new HS300 struct and initializes things so it can talk to
// the specified Kasa HS300 Power Strip
func NewHS300(kppAddress string) (*HS300, error) {
	var kpp, err = kasalink.NewKasaPowerPlug(kppAddress)
	if err != nil {
		return nil, err
	}
	var hs300 = HS300{kpp: kpp}
	return &hs300, nil
}
