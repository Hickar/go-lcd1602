package main

import (
	"errors"
	"time"

	"github.com/d2r2/go-i2c"
)

const (
	MODE_CMD  = 0
	MODE_DATA = 1
	BACKLIGHT = 0x08

	WIDTH = 16

	LCD_LINE_1 = 0x80
	LCD_LINE_2 = 0xC0

	LCD_LEFT  = "left"
	LCD_RIGHT = "right"

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
	l.ClearDisplay()
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

func (l *LCD) ReadByte(bits, mode byte) error {
	highBits := mode | (bits & 0xF0) | BACKLIGHT
	lowBits := mode | ((bits << 4) & 0xF0) | BACKLIGHT

	_, err := l.ReadBytes([]byte{highBits})
	if err != nil {
		return err
	}
	l.ToggleEnable(highBits)

	_, err = l.ReadBytes([]byte{lowBits})
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
		err = l.WriteByte(byte(rune(message[i])), MODE_DATA)
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

func (l *LCD) ClearDisplay() error {
	return l.WriteByte(0x01, MODE_CMD)
}

func (l *LCD) CursorHome() error {
	return l.WriteByte(0x02, MODE_CMD)
}

func (l *LCD) CursorMove(direction string) error {
	if direction != LCD_LEFT && direction != LCD_RIGHT {
		return errors.New("display direction must be \"left\" or \"right\"")
	}

	mask := 0b00010000

	if direction == LCD_LEFT {
		mask |= 0b00000100
	}

	return l.WriteByte(byte(mask), MODE_CMD)
}

func (l *LCD) CursorMoveLeft() error {
	return l.CursorMove(LCD_LEFT)
}

func (l *LCD) CursorMoveRight() error {
	return l.CursorMove(LCD_RIGHT)
}

func (l *LCD) FunctionSet(dataLength, lineNumber, font bool) error {
	mask := 0b00100000

	if dataLength {
		mask |= 0b00010000
	}

	if lineNumber {
		mask |= 0b00001000
	}

	if font {
		mask |= 0b00000100
	}

	return l.ReadByte(byte(mask), MODE_CMD)
}

func (l *LCD) DisplayShift(direction string) error {
	if direction != LCD_LEFT && direction != LCD_RIGHT {
		return errors.New("display direction must be \"left\" or \"right\"")
	}

	mask := 0b00011000

	if direction == LCD_LEFT {
		mask |= 0b00000100
	}

	return l.WriteByte(byte(mask), MODE_CMD)
}

func (l *LCD) DisplayShiftLeft() error {
	return l.DisplayShift(LCD_LEFT)
}

func (l *LCD) DisplayShiftRight() error {
	return l.DisplayShift(LCD_RIGHT)
}

func (l *LCD) Display(displayOn, cursorOn, cursorPositionOn bool) error {
	mask := 0b00001000

	if displayOn {
		mask |= 0b00000100
	}

	if cursorOn {
		mask |= 0b00000010
	}

	if cursorPositionOn {
		mask |= 0b00000001
	}

	return l.WriteByte(byte(mask), MODE_CMD)
}

func (l *LCD) SetCGRAMAddress(bits byte) error {
	return l.ReadByte(bits&(0x1<<6), MODE_CMD)
}

func (l *LCD) SetDDRAMAddress(bits byte) error {
	return l.ReadByte(bits&(0x1<<7), MODE_CMD)
}