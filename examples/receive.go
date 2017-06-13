package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/davent/nRF24"
)

const (
	APPLICATION_NAME string = "nRF24L01+ Receive example"
)

func main() {

	// Welcome children, are you sitting comfortably?
	log.Printf("%s (%s)\n", APPLICATION_NAME, nrf24.VERSION)

	// Then we shall begin.

	// We need to catch signals so that we can close down cleanly
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	// First create a new radio device
	// - Just pass the path to the Linux kernel module provided SPI device
	radio, err := nrf24.New("/dev/spidev0.0")
	if err != nil {
		panic(err)
	}
	defer radio.Close()

	// If a signal is received, ensure the device is powered down gracefully
	go func() {
		<-signalChannel

		if err = radio.PowerOn(false); err != nil {
			log.Printf("Could not power off: %s\n", err)
		}

		radio.Close()
		os.Exit(0)
	}()

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
	// - 250kbps (RF_DR_250KBPS)
	if err = radio.SetDataRate(nrf24.RF_DR_250KBPS); err != nil {
		log.Printf("Culd not set data rate: %s\n", err)
	}

	// Set the RF frequency
	// - This can be any frequency from 2.4GHz to 2.525GHz specified in MHz
	if err = radio.SetFrequency(2400); err != nil {
		log.Printf("Could not set frequency: %s\n", err)
	}

	// Set the power output of the transmitting radio (from most power to least powerful)
	// -  0dB  (RF_PWR_0DBM)
	// - -6dB  (RF_PWR_6DBM)
	// - -12dB (RF_PWR_12DBM)
	// - -18dB (RF_PWR_18DBM)
	if err = radio.SetPowerAmplifier(nrf24.RF_PWR_0DBM); err != nil {
		log.Printf("Could not set power amplifer: %s\n", err)
	}

	// Set the CRC accuracy
	// - 8bit  (CRC0_8BIT)
	// - 16bit (CRC0_16BIT)
	if err = radio.SetCRCO(nrf24.CRC0_16BIT); err != nil {
		log.Printf("Could not set CRC: %s\n", err)
	}

	// Now we are ready to send! Turn the power to the radio on
	// - true  = on
	// - false = off
	if err = radio.PowerOn(true); err != nil {
		log.Printf("Could not power on: %s\n", err)
	}

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
	if err = radio.PowerOn(false); err != nil {
		log.Printf("Could not power off: %s\n", err)
	}

}
