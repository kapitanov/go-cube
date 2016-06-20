package cube

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

const baudRate = 57600
const delay = time.Second * 2

const cmdReportAnalog = 0xC0
const cmdReportDigital = 0xD0
const cmdSetPinMode = 0xF4
const cmdDigitalMessage = 0x90

type pinMode int

const (
	input  pinMode = 0
	output pinMode = 1
	analog pinMode = 2
	pwm    pinMode = 3
	servo  pinMode = 4
)

type firmata interface {
	PinMode(pin int, mode pinMode) error
	DigitalWrite(pin int, value bool) error
	Close() error
}

type firmataImpl struct {
	port           *serial.Port
	digitalOutData []int
}

func findSerialPort() (string, error) {
	for i := 12; i > 0; i-- {
		port := fmt.Sprintf("COM%d", i)

		serialConfig := &serial.Config{
			Name:     port,
			Baud:     baudRate,
			Parity:   serial.ParityNone,
			StopBits: serial.Stop1,
		}

		serialPort, err := serial.OpenPort(serialConfig)
		if err != nil {
			continue
		}

		serialPort.Close()
		return port, nil
	}

	return "", errors.New("Unable to auto-detect COM port")
}

const autoDetectPort = "AUTO"

func newFirmata(port string) (firmata, error) {
	if port == autoDetectPort {
		trace.Printf("Detecting Cube port...")
		autoPort, err := findSerialPort()
		if err != nil {
			trace.Printf("Unable to auto detect port")
			return nil, err
		}

		trace.Printf("Detected Cube port is %s", autoPort)
		port = autoPort
	}

	serialConfig := &serial.Config{
		Name:     port,
		Baud:     baudRate,
		Parity:   serial.ParityNone,
		StopBits: serial.Stop1,
	}

	serialPort, err := serial.OpenPort(serialConfig)
	if err != nil {
		trace.Printf("Unable to open port %s: %s", port, err)
		return nil, err
	}

	c := &firmataImpl{serialPort, make([]int, 13)}
	err = c.initialize()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *firmataImpl) initialize() error {
	time.Sleep(delay)

	cmd := make([]byte, 2)

	for i := 0; i < 6; i++ {
		cmd[0] = (byte)(cmdReportAnalog | i)
		cmd[1] = 1

		_, err := c.port.Write(cmd)
		if err != nil {
			log.Printf("Write error: %s", err)
			return err
		}

		err = c.port.Flush()
		if err != nil {
			log.Printf("Flush error: %s", err)
			return err
		}
	}

	for i := 0; i < 2; i++ {
		cmd[0] = (byte)(cmdReportDigital | i)
		cmd[1] = 1

		_, err := c.port.Write(cmd)
		if err != nil {
			trace.Printf("Write error: %s", err)
			return err
		}

		err = c.port.Flush()
		if err != nil {
			log.Printf("Flush error: %s", err)
			return err
		}
	}

	return nil
}

func (c *firmataImpl) PinMode(pin int, mode pinMode) error {
	cmd := make([]byte, 3)
	cmd[0] = cmdSetPinMode
	cmd[1] = byte(pin)
	cmd[2] = byte(mode)

	_, err := c.port.Write(cmd)
	if err != nil {
		trace.Printf("Write error: %s", err)
		return err
	}

	err = c.port.Flush()
	if err != nil {
		log.Printf("Flush error: %s", err)
		return err
	}

	return nil
}

func (c *firmataImpl) DigitalWrite(pin int, value bool) error {
	portNumber := (pin >> 3) & 0x0F

	if !value {
		c.digitalOutData[portNumber] &= ^(1 << uint(pin&0x07))
	} else {
		c.digitalOutData[portNumber] |= (1 << uint(pin&0x07))
	}

	cmd := make([]byte, 3)
	cmd[0] = byte(cmdDigitalMessage | portNumber)
	cmd[1] = byte(c.digitalOutData[portNumber] & 0x7F)
	cmd[2] = (byte)(c.digitalOutData[portNumber] >> 7)

	_, err := c.port.Write(cmd)
	if err != nil {
		trace.Printf("Write error: %s", err)
		return err
	}

	err = c.port.Flush()
	if err != nil {
		log.Printf("Flush error: %s", err)
		return err
	}

	return nil
}

func (c *firmataImpl) Close() error {
	err := c.port.Close()
	if err != nil {
		log.Printf("Close error: %s", err)
		return err
	}

	return nil
}
