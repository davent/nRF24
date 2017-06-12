package nrf24

type FifoStatusRegister struct {
	TX_REUSE uint8
	TX_FULL  uint8
	TX_EMPTY uint8
	RX_FULL  uint8
	RX_EMPTY uint8
}

func (device *NRF24L01P) readFifoStatusRegister() (fifo_status_register *FifoStatusRegister, err error) {

	// read the register from the device
	var buf []byte
	buf, err = device.readRegister(FIFO_STATUS)
	if err != nil {
		return
	}

	fifo_status_register = &FifoStatusRegister{
		TX_REUSE: uint8((buf[0] & TX_REUSE) >> 6),
		TX_FULL:  uint8((buf[0] & TX_FULL_FIFO) >> 5),
		TX_EMPTY: uint8((buf[0] & TX_EMPTY) >> 4),
		RX_FULL:  uint8((buf[0] & RX_FULL) >> 1),
		RX_EMPTY: uint8(buf[0] & RX_EMPTY),
	}

	return
}
