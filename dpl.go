package nrf24

func (device *NRF24L01P) enableDPL() (err error) {

	// To enable DPL (Dynamic Payload Length), we need to
	// - Set the EN_DPL bit to 1 in the FEATURE register
	// - Set the DPL_P{n} bit to 1 for each pip in the DYNPD Register

	// Enable Dynamic Payload Length feature
	if err = device.writeRegister(FEATURE, EN_DPL, EN_DPL); err != nil {
		return
	}

	// Enable Dynamic Payload Length for pipe 0
	if err = device.writeRegister(DYNPD, DPL_P0, DPL_P0); err != nil {
		return
	}

	// Enable Dynamic Payload Length for pipe 1
	if err = device.writeRegister(DYNPD, DPL_P1, DPL_P1); err != nil {
		return
	}

	// Enable Dynamic Payload Length for pipe 2
	if err = device.writeRegister(DYNPD, DPL_P2, DPL_P2); err != nil {
		return
	}

	// Enable Dynamic Payload Length for pipe 3
	if err = device.writeRegister(DYNPD, DPL_P3, DPL_P3); err != nil {
		return
	}

	// Enable Dynamic Payload Length for pipe 4
	if err = device.writeRegister(DYNPD, DPL_P4, DPL_P4); err != nil {
		return
	}

	// Enable Dynamic Payload Length for pipe 5
	if err = device.writeRegister(DYNPD, DPL_P5, DPL_P5); err != nil {
		return
	}

	return
}
