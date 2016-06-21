# go-cube
Client API for [Amperka Cube](http://wiki.amperka.ru/device:techno-cube) in Golang

[![GoDoc](https://godoc.org/github.com/kapitanov/go-cube?status.svg)](https://godoc.org/github.com/kapitanov/go-cube)

## API
Please have a look at the [GoDoc documentation](https://godoc.org/github.com/kapitanov/go-cube) for a detailed API description.

## Example

```go
package main

import (
	"fmt"
	"github.com/kapitanov/go-cube"
)

func main() {
  c, err := cube.NewCube(cube.AutoDetectPort)
	if err != nil {
		panic(err)
	}
	
	fmt.Print("Cube is RED\n")
	c.Red()
	
	fmt.Print("Cube is GREEN\n")
	c.Green()
	
	fmt.Print("Cube is BLINKING FAST IN RED AND GREEN\n")
	c.Blink(cube.BlinkFast | cube.BlinkRed | cube.BlinkGreen)
	
	fmt.Print("Cube is OFF\n")
	c.Off()
}
```
