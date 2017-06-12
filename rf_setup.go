package nrf24

import (
	"fmt"
)

type RfSetupRegister struct {
	CONT_WAVE  uint8
	RF_DR_LOW  uint8
	PLL_LOCK   uint8
	RF_DR_HIGH uint8
	RF_PWR     uint8
}

func (device *NRF24L01P) readRfSetupRegister() (rf_setup_register *RfSetupRegister, err error) {

	// read the register from the device
	var buf []byte
	buf, err = device.readRegister(RF_SETUP)
	if err != nil {
		return
	}

	rf_setup_register = &RfSetupRegister{
		CONT_WAVE:  uint8((buf[0] & CONT_WAVE) >> 7),
		RF_DR_LOW:  uint8((buf[0] & RF_DR_LOW) >> 5),
		PLL_LOCK:   uint8((buf[0] & PLL_LOCK) >> 4),
		RF_DR_HIGH: uint8((buf[0] & RF_DR_HIGH) >> 3),
		RF_PWR:     uint8((buf[0] & RF_PWR) >> 1),
	}

	return

}

func (device *NRF24L01P) SetDataRate(rf_dr byte) (err error) {
	switch rf_dr {
	case RF_DR_1MBPS:
		device.Info("Setting data rate to 1Mbps")
		err = device.writeRegister(RF_SETUP, RF_DR_MASK, rf_dr)
	case RF_DR_2MBPS:
		device.Info("Setting data rate to 2Mbps")
		err = device.writeRegister(RF_SETUP, RF_DR_MASK, rf_dr)
	case RF_DR_250KBPS:
		device.Info("Setting data rate to 250kbps")
		err = device.writeRegister(RF_SETUP, RF_DR_MASK, rf_dr)
	default:
		err = fmt.Errorf("Unknown data rate specified")
	}
	return
}

func (device *NRF24L01P) SetPowerAmplifier(rf_pwr byte) (err error) {
	switch rf_pwr {
	case RF_PWR_18DBM:
		device.Info("Setting RF output power to -18dB")
		err = device.writeRegister(RF_SETUP, RF_PWR_MASK, rf_pwr)
	case RF_PWR_12DBM:
		device.Info("Setting RF output power to -12dB")
		err = device.writeRegister(RF_SETUP, RF_PWR_MASK, rf_pwr)
	case RF_PWR_6DBM:
		device.Info("Setting RF output power to -6dB")
		err = device.writeRegister(RF_SETUP, RF_PWR_MASK, rf_pwr)
	case RF_PWR_0DBM:
		device.Info("Setting RF output power to 0dB")
		err = device.writeRegister(RF_SETUP, RF_PWR_MASK, rf_pwr)
	default:
		err = fmt.Errorf("Unknown RF power output specified")
	}

	return
}
