package nrf24

import (
	"fmt"
	"time"

	"github.com/davent/bcm2835"
)

func (device *NRF24L01P) enterPRX() (err error) {

	device.Debug("Entering PRX mode")

	// Set PRIM_RX to 1 in the CONFIG register
	err = device.writeRegister(CONFIG, PRIM_RX, PRIM_RX)
	if err != nil {
		return
	}

	// PULL CE_PIN HIGH to enter PRX mode
	bcm2835.GpioSet(CE_PIN)

	// Wait for TSTBY_2_A to enter RX mode
	time.Sleep(TSTBY_2_A)

	return
}

func (device *NRF24L01P) flushRxFifo() (err error) {

	device.Info("Flushing RX_FIFO")

	// Enter PRX mode
	if err = device.enterPRX(); err != nil {
		return
	}

	// Write new register back to the device
	write_buf := make([]byte, 2)
	read_buf := make([]byte, 2)
	write_buf[0] = FLUSH_RX
	if err = device.spiTx(write_buf, read_buf); err != nil {
		return
	}

	// Return to Standby-I mode
	if err = device.standbyI(); err != nil {
		return
	}

	return
}

func (device *NRF24L01P) readRxFifo() (read_buf []byte, err error) {

	// Get size of received packet
	var dpl uint8
	dpl, err = device.readRxPlWid()
	if err != nil {
		return
	}

	// Read data from the RX_FIFO buffer
	read_buf = make([]byte, dpl+1)
	write_buf := make([]byte, dpl+1)
	write_buf[0] = R_RX_PAYLOAD

	err = device.spiTx(write_buf, read_buf)
	if err != nil {
		return
	}

	return read_buf[1:], nil
}

func (device *NRF24L01P) readRxPlWid() (dpl uint8, err error) {

	// Read DPL from the R_RX_PL_WID register
	read_buf := make([]byte, 2)
	write_buf := make([]byte, 2)
	write_buf[0] = R_RX_PL_WID

	err = device.spiTx(write_buf, read_buf)
	if err != nil {
		return
	}

	dpl = uint8(read_buf[1])

	device.Info("Received data Dynamic Packet Length: %d\n", dpl)

	// Flush RX FIFO if the read value is larger than 32 bytes
	if dpl > 32 {
		if err = device.flushRxFifo(); err != nil {
			return
		}

		return dpl, fmt.Errorf("Received data Dynamic Packet Length was greater than 32 bytes")
	}

	return
}

func (device *NRF24L01P) resetRxDr() (err error) {

	device.Info("Resetting Data Ready RX FIFO interrupt")
	err = device.writeRegister(STATUS, RX_DR, RX_DR)

	return
}

func (device *NRF24L01P) Receive() (data_chan chan []byte, err error) {

	data_chan = make(chan []byte)
	err_chan := make(chan error)

	// Get the ARD (TX retry delay) to use as interval to check for new RX data
	// - To be replaced with IRQ in the future
	_, ard, err := device.txAckTimeout()
	if err != nil {
		return
	}

	go func() {

		// Enter PRX mode
		if err = device.enterPRX(); err != nil {
			return
		}

		// wait for data to arrive
		for {

			// Perform a NOP registry lookup to get status
			var status *StatusRegister
			status, err = device.readStatusRegister()
			if err != nil {
				err_chan <- err
			}

			if status.RX_DR == 1 {
				var data []byte
				data, err = device.readRxFifo()
				if err != nil {
					err_chan <- err
				}

				device.Info("Received data: %s\n", data)

				// Add incoming data to the data channel
				data_chan <- data

				// Reset Data Ready RX FIFO interrupt
				if err = device.resetRxDr(); err != nil {
					err_chan <- err
				}
			}

			// Try to save some CPU cycles
			time.Sleep(ard)
		}

		// Return to Standby-I mode
		if err = device.standbyI(); err != nil {
			return
		}

	}()

	return
}
