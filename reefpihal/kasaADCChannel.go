package reefpihal

import (
	"fmt"
	"time"

	"github.com/PaulSRock/kasalink"

	"github.com/reef-pi/hal"
)

// PowerMe
type powerMetric int

const (
	voltage powerMetric = iota
	current
	power
	totalWatts
)

// HS300ADCChannel implements hal.ADCChannel for the HS300 Power Strip
type HS300ADCChannel struct {
	hs300 *HS300
	name  string
	id    int
	pm    powerMetric
}

// Name returns the name of the power metric this channel provides
func (h *HS300ADCChannel) Name() string {
	switch h.pm {
	case voltage:
		return fmt.Sprintf("%s voltage", h.name)
	case current:
		return fmt.Sprintf("%s current", h.name)
	case power:
		return fmt.Sprintf("%s power", h.name)
	case totalWatts:
		return fmt.Sprintf("%s total watts", h.name)
	default:
		return h.name
	}
}

// Read asks the HS300 power strip for updated metrics and returns those values.
func (h *HS300ADCChannel) Read() (float64, error) {
	var (
		rw  *kasalink.KasaResponse
		err error
	)
	if time.Now().After(h.hs300.childInfo[h.id].lastUpdate.Add(time.Second)) {
		rw, err = h.hs300.kpp.GetRealtimeCurrentAndVoltage(h.id)
		h.hs300.childInfo[h.id].lastUpdate = time.Now()
		if err != nil {
			return 0, err
		}
		h.hs300.childInfo[h.id].EnergyMeter = rw.EnergyMeter
	}
	if rw.EnergyMeter.Realtime.Error != 0 {
		return 0, fmt.Errorf("error gathering power stats from plug")
	}
	switch h.pm {
	case voltage:
		return float64(rw.EnergyMeter.Realtime.Voltage), err
	case current:
		return float64(rw.EnergyMeter.Realtime.Current), err
	case power:
		return float64(rw.EnergyMeter.Realtime.Power), err
	case totalWatts:
		return float64(rw.EnergyMeter.Realtime.Power), err
	default:
		return 0, err
	}
}

// Calibrate is not currently supported
func (h *HS300ADCChannel) Calibrate([]hal.Measurement) error {
	return nil
}

// Measure is the same as Read (it literally calls the Read method)
func (h *HS300ADCChannel) Measure() (float64, error) {
	return h.Read()
}
