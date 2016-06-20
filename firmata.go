package cube

import (
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
}

type firmataImpl struct {
	port           *serial.Port
	digitalOutData []int
}

func newFirmata(port string) (firmata, error) {
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

	trace.Printf("Port %s is open", port)
	c := &firmataImpl{serialPort, make([]int, 13)}
	err = c.initialize()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (mode pinMode) String() string {
	switch mode {
	case input:
		return "INPUT"
	case output:
		return "OUTPUT"
	case analog:
		return "ANALOG"
	case pwm:
		return "PWM"
	case servo:
		return "SERVO"
	default:
		return "?"
	}
}

func (c *firmataImpl) initialize() error {
	time.Sleep(delay)

	cmd := make([]byte, 2)

	for i := 0; i < 6; i++ {
		cmd[0] = (byte)(cmdReportAnalog | i)
		cmd[1] = 1

		trace.Printf("report analog %d", i)
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

		trace.Printf("report digital %d", i)
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

	trace.Printf("PinMode %d %s", pin, mode)
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

	if value {
		c.digitalOutData[portNumber] &= ^(1 << uint(pin&0x07))
	} else {
		c.digitalOutData[portNumber] |= (1 << uint(pin&0x07))
	}

	cmd := make([]byte, 3)
	cmd[0] = byte(cmdDigitalMessage | portNumber)
	cmd[1] = byte(c.digitalOutData[portNumber] & 0x7F)
	cmd[2] = (byte)(c.digitalOutData[portNumber] >> 7)

	if value {
		trace.Printf("DigitalWrite %d 1", pin)
	} else {
		trace.Printf("DigitalWrite %d 0", pin)
	}

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
