package nrf24

import (
	"fmt"
	"log"
)

const (
	CRITICAL uint8 = 1
	ERROR    uint8 = 2
	WARNING  uint8 = 3
	INFO     uint8 = 4
	DEBUG    uint8 = 5
)

func (device *NRF24L01P) Critical(format string, v ...interface{}) {
	if device.LogLevel >= CRITICAL {
		printf("Critical", format, v...)
	}
}

func (device *NRF24L01P) Error(format string, v ...interface{}) {
	if device.LogLevel >= ERROR {
		printf("Error", format, v...)
	}
}

func (device *NRF24L01P) Warning(format string, v ...interface{}) {
	if device.LogLevel >= WARNING {
		printf("Warning", format, v...)
	}
}

func (device *NRF24L01P) Info(format string, v ...interface{}) {
	if device.LogLevel >= INFO {
		printf("Info", format, v...)
	}
}

func (device *NRF24L01P) Debug(format string, v ...interface{}) {
	if device.LogLevel >= DEBUG {
		printf("Debug", format, v...)
	}
}

func printf(level, format string, v ...interface{}) {
	log.Printf(fmt.Sprintf("[%s/%s] %s", DEVICE_NAME, level, format), v...)
}
