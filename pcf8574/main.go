package main

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/exp/io/i2c"
	"simonwaldherr.de/go/golibs/bitmask"
)

const (
	I2C_ADDR = "/dev/i2c-1"
	ADDR_01  = 0x20
)

func main() {
	i2cDevice, err := i2c.Open(&i2c.Devfs{Dev: I2C_ADDR}, ADDR_01)
	if err != nil {
		panic(err)
	}
	defer i2cDevice.Close()

	bitmask := bitmask.New(0b11111111)

	i2cDevice.Write([]byte(strconv.Itoa(bitmask.Int())))
	time.Sleep(1500 * time.Millisecond)

	for {
		for i := 0; i < 8; i++ {
			bitmask.Set(i, false)
			i2cDevice.Write([]byte{byte(bitmask.Int())})
		}
		time.Sleep(1500 * time.Millisecond)
		for i := 0; i < 8; i++ {
			bitmask.Set(i, true)
			i2cDevice.Write([]byte{byte(bitmask.Int())})
		}
		time.Sleep(1500 * time.Millisecond)
	}

	for i := 0; i < 8; i++ {
		bitmask.Set(i, false)
		i2cDevice.Write([]byte{byte(bitmask.Int())})
		fmt.Println(bitmask)
		fmt.Printf("%b\n", []byte{byte(bitmask.Int())})
		bitmask.Set(i, true)
		time.Sleep(8000 * time.Millisecond)
	}
}
