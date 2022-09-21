package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/SimonWaldherr/hx711go"
	"simonwaldherr.de/go/golibs/gcurses"
	"simonwaldherr.de/go/golibs/xmath"
)

var TargetWeight int
var AdjustZero int
var AdjustScale float64

func scaleDelay(scaleDelta int, timeout time.Duration) {
	runtime.GC()
	hx711, err := hx711.NewHx711("6", "5")

	if err != nil {
		fmt.Println("NewHx711 error:", err)
		return
	}

	defer hx711.Shutdown()

	for {
		err = hx711.Reset()
		if err == nil {
			break
		}
		log.Print("hx711 BackgroundReadMovingAvgs Reset error:", err)
		time.Sleep(time.Second)
	}

	hx711.AdjustZero = AdjustZero
	hx711.AdjustScale = AdjustScale

	c1 := make(chan bool, 1)
	go func() {
		var tara = []float64{}
		var i int
		var data, predata float64

		fmt.Println("Tara")

		for i = 0; i < 5; i++ {
			time.Sleep(600 * time.Microsecond)

			data, err := hx711.ReadDataMedian(3)
			if err != nil {
				fmt.Println("ReadDataRaw error:", err)
				continue
			}

			tara = append(tara, data)
		}
		taraAvg := float64(xmath.Round(xmath.Arithmetic(tara)))

		fmt.Printf("New tara set to: %v\n", taraAvg)

		writer := gcurses.New()
		writer.Start()

		for {
			time.Sleep(10 * time.Millisecond)
			data2, err := hx711.ReadDataRaw()

			if err != nil {
				if fmt.Sprintln("ReadDataRaw error:", err) != "ReadDataRaw error: waitForDataReady error: timeout\n" {
					fmt.Println("ReadDataRaw error:", err)
				}
				continue
			}

			predata = data
			data = (float64(data2-hx711.AdjustZero) / hx711.AdjustScale) - taraAvg

			fmt.Fprintf(writer, "scale value: %d\n", xmath.Round((data+predata)/2))
			if int(data) > scaleDelta && int(predata) > scaleDelta {
				writer.Stop()
				fmt.Printf("set weight reached. weight is: %d\n", xmath.Round(data))
				c1 <- true
				return
			}
		}
	}()

	select {
	case _ = <-c1:
		return
	case <-time.After(timeout):
		fmt.Println("timeout")
		return
	}
}

func main() {
	runtime.GOMAXPROCS(3)
	flag.IntVar(&TargetWeight, "target", 100, "weight to be measured")
	flag.IntVar(&AdjustZero, "zero", -94932, "adjust zero value")
	flag.Float64Var(&AdjustScale, "scale", 62.8, "adjust scale value")
	flag.Parse()

	err := hx711.HostInit()
	if err != nil {
		fmt.Println("HostInit error:", err)
		return
	}

	fmt.Printf("measurement target set to %d\n", TargetWeight)

	scaleDelay(TargetWeight, 5*time.Minute)

	fmt.Println("measurement completed")
}
