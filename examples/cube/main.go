package main

import (
	"flag"

	"github.com/kapitanov/go-cube"
)

var port = "COM1"
var command = ""

func init() {
	flag.StringVar(&port, "p", "COM1", "Port number")
	flag.StringVar(&command, "c", "", "Command")
}

func main() {
	flag.Parse()

	cube, err := cube.NewCube(port)
	if err != nil {
		panic(err)
	}

	switch command {
	case "off":
		cube.Off()
		break
	case "red":
		cube.Red()
		break
	case "green":
		cube.Green()
		break
	case "blink":
		cube.Blink()
		break

	default:
		break
	}
}
