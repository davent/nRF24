package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/davent/nrf24l01p"
)

const (
	APPLICATION_NAME string = "nRF24L01+ Receive example"
)

func main() {

	// Welcome children, are you sitting comfortably?
	log.Printf("%s (%s)\n", APPLICATION_NAME, nrf24l01p.VERSION)

	// Then we shall begin.

	// We need to catch signals so that we can close down cleanly
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	// First create a new radio device
	// - Just pass the path to the Linux kernel module provided SPI device
	radio, err := nrf24l01p.New("/dev/spidev0.0")
	if err != nil {
		panic(err)
	}
	defer radio.Close()

	// If a signal is received, ensure the device is powered down gracefully
	go func() {
		<-signalChannel
		radio.PowerOn(false)
		radio.Close()
		os.Exit(0)
	}()

	// Set log level to help with debug
	// - CRITICAL
	// - ERROR (default)
	// - WARNING
	// - INFO
	// - DEBUG
	radio.LogLevel = nrf24l01p.INFO

	// Set the data to eiter
	// - 2Mbps (RF_DR_2MBPS)
	// - 1Mbps (RF_DR_1MBPS)
	// - 512kbps (RF_DR_512KBPS)
	radio.SetDataRate(nrf24l01p.RF_DR_2MBPS)

	// Set the RF frequency
	// - This can be any frequency from 2.4GHz to 2.525GHz in 2MHz divisions
	// - RF_CH_2400MHZ, RF_CH_2402MHZ, RF_CH_2404MHZ ... RF_CH_2525MHZ
	radio.SetChannel(nrf24l01p.RF_CH_2400MHZ)

	// Set the power output of the transmitting radio (from most power to least powerful)
	// -  0dB  (RF_PWR_0DBM)
	// - -6dB  (RF_PWR_6DBM)
	// - -12dB (RF_PWR_12DBM)
	// - -18dB (RF_PWR_18DBM)
	radio.SetPowerAmplifier(nrf24l01p.RF_PWR_18DBM)

	// Set the CRC accuracy
	// - 8bit  (CRC0_8BIT)
	// - 16bit (CRC0_16BIT)
	radio.SetCRCO(nrf24l01p.CRC0_16BIT)

	// Now we are ready to send! Turn the power to the radio on
	// - true  = on
	// - false = off
	radio.PowerOn(true)

	// Get the incoming data channel
	incoming, err := radio.Receive()
	if err != nil {
		log.Println(err)
	} else {

		// Wait for messages to be received
		var message []byte
		for {
			message = <-incoming
			log.Printf("Message received: %s\n", message)
		}

	}

	// Don't forget to turn the radio off before you leave!
	radio.PowerOn(false)

}
