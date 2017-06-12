package nrf24

import (
	"sync"

	"github.com/davent/bcm2835"
	"golang.org/x/exp/io/spi"
)

type NRF24L01P struct {
	cond *sync.Cond
	spi  *spi.Device

	LogLevel uint8
}

func init() {

	// Set the CE_PIN to Output
	err := bcm2835.Init()
	if err != nil {
		panic(err)
	}
	bcm2835.GpioFsel(CE_PIN, bcm2835.Output)

	return
}

func New(spi_device string) (device *NRF24L01P, err error) {

	device = &NRF24L01P{
		cond:     sync.NewCond(new(sync.Mutex)),
		LogLevel: ERROR,
	}

	// Open the SPI device
	if device.spi, err = spi.Open(&spi.Devfs{
		Dev:      spi_device,
		Mode:     SPI_MODE,
		MaxSpeed: SPI_SPEED,
	}); err != nil {
		return
	}

	// Enable Dynamic Payload Length
	if err = device.enableDPL(); err != nil {
		return
	}

	// Perform a NOP register lookup to get status
	if _, err = device.readRegister(NOP); err != nil {
		return
	}

	return
}

func (device *NRF24L01P) Close() (err error) {

	device.Info("Closing device")

	// Shutdown the device gracefully
	device.spi.Close()

	return
}

func (device *NRF24L01P) spiTx(w, r []byte) (err error) {

	// Acquire lock
	device.cond.L.Lock()

	// Write to the SPI device
	err = device.spi.Tx(w, r)

	// Release lock
	device.cond.L.Unlock()

	return
}

func (device *NRF24L01P) readRegister(register byte) (read_buf []byte, err error) {

	// Specify size of register in bytes
	var register_size uint8
	switch register {
	case NOP:
		register_size = 0
	case RX_ADDR_P0, RX_ADDR_P1, TX_ADDR:
		register_size = 5
	default:
		register_size = 1
	}

	// Create a buffer to send []byte to SPI device
	// (+1 to accomadate the SPI command word at byte 0)
	write_buf := make([]byte, register_size+1)

	// Create a buffer to receive []byte from the SPI device
	// (+1 to accomadate the STATUS register returned at byte 0)
	read_buf = make([]byte, register_size+1)

	// Set the byte of the register we wish to read at byte 0
	write_buf[0] = register

	if err = device.spiTx(write_buf, read_buf); err != nil {
		return
	}

	device.Debug("Reading register %X: %X %08b\n", register, read_buf, read_buf)

	// Return the response buffer minus the 1st byte (STATUS register)
	return read_buf[1:], nil
}

func (device *NRF24L01P) writeRegister(register byte, write_mask byte, write_value byte) (err error) {

	// Get the current register values
	var buf []byte
	buf, err = device.readRegister(register)
	if err != nil {
		return
	}

	// Amend register with new value
	buf[0] = buf[0] ^ ((buf[0] ^ write_value) & write_mask)

	// Write new register back to the device
	read_buf := make([]byte, 1) // Only need 1 byte for the status register
	write_buf := make([]byte, 1)
	write_buf[0] = (register | W_REGISTER_MASK)
	for _, n := range buf {
		write_buf = append(write_buf, n)
	}

	device.Debug("Writing register %X: %X %08b\n", register, write_buf, write_buf)

	// Write the buffer to the SPI device
	if err = device.spiTx(write_buf, read_buf); err != nil {
		return
	}

	// If all went well, we can return a nil err
	return
}
