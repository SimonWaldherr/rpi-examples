package main

import (
	"flag"
	"fmt"

	"github.com/SimonWaldherr/hx711go"
)

var TargetWeight int
var AdjustZero int
var AdjustScale float64

func main() {
	flag.IntVar(&TargetWeight, "target", 100, "weight to be measured")
	flag.IntVar(&AdjustZero, "zero", -94932, "adjust zero value")
	flag.Float64Var(&AdjustScale, "scale", 62.8, "adjust scale value")
	flag.Parse()

	err := hx711.HostInit()
	if err != nil {
		fmt.Println("HostInit error:", err)
		return
	}

	hx711, err := hx711.NewHx711("6", "5")

	defer hx711.Shutdown()

	hx711.AdjustZero = AdjustZero
	hx711.AdjustScale = AdjustScale

	data, err := hx711.ReadDataMedian(3)
	if err != nil {
		fmt.Println("ReadDataRaw error:", err)
		return
	}
	fmt.Printf("measurement: %v\n", data)
}
