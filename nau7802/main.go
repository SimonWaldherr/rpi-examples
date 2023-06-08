package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/exp/io/i2c"
)

const (
	DEVICE_ADDRESS = 0x2A

	// Register Map
	NAU7802_PU_CTRL     = 0x00
	NAU7802_CTRL1       = 0x01
	NAU7802_CTRL2       = 0x02
	NAU7802_OCAL1_B2    = 0x03
	NAU7802_OCAL1_B1    = 0x04
	NAU7802_OCAL1_B0    = 0x05
	NAU7802_GCAL1_B3    = 0x06
	NAU7802_GCAL1_B2    = 0x07
	NAU7802_GCAL1_B1    = 0x08
	NAU7802_GCAL1_B0    = 0x09
	NAU7802_OCAL2_B2    = 0x0A
	NAU7802_OCAL2_B1    = 0x0B
	NAU7802_OCAL2_B0    = 0x0C
	NAU7802_GCAL2_B3    = 0x0D
	NAU7802_GCAL2_B2    = 0x0E
	NAU7802_GCAL2_B1    = 0x0F
	NAU7802_GCAL2_B0    = 0x10
	NAU7802_I2C_CONTROL = 0x11
	NAU7802_ADCO_B2     = 0x12
	NAU7802_ADCO_B1     = 0x13
	NAU7802_ADCO_B0     = 0x14
	NAU7802_ADC         = 0x15 // Shared ADC and OTP 32:24
	NAU7802_OTP_B1      = 0x16 // OTP 23:16 or 7:0?
	NAU7802_OTP_B0      = 0x17 // OTP 15:8
	NAU7802_PGA         = 0x1B
	NAU7802_PGA_PWR     = 0x1C
	NAU7802_DEVICE_REV  = 0x1F

	// Bits within the PU_CTRL register
	NAU7802_PU_CTRL_RR    = 0
	NAU7802_PU_CTRL_PUD   = 1
	NAU7802_PU_CTRL_PUA   = 2
	NAU7802_PU_CTRL_PUR   = 3
	NAU7802_PU_CTRL_CS    = 4
	NAU7802_PU_CTRL_CR    = 5
	NAU7802_PU_CTRL_OSCS  = 6
	NAU7802_PU_CTRL_AVDDS = 7

	// Bits within the CTRL1 register
	NAU7802_CTRL1_GAIN     = 2
	NAU7802_CTRL1_VLDO     = 5
	NAU7802_CTRL1_DRDY_SEL = 6
	NAU7802_CTRL1_CRP      = 7

	// Bits within the CTRL2 register
	NAU7802_CTRL2_CALMOD    = 0
	NAU7802_CTRL2_CALS      = 2
	NAU7802_CTRL2_CAL_ERROR = 3
	NAU7802_CTRL2_CRS       = 4
	NAU7802_CTRL2_CHS       = 7

	// Bits within the PGA register
	NAU7802_PGA_CHP_DIS    = 0
	NAU7802_PGA_INV        = 3
	NAU7802_PGA_BYPASS_EN  = 4
	NAU7802_PGA_OUT_EN     = 5
	NAU7802_PGA_LDOMODE    = 6
	NAU7802_PGA_RD_OTP_SEL = 7

	// Bits within the PGA PWR register
	NAU7802_PGA_PWR_PGA_CURR       = 0
	NAU7802_PGA_PWR_ADC_CURR       = 2
	NAU7802_PGA_PWR_MSTR_BIAS_CURR = 4
	NAU7802_PGA_PWR_PGA_CAP_EN     = 7

	// Allowed Low drop out regulator voltages
	NAU7802_LDO_2V4 = 0b111
	NAU7802_LDO_2V7 = 0b110
	NAU7802_LDO_3V0 = 0b101
	NAU7802_LDO_3V3 = 0b100
	NAU7802_LDO_3V6 = 0b011
	NAU7802_LDO_3V9 = 0b010
	NAU7802_LDO_4V2 = 0b001
	NAU7802_LDO_4V5 = 0b000

	// Allowed gains
	NAU7802_GAIN_128 = 0b111
	NAU7802_GAIN_64  = 0b110
	NAU7802_GAIN_32  = 0b101
	NAU7802_GAIN_16  = 0b100
	NAU7802_GAIN_8   = 0b011
	NAU7802_GAIN_4   = 0b010
	NAU7802_GAIN_2   = 0b001
	NAU7802_GAIN_1   = 0b000

	// Allowed samples per second
	NAU7802_SPS_320 = 0b111
	NAU7802_SPS_80  = 0b011
	NAU7802_SPS_40  = 0b010
	NAU7802_SPS_20  = 0b001
	NAU7802_SPS_10  = 0b000

	// Select between channel values
	NAU7802_CHANNEL_1 = 0
	NAU7802_CHANNEL_2 = 1

	// Calibration state
	NAU7802_CAL_SUCCESS     = 0
	NAU7802_CAL_IN_PROGRESS = 1
	NAU7802_CAL_FAILURE     = 2
)

type NAU7802 struct {
	Dev               *i2c.Device
	zeroOffset        int32
	calibrationFactor float64
}

func NewNAU7802() (*NAU7802, error) {
	dev, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, DEVICE_ADDRESS)
	if err != nil {
		return nil, err
	}

	return &NAU7802{Dev: dev}, nil
}

func (n *NAU7802) isConnected() bool {
	data := make([]byte, 1)
	err := n.Dev.ReadReg(DEVICE_ADDRESS, data)
	if err != nil {
		return false // Sensor did not ACK
	}
	return true // All good
}

func (n *NAU7802) getBit(bit, register byte) (bool, error) {
	buf := make([]byte, 1)

	err := n.Dev.ReadReg(register, buf)
	if err != nil {
		return false, err
	}

	return (buf[0]>>bit)&1 == 1, nil
}

func (n *NAU7802) setBit(bit, register byte, value bool) error {
	buf := make([]byte, 15)
	err := n.Dev.ReadReg(register, buf)
	if err != nil {
		return err
	}

	data := buf[0]

	if value {
		data |= (1 << bit)
	} else {
		data &= ^(1 << bit)
	}

	writeBuf := []byte{data}

	return n.Dev.WriteReg(byte(register), writeBuf)
}

func (n *NAU7802) getRegister(register byte) ([]byte, error) {
	buf := make([]byte, 1)

	err := n.Dev.ReadReg(register, buf)

	if err != nil {
		fmt.Printf("Error reading register %d: %v\n", register, err)
		return []byte{}, err
	}

	return buf, nil
}

func (n *NAU7802) setRegister(register byte, value []byte) error {
	return n.Dev.WriteReg(register, value)
}

func (n *NAU7802) available() bool {
	data := make([]byte, 1)
	err := n.Dev.ReadReg(NAU7802_PU_CTRL, data)
	if err != nil {
		return false
	}
	return data[0]&(1<<NAU7802_PU_CTRL_CR) != 0
}

func (n *NAU7802) getReading() (int32, error) {
	data := make([]byte, 3)

	err := n.Dev.ReadReg(NAU7802_ADCO_B2, data)
	if err != nil {
		return 0, err // Sensor did not ACK
	}

	value := int32((uint(data[0]) << 16) | (uint(data[1]) << 8) | uint(data[2]))
	return value, nil
}

func (n *NAU7802) getAverage(average int) (int32, error) {
	var sum int32
	for i := 0; i < average; i++ {
		data, err := n.getReading()
		if err != nil {
			return 0, err
		}

		time.Sleep(1 * time.Millisecond)

		sum += data
	}

	return sum / int32(average), nil
}

func (n *NAU7802) calculateZeroOffset(average int) error {
	data, err := n.getAverage(average)
	if err != nil {
		return err
	}

	n.zeroOffset = data

	return nil
}

func (n *NAU7802) setZeroOffset(offset int32) {
	n.zeroOffset = offset
}

func (n *NAU7802) getZeroOffset() int32 {
	return n.zeroOffset
}

func (n *NAU7802) calculateCalibrationFactor(knowWeight float64, average int) error {
	data, err := n.getAverage(average)
	if err != nil {
		return err
	}

	n.calibrationFactor = float64(data-n.zeroOffset) / knowWeight

	return nil
}

func (n *NAU7802) setCalibrationFactor(factor float64) {
	n.calibrationFactor = factor
}

func (n *NAU7802) getCalibrationFactor() float64 {
	return n.calibrationFactor
}

func (n *NAU7802) getWeight(allowNegative bool, samples int) (float64, error) {
	data, err := n.getAverage(samples)
	if err != nil {
		return 0, err
	}

	if !allowNegative && data < 0 {
		return 0, errors.New("negative weight not allowed")
	}

	return float64(data-n.zeroOffset) / n.calibrationFactor, nil
}

func (n *NAU7802) setGain(gain int) error {
	if gain < 0 || gain > 7 {
		return errors.New("invalid gain value")
	}

	value, err := n.getRegister(NAU7802_CTRL1)

	if err != nil {
		return err
	}

	val := make([]byte, 1)
	val[0] = value[0]
	val[0] &= 0b11111000
	val[0] |= uint8(gain)

	return n.setRegister(NAU7802_CTRL1, val)
}

func (n *NAU7802) setLDO(ldo int) error {
	if ldo < 0 || ldo > 7 {
		return errors.New("invalid ldo value")
	}

	value, err := n.getRegister(NAU7802_CTRL1)

	if err != nil {
		return err
	}

	val := make([]byte, 1)
	val[0] = value[0]
	val[0] &= 0b11000111
	val[0] |= uint8(ldo << 3)

	n.setRegister(NAU7802_CTRL1, val)

	return n.setBit(NAU7802_PU_CTRL_AVDDS, NAU7802_PU_CTRL, true)
}

func (n *NAU7802) setSampleRate(rate int) error {
	if rate < 0 || rate > 7 {
		return errors.New("invalid sample rate value")
	}

	value, err := n.getRegister(NAU7802_CTRL2)

	if err != nil {
		return err
	}

	val := make([]byte, 1)
	val[0] = value[0]
	val[0] &= 0b10001111
	val[0] |= uint8(rate << 4)

	return n.setRegister(NAU7802_CTRL2, val)
}

func (n *NAU7802) setChannel(channel int) error {
	if channel == NAU7802_CHANNEL_1 {
		return n.setBit(NAU7802_CTRL2_CHS, NAU7802_CTRL2, false)
	} else {
		return n.setBit(NAU7802_CTRL2_CHS, NAU7802_CTRL2, true)
	}
}

func (n *NAU7802) calAFEinProgress() (bool, error) {
	if val, err := n.getBit(NAU7802_CTRL2_CALS, NAU7802_CTRL2); err != nil && val {
		return true, nil
	}

	if val, err := n.getBit(NAU7802_CTRL2_CAL_ERROR, NAU7802_CTRL2); err != nil && val {
		return false, errors.New("calibration error")
	}

	return false, nil
}

func (n *NAU7802) beginnCalibrateAFE() error {
	return n.setBit(NAU7802_CTRL2_CALS, NAU7802_CTRL2, true)
}

func (n *NAU7802) waitForCalibrateAFE(timeout time.Duration) error {
	start := time.Now()

	for {
		if time.Since(start) > timeout {
			return errors.New("timeout")
		}

		if inProgress, err := n.calAFEinProgress(); err != nil {
			return err
		} else if !inProgress {
			return nil
		}

		time.Sleep(1 * time.Millisecond)
	}
	return nil
}

func (n *NAU7802) calibrateAFE() error {
	if err := n.beginnCalibrateAFE(); err != nil {
		return err
	}

	return n.waitForCalibrateAFE(100 * time.Millisecond)
}

func (n *NAU7802) reset() error {
	n.setBit(NAU7802_PU_CTRL_RR, NAU7802_PU_CTRL, true)
	time.Sleep(1 * time.Millisecond)

	return n.setBit(NAU7802_PU_CTRL_RR, NAU7802_PU_CTRL, false)
}

func (n *NAU7802) powerUp() error {
	err := n.setBit(NAU7802_PU_CTRL_PUD, NAU7802_PU_CTRL, true)
	if err != nil {
		return err
	}
	err = n.setBit(NAU7802_PU_CTRL_PUA, NAU7802_PU_CTRL, true)
	if err != nil {
		return err
	}
	counter := 0
	for {
		data := make([]byte, 1)
		err = n.Dev.ReadReg(NAU7802_PU_CTRL, data)
		if err != nil {
			return err
		}
		if data[0]&(1<<NAU7802_PU_CTRL_PUR) != 0 {
			break
		}
		time.Sleep(1 * time.Millisecond)
		counter++
		if counter > 100 {
			return fmt.Errorf("PowerUp failed")
		}
	}
	return nil
}

func (n *NAU7802) powerDown() error {
	n.setBit(NAU7802_PU_CTRL_PUD, NAU7802_PU_CTRL, false)
	n.setBit(NAU7802_PU_CTRL_PUA, NAU7802_PU_CTRL, false)

	for i := 0; i < 100; i++ {
		time.Sleep(1 * time.Millisecond)

		if val, err := n.getBit(NAU7802_PU_CTRL_PUR, NAU7802_PU_CTRL); err != nil && !val {
			return nil
		}
	}

	return errors.New("timeout")
}

func (n *NAU7802) setIntPolarityHigh() error {
	return n.setBit(NAU7802_CTRL1_CRP, NAU7802_CTRL1, false)
}

func (n *NAU7802) setIntPolarityLow() error {
	return n.setBit(NAU7802_CTRL1_CRP, NAU7802_CTRL1, true)
}

func (n *NAU7802) getRevisionCode() ([]byte, error) {
	return n.getRegister(byte(NAU7802_DEVICE_REV))
}

var (
	nau7802 = &NAU7802{}
)

func Initialize() (*NAU7802, error) {
	nau7802, err := NewNAU7802()

	nau7802.zeroOffset = 0
	nau7802.calibrationFactor = 1.0

	if err != nil {
		log.Fatal("new nau", err)
		return nil, err
	}

	if !nau7802.isConnected() {
		log.Fatal("nau not connected")
		if !nau7802.isConnected() {
			log.Fatal("nau not connected")
			return nil, err
		}
	}

	nau7802.setChannel(NAU7802_CHANNEL_2)

	err = nau7802.reset()
	if err != nil {
		log.Fatal("reset", err)
		return nil, err
	}

	if err = nau7802.powerUp(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = nau7802.setLDO(NAU7802_LDO_3V3); err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = nau7802.setGain(NAU7802_GAIN_128); err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = nau7802.setSampleRate(NAU7802_SPS_80); err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = nau7802.setRegister(NAU7802_ADC, []byte{0x30}); err != nil {
		log.Fatal(err)
		return nil, err
	}

	time.Sleep(100 * time.Millisecond)

	if err = nau7802.setBit(NAU7802_PGA_PWR_PGA_CAP_EN, NAU7802_PGA_PWR, true); err != nil {
		log.Fatal(err)
		return nil, err
	}

	time.Sleep(100 * time.Millisecond)

	nau7802.setGain(NAU7802_GAIN_128)

	if err = nau7802.calibrateAFE(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	time.Sleep(100 * time.Millisecond)

	nau7802.setCalibrationFactor(153.52 / 2.5)
	nau7802.setZeroOffset(16754344)

	return nau7802, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	nau7802, err := Initialize()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	initWeight, _ := nau7802.getWeight(true, 1)

	if initWeight == 0 {
		initWeight, _ = nau7802.getWeight(true, 1)

		if initWeight == 0 {
			nau7802.reset()

			nau7802, err = Initialize()
			if err != nil {
				log.Fatal(err)
			}

			initWeight, _ = nau7802.getWeight(true, 1)
		}
	}

	for {
		weight, err := nau7802.getWeight(true, 1)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(weight - initWeight)
		time.Sleep(500 * time.Millisecond)
	}
}
