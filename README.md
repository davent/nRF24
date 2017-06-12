# Go nRF24
Go driver for the Nordic nRF24L01+ transceiver.

This package provides functionality for interfacing with the Nordic nRF24L01+ module using Go.

It is designed primarily for use on Raspberry Pi boards but should be easily ported to any Linux based systems running as an SPI master.

### CHANGELOG
#### 0.0.1
- Very basic send/receive functionality

### TODO
- IRQ
- Reset to defaults
- Multi-channel TX/RX
- Split out logging to seperate package
- Pretty much everything else. Feel free to submit feature requests!

### Usage
The main goal of this package is to make using the nRF24L01+ module as easy as possible. It should be simple to send a message from one tranceiver to another using default configuration options.

#### Example message sending
```
package main

import (
  "github.com/davent/nRF24"
)

func main() {

  radio, err := nrf24.New("/dev/spidev0.0")
  if err != nil {
    panic(err)
  }
  defer radio.Close()

  // Turn the power to the radio on
  radio.PowerOn(true)

  // Send our "Hello world!" messages
  if err = radio.Send([]byte("Hello world!")); err != nil {
    panic(err)
  }

  // Don't forget to turn the radio off
  radio.PowerOn(false)

}
```

#### Example message receiving
```
package main

import (
  "fmt"
  "github.com/davent/nRF24"
)

func main() {

  radio, err := nrf24.New("/dev/spidev0.0")
  if err != nil {
    panic(err)
  }
  defer radio.Close()

  // Turn the power to the radio on
  radio.PowerOn(true)

  // Wait for messages to be received
  incoming, err := radio.Receive()
  if err != nil {
    panic(err)
  }

  message := <-incoming
  fmt.Printf("Message received: %s\n", message)

  // Don't forget to turn the radio off
  radio.PowerOn(false)

}
```
