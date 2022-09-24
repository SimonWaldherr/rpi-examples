package main

import (
	"fmt"
	"time"

	"golang.org/x/exp/io/i2c"
	"simonwaldherr.de/go/golibs/bitmask"
)

const (
	I2C_ADDR = "/dev/i2c-1"
)

var pins map[int]int
var i2cDev1, i2cDev2 *i2c.Device
var bm1, bm2 *bitmask.Bitmask

func setValve(valve int, status bool) {
	var pin int
	pin = pins[valve]
	
	if pin > 7 {
		pin = pin-6
		bm2.Set(pin, !status)
		i2cDev2.Write([]byte{byte(bm2.Int())})
		return
	}
	
	bm1.Set(pin, !status)
	i2cDev1.Write([]byte{byte(bm1.Int())})
}

func setPump(status bool) {
	bm2.Set(0, !status)
	i2cDev2.Write([]byte{byte(bm2.Int())})
}

func setMasterValve(status bool) {
	bm2.Set(1, !status)
	i2cDev2.Write([]byte{byte(bm2.Int())})
}

func init() {
	pins = map[int]int{
		1:  0,
		2:  1,
		3:  2,
		4:  3,
		5:  4,
		6:  5,
		7:  8,
		8:  9,
		9:  10,
		10: 11,
		11: 12,
		12: 13,
	}
	
	var err error
	i2cDev1, err = i2c.Open(&i2c.Devfs{Dev: I2C_ADDR}, 0x20)
	if err != nil {
		panic(err)
	}
	
	
	i2cDev2, err = i2c.Open(&i2c.Devfs{Dev: I2C_ADDR}, 0x21)
	if err != nil {
		panic(err)
	}
	
	bm1 = bitmask.New(0b11111111)
	bm2 = bitmask.New(0b11111111)
	
	i2cDev1.Write([]byte{byte(bm1.Int())})
	i2cDev2.Write([]byte{byte(bm2.Int())})
	
	time.Sleep(10 * time.Millisecond)
}

func main() {
	defer i2cDev1.Close()
	defer i2cDev2.Close()
	
	for {
		for i := 1; i < 13; i++ {
			setValve(i, true)
			fmt.Printf("set Valve %d on, bitmask is now: %b,%b\n", i, []byte{byte(bm1.Int())}, []byte{byte(bm2.Int())})
			time.Sleep(1 * time.Second)
			
			setValve(i, false)
			fmt.Printf("set Valve %d off, bitmask is now: %b,%b\n", i, []byte{byte(bm1.Int())}, []byte{byte(bm2.Int())})
			time.Sleep(350 * time.Millisecond)
			
			fmt.Println()
		}
		fmt.Println()
	}
}
