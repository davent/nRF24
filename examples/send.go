package main

import (
	"log"

	"github.com/davent/nRF24"
)

const (
	APPLICATION_NAME string = "nrf24L01+ Send example"
)

func main() {

	// Welcome children, are you sitting comfortably?
	log.Printf("%s (%s)\n", APPLICATION_NAME, nrf24.VERSION)

	// Then we shall begin.

	// First create a new radio device
	// - Just pass the path to the Linux kernel module provided SPI device
	radio, err := nrf24.New("/dev/spidev0.0")
	if err != nil {
		panic(err)
	}
	defer radio.Close()

	// Set log level to help with debug
	// - CRITICAL
	// - ERROR (default)
	// - WARNING
	// - INFO
	// - DEBUG
	radio.LogLevel = nrf24.INFO

	// Set the data to eiter
	// - 2Mbps (RF_DR_2MBPS)
	// - 1Mbps (RF_DR_1MBPS)
	// - 512kbps (RF_DR_512KBPS)
	radio.SetDataRate(nrf24.RF_DR_2MBPS)

	// Set the RF frequency
	// - This can be any frequency from 2.4GHz to 2.525GHz in 2MHz divisions
	// - RF_CH_2400MHZ, RF_CH_2402MHZ, RF_CH_2404MHZ ... RF_CH_2525MHZ
	radio.SetChannel(nrf24.RF_CH_2400MHZ)

	// Set the power output of the transmitting radio (from most power to least powerful)
	// -  0dB  (RF_PWR_0DBM)
	// - -6dB  (RF_PWR_6DBM)
	// - -12dB (RF_PWR_12DBM)
	// - -18dB (RF_PWR_18DBM)
	radio.SetPowerAmplifier(nrf24.RF_PWR_18DBM)

	// Set the CRC accuracy
	// - 8bit  (CRC0_8BIT)
	// - 16bit (CRC0_16BIT)
	radio.SetCRCO(nrf24.CRC0_16BIT)

	// Now we are ready to send! Turn the power to the radio on
	// - true  = on
	// - false = off
	radio.PowerOn(true)

	// Now, let's send our "Hello world!" messages
	// - 'messages' is []byte of up to 32 bytes in length
	message := []byte("Hello world!")
	if err = radio.Send(message); err != nil {
		log.Println(err)
	}

	// Don't forget to turn the radio off before you leave!
	radio.PowerOn(false)

}
