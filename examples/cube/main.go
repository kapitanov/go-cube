package main

import (
	"fmt"

	"github.com/eiannone/keyboard"

	"github.com/kapitanov/go-cube"
)

func main() {
	c, err := cube.NewCube(cube.AutoDetectPort)
	if err != nil {
		panic(err)
	}

	defer c.Close()

	fmt.Println("Press <1> key to make cube RED")
	fmt.Println("Press <2> key to make cube GREEN")
	fmt.Println("Press <2> key to turn cube OFF")

	fmt.Println("Press <4> key to make cube BLINK RED           FAST")
	fmt.Println("Press <5> key to make cube BLINK GREEN         FAST")
	fmt.Println("Press <6> key to make cube BLINK RED AND GREEN FAST")

	fmt.Println("Press <7> key to make cube BLINK RED           SLOW")
	fmt.Println("Press <8> key to make cube BLINK GREEN         SLOW")
	fmt.Println("Press <9> key to make cube BLINK RED AND GREEN SLOW")

	fmt.Println("Press <Q> key to exit")
	fmt.Println("")

	for {
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}

		switch char {
		case '1':
			c.Red()
		case '2':
			c.Green()
		case '3':
			c.Off()

		case '4':
			c.Blink(cube.BlinkSlow | cube.BlinkRed)
		case '5':
			c.Blink(cube.BlinkSlow | cube.BlinkGreen)
		case '6':
			c.Blink(cube.BlinkSlow | cube.BlinkRed | cube.BlinkGreen)

		case '7':
			c.Blink(cube.BlinkFast | cube.BlinkRed)
		case '8':
			c.Blink(cube.BlinkFast | cube.BlinkGreen)
		case '9':
			c.Blink(cube.BlinkFast | cube.BlinkRed | cube.BlinkGreen)

		case 'Q', 'q':
			c.Off()
			fmt.Println("Bye!")
			return
		}
	}
}
