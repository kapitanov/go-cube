package cube

import (
	"log"
	"os"
	"time"
)

var trace = log.New(os.Stdout, "[go-cube] ", log.Ltime)

var redPins = [...]int{6, 10}
var greenPins = [...]int{5, 9}

type Cube interface {
	Off() error
	Red() error
	Green() error
	Blink() error
}

type cube struct {
	fm firmata
}

func NewCube(port string) (Cube, error) {
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

	return c, nil
}

func (c cube) Off() error {
	return c.write(false, false)
}

func (c cube) Red() error {
	return c.write(false, true)
}

func (c cube) Green() error {
	return c.write(true, false)
}

func (c cube) Blink() error {
	c.Off()

	for i := 0; i < 5; i++ {
		err := c.Red()
		if err != nil {
			return err
		}

		time.Sleep(500 * time.Millisecond)
		err = c.Green()
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}

	err := c.Off()
	if err != nil {
		return err
	}

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
