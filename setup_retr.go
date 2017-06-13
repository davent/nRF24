package nrf24

import (
	"fmt"
)

type SetupRetrRegister struct {
	ARD uint8
	ARC uint8
}

func (device *NRF24L01P) readSetupRetrRegister() (setup_retr_register *SetupRetrRegister, err error) {

	// read the register from the device
	var buf []byte
	buf, err = device.readRegister(SETUP_RETR)
	if err != nil {
		return
	}

	setup_retr_register = &SetupRetrRegister{
		ARD: uint8((buf[0] & ARD) >> 4),
		ARC: uint8(buf[0] & ARC),
	}

	return
}

func (device *NRF24L01P) SetAutoRetransmitCount(arc uint8) (err error) {

	device.Info("Setting Auto Retransmit Count to %d\n", arc)

	// arc must be between 0 (disabled) and 15
	if arc > 15 {
		fmt.Errorf("Invalid arc")
	}

	if err = device.writeRegister(SETUP_RETR, ARC, arc); err != nil {
		return
	}

	return
}
