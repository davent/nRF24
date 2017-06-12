package nrf24

import (
	"fmt"
)

func (device *NRF24L01P) PrintRegisters() {

	// Status
	status, err := device.readStatusRegister()
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
	} else {
		fmt.Println("======\nSTATUS\n======")
		fmt.Printf("RX_DR: %d\n", status.RX_DR)
		fmt.Printf("TX_DS: %d\n", status.TX_DS)
		fmt.Printf("MAX_RT: %d\n", status.MAX_RT)
		fmt.Printf("RX_P_NO: %d\n", status.RX_P_NO)
		fmt.Printf("TX_FULL: %d\n", status.TX_FULL)
		fmt.Println()
	}

	// Config
	config, err := device.readConfigRegister()
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
	} else {
		fmt.Println("======\nCONFIG\n======")
		fmt.Printf("MASK_RX_DR: %d\n", config.MASK_RX_DR)
		fmt.Printf("MASK_TX_DS: %d\n", config.MASK_TX_DS)
		fmt.Printf("MASK_MAX_RT: %d\n", config.MASK_MAX_RT)
		fmt.Printf("EN_CRC: %d\n", config.EN_CRC)
		fmt.Printf("CRCO: %d\n", config.CRCO)
		fmt.Printf("PWR_UP: %d\n", config.PWR_UP)
		fmt.Printf("PRIM_RX: %d\n", config.PRIM_RX)
		fmt.Println()
	}

	// RF CH
	rf_ch, err := device.readRfChRegister()
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
	} else {
		fmt.Println("=====\nRF_CH\n=====")
		fmt.Printf("RF_CH: %d\n", rf_ch.RF_CH)
		fmt.Println()

	}
	// RF Setup
	rf_setup, err := device.readRfSetupRegister()
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
	} else {
		fmt.Println("========\nRF_SETUP\n========")
		fmt.Printf("CONT_WAVE: %d\n", rf_setup.CONT_WAVE)
		fmt.Printf("RF_DR_LOW: %d\n", rf_setup.RF_DR_LOW)
		fmt.Printf("PLL_LOCK: %d\n", rf_setup.PLL_LOCK)
		fmt.Printf("RF_DR_HIGH: %d\n", rf_setup.RF_DR_HIGH)
		fmt.Printf("RF_PWR: %d\n", rf_setup.RF_PWR)
		fmt.Println()
	}

}
