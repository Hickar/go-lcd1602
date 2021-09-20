package main

import (
	"fmt"
	"strconv"
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

	E_PULSE = time.Duration(time.Microsecond * 50)
	E_DELAY = time.Duration(time.Microsecond * 50)

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

func (l *LCD) WriteByte(bits byte, mode byte) {
	highBits := mode | (bits & 0xF0) | BACKLIGHT
	lowBits := mode | ((bits << 4) & 0xF0) | BACKLIGHT

	l.WriteBytes([]byte{highBits})
	l.ToggleEnable(highBits)

	l.WriteBytes([]byte{lowBits})
	l.ToggleEnable(lowBits)
}

func (l *LCD) WriteString(message string, line byte) {
	message = LeftAlign(message)

	l.WriteByte(line, MODE_CMD)

	for i := range message {
		l.WriteByte(byte(rune(message[i])), MODE_CHR)
	}
}

func (l *LCD) ToggleEnable(bits byte) {
	time.Sleep(E_DELAY)
	l.WriteBytes([]byte{bits | ENABLE})
	time.Sleep(E_PULSE)
	l.WriteBytes([]byte{bits & ^ENABLE})
	time.Sleep(E_DELAY)
}

func LeftAlign(str string) string {
	return fmt.Sprintf("%" + strconv.Itoa(WIDTH) + "v", str)
}

func RightAlign(str string) string {
	return fmt.Sprintf("%-" + strconv.Itoa(WIDTH) + "v", str)
}