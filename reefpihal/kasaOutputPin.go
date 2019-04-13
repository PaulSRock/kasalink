package reefpihal

import "fmt"

const errPinNotInitialized = "kasa device not initialized"

//HS300OutputPin is the hal.OutputPin implementation for the Kasa HS300 PowerStrip
type HS300OutputPin struct {
	childID int
	hs300   *HS300
}

// Close does nothing, exists only because the hal.Pin interface requires it.
func (h *HS300OutputPin) Close() error {
	return nil
}

// Name returns the Alias that has been set for the plug in question
func (h *HS300OutputPin) Name() string {
	if h.hs300.kpp.SysInfo != nil {
		return h.hs300.kpp.SysInfo.Children[h.childID].Alias
	}
	return ""
}

// Write changes the device state between on and off.
func (h *HS300OutputPin) Write(state bool) error {
	if h.hs300.kpp.SysInfo == nil {
		return fmt.Errorf(errPinNotInitialized)
	}
	if state {
		_, err := h.hs300.kpp.TurnDeviceOn(h.childID)
		if err != nil {
			return err
		}
	} else {
		_, err := h.hs300.kpp.TurnDeviceOff(h.childID)
		if err != nil {
			return err
		}
	}
	return nil
}

// LastState returns the last reported state of the outlet, unless the HS300 hasn't been properly initialized,
// and then it will always return false.
func (h *HS300OutputPin) LastState() bool {
	if h.hs300.kpp.SysInfo == nil {
		return false
	}
	return h.hs300.kpp.SysInfo.Children[h.childID].State == 1
}
