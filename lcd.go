package main

import (
	"time"

	"github.com/d2r2/go-i2c"
)

const (
	MODE_CMD = 0
	MODE_CHR = 1
	BACKLIGHT = 0x08

	WIDTH = 16

	LINE_1 = 0x80
	LINE_2 = 0xC0

	E_PULSE = time.Microsecond * 50
	E_DELAY = time.Microsecond * 50

	ENABLE = byte(0b00000100)
)

type LCD struct {
	*i2c.I2C
}

func NewLCD(addr uint8, bus int) (*LCD, error) {
	ic2, err := i2c.NewI2C(addr, bus)
	return &LCD{ic2}, err
}

func (l *LCD) Init() {
	l.WriteByte(0x33, MODE_CMD)
	l.WriteByte(0x32, MODE_CMD)
	l.WriteByte(0x06, MODE_CMD)
	l.WriteByte(0x0C, MODE_CMD)
	l.WriteByte(0x28, MODE_CMD)
	l.WriteByte(0x01, MODE_CMD)
	time.Sleep(E_DELAY)
}

func (l *LCD) WriteByte(bits, mode byte) error {
	highBits := mode | (bits & 0xF0) | BACKLIGHT
	lowBits := mode | ((bits << 4) & 0xF0) | BACKLIGHT

	_, err := l.WriteBytes([]byte{highBits})
	if err != nil {
		return err
	}
	l.ToggleEnable(highBits)

	_, err = l.WriteBytes([]byte{lowBits})
	if err != nil {
		return err
	}
	l.ToggleEnable(lowBits)

	return nil
}

func (l *LCD) WriteString(message string, line byte) error {
	message = LeftAlign(message)

	err := l.WriteByte(line, MODE_CMD)
	if err != nil {
		return err
	}

	for i := range message {
		err = l.WriteByte(byte(rune(message[i])), MODE_CHR)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *LCD) ToggleEnable(bits byte) error {
	time.Sleep(E_DELAY)
	_, err := l.WriteBytes([]byte{bits | ENABLE})
	if err != nil {
		return err
	}

	time.Sleep(E_PULSE)
	_, err = l.WriteBytes([]byte{bits & ^ENABLE})
	if err != nil {
		return err
	}
	time.Sleep(E_DELAY)

	return nil
}