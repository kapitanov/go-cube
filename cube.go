package cube

import (
	"log"
	"os"
	"time"
)

var trace = log.New(os.Stdout, "[go-cube] ", log.Ltime)

var greenPins = [...]int{6, 10}
var redPins = [...]int{5, 9}

type BlinkMode int

const (
	BlinkFast BlinkMode = 0x01
	BlinkSlow BlinkMode = 0x02

	BlinkRed   BlinkMode = 0x10
	BlinkGreen BlinkMode = 0x20
)

type Cube interface {
	Off() error
	Red() error
	Green() error
	Blink(mode BlinkMode) error

	Close()
}

type cube struct {
	fm firmata
}

const AutoDetectPort = autoDetectPort

func NewCube(port string) (Cube, error) {
	trace.Printf("Connecting to Cube...", port)
	fm, err := newFirmata(port)
	if err != nil {
		return nil, err
	}

	c := cube{fm}

	err = c.initiaize()
	if err != nil {
		return nil, err
	}

	err = c.Off()
	if err != nil {
		return nil, err
	}

	trace.Printf("Connected to Cube at %s", port)
	return c, nil
}

func (c cube) Close() {
	c.off()
	c.fm.Close()
}

func (c cube) Off() error {
	trace.Printf("Cube is OFF")
	return c.off()
}

func (c cube) Red() error {
	trace.Printf("Cube is RED")
	return c.red()
}

func (c cube) Green() error {
	trace.Printf("Cube is GREEN")
	return c.green()
}

const fastDuration = 120 * time.Millisecond
const slowDuration = 250 * time.Millisecond

func (c cube) Blink(mode BlinkMode) error {
	duration := fastDuration
	name := "BLINKING "
	enableRed := false
	enableGreen := false

	if (mode & BlinkFast) == BlinkFast {
		duration = fastDuration
		name += "FAST "
	} else if (mode & BlinkSlow) == BlinkSlow {
		duration = slowDuration
		name += "SLOW "
	}

	if (mode & BlinkRed) == BlinkRed {
		enableRed = true
	}
	if (mode & BlinkGreen) == BlinkGreen {
		enableGreen = true
	}

	if enableRed && enableGreen {
		name += "RED/GREEN"
	} else if enableRed {
		name += "RED"
	} else if enableGreen {
		name += "GREEN"
	}

	trace.Printf("Cube is %s", name)
	c.off()

	for i := 0; i < 5; i++ {
		if enableRed {
			err := c.red()
			if err != nil {
				return err
			}

			time.Sleep(duration)
		}

		if enableGreen {
			err := c.green()
			if err != nil {
				return err
			}
			time.Sleep(duration)
		}

		if !enableRed || !enableGreen {
			err := c.off()
			if err != nil {
				return err
			}
			time.Sleep(duration)
		}
	}

	err := c.off()
	if err != nil {
		return err
	}

	trace.Printf("Cube is OFF")
	return nil
}

func (c cube) initiaize() error {
	for _, pin := range redPins {
		err := c.fm.PinMode(pin, output)
		if err != nil {
			return err
		}
	}

	for _, pin := range greenPins {
		err := c.fm.PinMode(pin, output)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c cube) off() error {
	return c.write(false, false)
}

func (c cube) red() error {
	return c.write(false, true)
}

func (c cube) green() error {
	return c.write(true, false)
}

func (c cube) write(green, red bool) error {
	for _, pin := range redPins {
		err := c.fm.DigitalWrite(pin, red)
		if err != nil {
			return err
		}
	}

	for _, pin := range greenPins {
		err := c.fm.DigitalWrite(pin, green)
		if err != nil {
			return err
		}
	}

	return nil
}
