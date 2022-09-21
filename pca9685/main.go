package main

import (
	"os"
	"time"

	"github.com/op/go-logging"
	"github.com/sergiorb/pca9685-golang/device"
	"golang.org/x/exp/io/i2c"
)

const (
	I2C_ADDR      = "/dev/i2c-1"
	ADDR_01       = 0x40
	SERVO_CHANNEL = 4
	MIN_PULSE     = 150
	MAX_PULSE     = 650
)

func init() {

	stderrorLog := logging.NewLogBackend(os.Stderr, "", 0)

	stderrorLogLeveled := logging.AddModuleLevel(stderrorLog)
	stderrorLogLeveled.SetLevel(logging.ERROR, "")

	logging.SetBackend(stderrorLogLeveled)
}

func setPercentage(p *device.Pwm, percent float32) {
	pulseLength := int((MAX_PULSE-MIN_PULSE)*percent/100 + MIN_PULSE)
	p.SetPulse(0, pulseLength)
}

func main() {

	var mainLog = logging.MustGetLogger("PCA9685 Demo")

	i2cDevice, err := i2c.Open(&i2c.Devfs{Dev: I2C_ADDR}, ADDR_01)
	defer i2cDevice.Close()

	if err != nil {

		mainLog.Error(err)

	} else {

		var deviceLog = logging.MustGetLogger("PCA9685")

		pca9685 := device.NewPCA9685(i2cDevice, "PWM Controller", MIN_PULSE, MAX_PULSE, deviceLog)

		pca9685.Frequency = 60.0

		pca9685.Init()

		for i := 0; i < 16; i++ {
			servo := pca9685.NewPwm(i)
			setPercentage(servo, 100.0)
		}

		time.Sleep(2 * time.Second)

		for i := 0; i < 16; i++ {
			servo := pca9685.NewPwm(i)
			setPercentage(servo, 0.0)
		}

		time.Sleep(2 * time.Second)
	}
}
