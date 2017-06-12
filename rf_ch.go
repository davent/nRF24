package nrf24

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

func (device *NRF24L01P) SetChannel(rf_ch byte) (err error) {

	device.Info("Setting RF frequency to %0.3fGHz\n", (2400.00+float64(rf_ch))/1000)
	err = device.writeRegister(RF_CH, RF_CH_MASK, rf_ch)

	return
}
