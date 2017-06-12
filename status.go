package nrf24

type StatusRegister struct {
	RX_DR   uint8
	TX_DS   uint8
	MAX_RT  uint8
	RX_P_NO uint8
	TX_FULL uint8
}

func parseStatus(register byte) (status *StatusRegister, err error) {
	status = &StatusRegister{
		RX_DR:   uint8((register & RX_DR) >> 6),
		TX_DS:   uint8((register & TX_DS) >> 5),
		MAX_RT:  uint8((register & MAX_RT) >> 4),
		RX_P_NO: uint8((register & RX_P_NO) >> 1),
		TX_FULL: uint8(register & TX_FULL),
	}
	return
}

func (device *NRF24L01P) readStatusRegister() (status_register *StatusRegister, err error) {

	// read the register from the device
	var buf []byte
	buf, err = device.readRegister(STATUS)
	if err != nil {
		return
	}

	status_register, err = parseStatus(buf[0])
	if err != nil {
		return
	}

	return
}
