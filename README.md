# Raspberry Pi Golang Examples

[![Go Report Card](https://goreportcard.com/badge/github.com/simonwaldherr/rpi-examples)](https://goreportcard.com/report/github.com/simonwaldherr/rpi-examples)  
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)  

If you liked this project, you may also like one of these projects [golang-examples](https://github.com/SimonWaldherr/golang-examples), [golang-benchmarks](https://github.com/SimonWaldherr/golang-benchmarks) or [sql-examples](https://github.com/SimonWaldherr/sql-examples).

## About

These examples explain how to use different [Raspberry Pi](https://www.raspberrypi.org) accessories and hardware with [Golang](https://golang.org). There will be more examples from time to time.

if you like, feel free to add more examples. Many thanks to all [contributors](https://github.com/SimonWaldherr/rpi-examples/graphs/contributors).

## Install go(lang)

with [homebrew](http://mxcl.github.io/homebrew/):

```Shell
sudo brew install go
```

with [apt](http://packages.qa.debian.org/a/apt.html)-get:

```Shell
sudo apt-get install golang
```

[install Golang manually](https://golang.org/doc/install)
or
[compile it yourself](https://golang.org/doc/install/source)

## Get a Pi

of course you also need a Raspberry Pi, normaly you can get one at [Amazon](https://amzn.to/3xDegoT), [BerryBase](https://www.berrybase.de/raspberry-pi/) or [Reichelt](https://www.reichelt.de/raspberry-pi-compute-modul-4-8gb-ram-8gb-emmc-wlan-rpi-cm4w-8gb8gb-p290550.html?&nbc=1).  
Sometimes Raspberrys are unfortunately difficult to obtain, but there is the great website [rpilocator](https://rpilocator.com), which shows very clearly which dealers have the various Raspberry models on stock. 

## Examples

### [HX711](https://github.com/SimonWaldherr/rpi-examples/tree/master/hx711) 
The hx711 is a chip that makes it possible to query load cells with the RaspberryPi (or other systems, e.g. the Arduino). 
You can [buy boards with the hx711-chip on Amazon](https://amzn.to/3LyGWFl). 
There are also complete [sets with a hx711 board and a load cell](https://amzn.to/3xHaFWY). 

### [PCA9685](https://github.com/SimonWaldherr/rpi-examples/tree/master/pca9685) 
The pca9685 is a PWM driver with 12-bit resolution (4096 steps) for up to 16 separately controllable devices with an operating voltage of up to 6V. This makes it possible to control up to 16 PWM outputs with just two pins on the RaspberryPi. 
The pca9685 is controlled via I2C, which means that several pca9685 can be connected in a row and with up to 62 selectable addresses, up to 992 PWM outputs with 2 pins can be controlled. 
You can [buy a great board with the pca9685-chip on Amazon](https://amzn.to/3DGVCAm). 

### [PCF8574](https://github.com/SimonWaldherr/rpi-examples/tree/master/pcf8574) 
The PCF8574 is an 8-bit I/O port expander connected via the I2C bus. Anyone who has ever suffered from "lack of pins" in one of their applications knows what is meant. Here, too, only two pins are required to control 8 pins (per board). 
There are some [pcf8574 boards available on Amazon](https://amzn.to/3R7sTaV).

### [WS2812](https://github.com/SimonWaldherr/rpi-examples/tree/master/ws2812) 
The ws2812 is an "intelligent" LED, the chip not only contains 3 LEDs (in the colors red, green and blue), but also an IC which enables the control of the LEDs. The LEDs can be controlled in brightness and combination. The ws2812 light chains are available in a wide variety of variants, they differ in the distance between the LEDs, there are waterproof light chains, different colors of the circuit board, ... 
For example, Amazon has this [144 "pixel" per meter ws2812 light chain with white PCB](https://amzn.to/3Sk0Hmm).

