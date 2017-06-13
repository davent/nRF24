package nrf24

import (
	"time"

	"github.com/davent/bcm2835"
	"golang.org/x/exp/io/spi"
)

const (
	VERSION     string = "0.0.1"
	DEVICE_NAME string = "nRF24L01+"

	/*
		SPI config
	*/
	SPI_MODE  spi.Mode = spi.Mode0
	SPI_SPEED int64    = 8000000
	CE_PIN             = bcm2835.RPI_V2_GPIO_P1_22

	/*
		Timings
	*/
	TPD_2_STBY time.Duration = 5 * time.Millisecond
	TSTBY_2_A  time.Duration = 130 * time.Microsecond
	THCE       time.Duration = 10 * time.Microsecond

	/*
	  SPI commands
	*/
	R_REGISTER_MASK    byte = 0x00 // Read command and status registers.
	W_REGISTER_MASK    byte = 0x20 // Write command and status registers.
	R_RX_PAYLOAD       byte = 0x61 // Read the RX_FIFO
	W_TX_PAYLOAD       byte = 0xA0 // Write to the TX_FIFO
	FLUSH_TX           byte = 0xE1 // Flush the TX_FIFO
	FLUSH_RX           byte = 0xE2 // Flush the RX_FIFO
	REUSE_TX_PL        byte = 0xE3
	R_RX_PL_WID        byte = 0x60
	W_ACK_PAYLOAD_MASK byte = 0xA8
	W_TX_PAYLOAD_NOACK byte = 0xB0
	NOP                byte = 0xFF

	/*
		Command and status registers
	*/
	CONFIG      byte = 0x00 // Configuration Register
	EN_AA       byte = 0x01 // Enhanced ShockBurstTM - Enable 'Auto Acknowledgment' Function Disable this functionality to be compatible with nRF2401
	EN_RXADDR   byte = 0x02 // Enabled RX Addresses
	SETUP_AW    byte = 0x03 // Setup of Address Widths (common for all data pipes)
	SETUP_RETR  byte = 0x04 // Setup of Automatic Retransmission
	RF_CH       byte = 0x05 // RF Channel
	RF_SETUP    byte = 0x06 // RF Setup Register
	STATUS      byte = 0x07 // Status Register (In parallel to the SPI command word applied on the MOSI pin, the STATUS register is shifted serially out on the MISO pin)
	OBSERVE_TX  byte = 0x08 // Transmit observe register
	RPD         byte = 0x09 // Received Power Detector.
	RX_ADDR_P0  byte = 0x0A // Receive address data pipe 0. 5 Bytes maximum length. (LSByte is written first. Write the number of bytes defined by SETUP_AW)
	RX_ADDR_P1  byte = 0x0B // Receive address data pipe 1. 5 Bytes maximum length. (LSByte is written first. Write the number of bytes defined by SETUP_AW)
	RX_ADDR_P2  byte = 0x0C // Receive address data pipe 2. Only LSB. MSBytes are equal to RX_ADDR_P1[39:8]
	RX_ADDR_P3  byte = 0x0D // Receive address data pipe 3. Only LSB. MSBytes are equal to RX_ADDR_P1[39:8]
	RX_ADDR_P4  byte = 0x0E // Receive address data pipe 4. Only LSB. MSBytes are equal to RX_ADDR_P1[39:8]
	RX_ADDR_P5  byte = 0x0F // Receive address data pipe 5. Only LSB. MSBytes are equal to RX_ADDR_P1[39:8]
	TX_ADDR     byte = 0x10 // Transmit address. Used for a PTX device only. (LSByte is written first). Set RX_ADDR_P0 equal to this address to handle automatic acknowledge if this is a PTX device with Enhanced ShockBurstTM
	RX_PW_P0    byte = 0x11 // Number of bytes in RX payload in data pipe 0 (1 to 32 bytes).
	RX_PW_P1    byte = 0x12 // Number of bytes in RX payload in data pipe 1 (1 to 32 bytes).
	RX_PW_P2    byte = 0x13 // Number of bytes in RX payload in data pipe 2 (1 to 32 bytes).
	RX_PW_P3    byte = 0x14 // Number of bytes in RX payload in data pipe 3 (1 to 32 bytes).
	RX_PW_P4    byte = 0x15 // Number of bytes in RX payload in data pipe 4 (1 to 32 bytes).
	RX_PW_P5    byte = 0x16 // Number of bytes in RX payload in data pipe 5 (1 to 32 bytes).
	FIFO_STATUS byte = 0x17 // FIFO Status Register
	DYNPD       byte = 0x1C // Enable dynamic payload length
	FEATURE     byte = 0x1D // Feature Register

	/*
		Register bit masks
	*/

	// CONFIG
	MASK_RX_DR  byte = 0x40
	MASK_TX_DS  byte = 0x20
	MASK_MAX_RT byte = 0x10
	EN_CRC      byte = 0x08
	CRCO        byte = 0x04
	PWR_UP      byte = 0x02
	PRIM_RX     byte = 0x01

	CRC0_8BIT  byte = 0x00
	CRC0_16BIT byte = 0x04

	// SETUP_RETR
	ARD byte = 0xF0
	ARC byte = 0x0F

	// RF_CH
	RF_CH_MASK byte = 0x7F
	RF_CH_BASE uint = 2400

	// RF_SETUP
	CONT_WAVE  byte = 0x80
	RF_DR_LOW  byte = 0x20
	PLL_LOCK   byte = 0x10
	RF_DR_HIGH byte = 0x08
	RF_PWR     byte = 0x06

	RF_DR_MASK    byte = 0x28
	RF_DR_1MBPS   byte = 0x00
	RF_DR_2MBPS   byte = 0x08
	RF_DR_250KBPS byte = 0x20

	RF_PWR_MASK  byte = 0x06
	RF_PWR_18DBM byte = 0x00
	RF_PWR_12DBM byte = 0x02
	RF_PWR_6DBM  byte = 0x04
	RF_PWR_0DBM  byte = 0x06

	// STATUS
	RX_DR   byte = 0x40
	TX_DS   byte = 0x20
	MAX_RT  byte = 0x10
	RX_P_NO byte = 0x0E
	TX_FULL byte = 0x01

	//RX_PW_P{n}
	RX_PW_MASK byte = 0x3F

	// FIFO_STATUS
	TX_REUSE     byte = 0x40
	TX_FULL_FIFO byte = 0x20
	TX_EMPTY     byte = 0x10
	RX_FULL      byte = 0x02
	RX_EMPTY     byte = 0x01

	// DYNPD
	DPL_P5 byte = 0x20
	DPL_P4 byte = 0x10
	DPL_P3 byte = 0x08
	DPL_P2 byte = 0x04
	DPL_P1 byte = 0x02
	DPL_P0 byte = 0x01

	// FEATURE
	EN_DPL     byte = 0x04
	EN_ACK_PAY byte = 0x02
	EN_SYN_ACK byte = 0x01
)
