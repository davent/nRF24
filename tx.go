package nrf24

import (
	"fmt"
	"time"

	"github.com/davent/bcm2835"
)

func (device *NRF24L01P) enterPTX() (err error) {

	device.Debug("Entering PTX mode")

	// Set PRIM_RX to 0 in the CONFIG register
	err = device.writeRegister(CONFIG, PRIM_RX, 0x00)
	if err != nil {
		return
	}

	return
}

func (device *NRF24L01P) flushTxFifo() (err error) {

	device.Info("Flushing TX_FIFO")

	// Enter PTX mode
	if err = device.enterPTX(); err != nil {
		return
	}

	// PULL CE_PIN HIGH to enter PTX mode
	bcm2835.GpioSet(CE_PIN)

	// Wait for TSTBY_2_A to enter TX mode
	time.Sleep(TSTBY_2_A)

	// Write new register back to the device
	write_buf := make([]byte, 2)
	read_buf := make([]byte, 2) // Only need 1 byte for the status register
	write_buf[0] = FLUSH_TX

	if err = device.spiTx(write_buf, read_buf); err != nil {
		return
	}

	// Return to Standby-I mode
	if err = device.standbyI(); err != nil {
		return
	}

	return
}

func (device *NRF24L01P) resetTxDs() (err error) {

	device.Info("Resetting Data Sent TX FIFO interrupt")
	err = device.writeRegister(STATUS, TX_DS, TX_DS)

	return
}

func (device *NRF24L01P) resetMaxRetries() (err error) {

	device.Info("Resetting Maximum number of TX retransmits interrupt")
	err = device.writeRegister(STATUS, MAX_RT, MAX_RT)

	return
}

func (device *NRF24L01P) Send(data []byte) (err error) {

	// Is the TC_FIFO full?
	status, err := device.readStatusRegister()
	if err != nil {
		return
	} else {
		if status.TX_FULL == 1 {
			return fmt.Errorf("TX Output buffer is full!")
		}
	}

	// Enter PTX mode
	if err = device.enterPTX(); err != nil {
		return
	}

	// Write data buffer to TX_FIFO
	read_buf := make([]byte, 1)
	write_buf := make([]byte, 1)
	write_buf[0] = W_TX_PAYLOAD
	for _, buf := range data {
		write_buf = append(write_buf, buf)
	}

	device.Debug("Writing to TX_FIFO: %08b\n", write_buf)
	if err = device.spiTx(write_buf, read_buf); err != nil {
		return
	}

	// PULL CE_PIN HIGH for THCE to transmit
	device.Info("Sending message: %s\n", data)
	bcm2835.GpioSet(CE_PIN)
	time.Sleep(THCE)
	bcm2835.GpioClr(CE_PIN)

	// Wait for TSTBY_2_A to enter TX mode
	time.Sleep(TSTBY_2_A)

	// Wait for ACK
	ack_timeout, ard, err := device.txAckTimeout()
	if err != nil {
		return
	}

	ack_chan := make(chan bool, 1)
	err_chan := make(chan error, 1)
	go func() {
		for {
			status_register, err := device.readStatusRegister()
			if err != nil {
				err_chan <- err
			} else {
				if status_register.TX_DS == 1 {
					ack_chan <- true
					break
				}
			}
			time.Sleep(ard)
		}

	}()

	select {
	case <-ack_chan:
		device.Info("Acknowledgement received")

		// Flush RX_FIFO
		if err = device.flushRxFifo(); err != nil {
			return
		}

		// Reset Data Sent TX FIFO interrupt
		if err = device.resetTxDs(); err != nil {
			return
		}

	case err = <-err_chan:
		return
	case <-time.After(ack_timeout):

		// Flush TX_FIFO
		if err = device.flushTxFifo(); err != nil {
			device.Error(err.Error())
		}

		// Reset max retries
		if err = device.resetMaxRetries(); err != nil {
			device.Error(err.Error())
		}

		// Reset Data Sent TX FIFO interrupt
		if err = device.resetTxDs(); err != nil {
			return
		}

		return fmt.Errorf("No ACK received in %s", ack_timeout)
	}

	device.Info("Sent successfully")

	return
}

func (device *NRF24L01P) txAckTimeout() (timeout time.Duration, ard time.Duration, err error) {

	setup_retr, err := device.readSetupRetrRegister()
	if err != nil {
		return
	}

	if setup_retr.ARD == 0 {
		ard = time.Duration(250 * time.Microsecond)
	} else {
		ard = time.Duration(time.Duration(setup_retr.ARD*250) * time.Microsecond)
	}
	timeout = ard * time.Duration(uint16(setup_retr.ARC))

	return
}
