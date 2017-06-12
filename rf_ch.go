package nrf24

import (
	"fmt"
)

type RfChRegister struct {
	RF_CH uint8
}

func (device *NRF24L01P) readRfChRegister() (rf_ch_register *RfChRegister, err error) {

	// read the register from the device
	var buf []byte
	buf, err = device.readRegister(RF_CH)
	if err != nil {
		return
	}

	rf_ch_register = &RfChRegister{
		RF_CH: uint8(buf[0]),
	}

	return
}

func (device *NRF24L01P) setChannel(rf_ch byte) (err error) {

	device.Info("Setting RF frequency to %0.3fGHz\n", (float64(RF_CH_BASE)+float64(rf_ch))/1000)
	err = device.writeRegister(RF_CH, RF_CH_MASK, rf_ch)

	return
}

func (device *NRF24L01P) SetFrequency(freq uint) (err error) {

	// freq must be between 2400 and 2525 MHz
	if (freq < 2400) || (freq > 2525) {
		return fmt.Errorf("Invalid frequency")
	}

	rf_ch := byte(freq - RF_CH_BASE)
	if err = device.setChannel(rf_ch); err != nil {
		return
	}

	return
}
