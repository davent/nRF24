package nrf24

import (
	"fmt"
	"time"

	"github.com/davent/bcm2835"
)

type ConfigRegister struct {
	MASK_RX_DR  uint8
	MASK_TX_DS  uint8
	MASK_MAX_RT uint8
	EN_CRC      uint8
	CRCO        uint8
	PWR_UP      uint8
	PRIM_RX     uint8
}

func (device *NRF24L01P) readConfigRegister() (config_register *ConfigRegister, err error) {

	// read the register from the device
	var buf []byte
	buf, err = device.readRegister(CONFIG)
	if err != nil {
		return
	}

	config_register = &ConfigRegister{
		MASK_RX_DR:  uint8((buf[0] & MASK_RX_DR) >> 6),
		MASK_TX_DS:  uint8((buf[0] & MASK_TX_DS) >> 5),
		MASK_MAX_RT: uint8((buf[0] & MASK_MAX_RT) >> 4),
		EN_CRC:      uint8((buf[0] & EN_CRC) >> 3),
		CRCO:        uint8((buf[0] & CRCO) >> 2),
		PWR_UP:      uint8((buf[0] & PWR_UP) >> 1),
		PRIM_RX:     uint8(buf[0] & PRIM_RX),
	}

	return
}

func (device *NRF24L01P) standbyI() (err error) {

	device.Debug("Entering Standby-I mode")

	// PULL CE_PIN LOW to leave PRX/PTX mode
	bcm2835.GpioClr(CE_PIN)

	// Set PRIM_RX to 0 in the CONFIG register
	err = device.writeRegister(CONFIG, PRIM_RX, 0x00)
	if err != nil {
		return
	}

	time.Sleep(TSTBY_2_A)

	return
}

func (device *NRF24L01P) SetCRCO(crc0 byte) (err error) {
	switch crc0 {
	case CRC0_8BIT:
		device.Info("Setting CRC to 8bit")
		err = device.writeRegister(CONFIG, CRC0_16BIT, crc0)
	case CRC0_16BIT:
		device.Info("Setting CRC to 16bit")
		err = device.writeRegister(CONFIG, CRC0_16BIT, crc0)
	default:
		err = fmt.Errorf("Unknown CRCO specified")
	}

	return
}

func (device *NRF24L01P) PowerOn(on bool) (err error) {

	if on {
		// set PWR_UP to 1 in the CONFIG register
		if err = device.writeRegister(CONFIG, PWR_UP, PWR_UP); err != nil {
			return
		}

		// Enter Standby-I mode
		if err = device.standbyI(); err != nil {
			return
		}
	} else {
		device.Info("Enabling Power Down mode")

		// set PWR_UP to 1 in the CONFIG register
		if err = device.writeRegister(CONFIG, PWR_UP, 0x00); err != nil {
			return
		}
	}

	time.Sleep(TPD_2_STBY)
	return
}
