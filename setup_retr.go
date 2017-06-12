package nrf24

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
