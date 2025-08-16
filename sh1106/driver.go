package sh1106

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"picp/go-i2c"
	"sync"
	"time"
)

// Registers
const (
	SETCONTRAST                          = 0x81
	DISPLAYALLON_RESUME                  = 0xA4
	DISPLAYALLON                         = 0xA5
	NORMALDISPLAY                        = 0xA6
	INVERTDISPLAY                        = 0xA7
	DISPLAYOFF                           = 0xAE
	DISPLAYON                            = 0xAF
	SETDISPLAYOFFSET                     = 0xD3
	SETCOMPINS                           = 0xDA
	SETVCOMDETECT                        = 0xDB
	SETDISPLAYCLOCKDIV                   = 0xD5
	SETPRECHARGE                         = 0xD9
	SETMULTIPLEX                         = 0xA8
	SETLOWCOLUMN                         = 0x00
	SETHIGHCOLUMN                        = 0x10
	SETSTARTLINE                         = 0x40
	MEMORYMODE                           = 0x20
	COLUMNADDR                           = 0x21
	PAGEADDR                             = 0x22
	COMSCANINC                           = 0xC0
	COMSCANDEC                           = 0xC8
	SEGREMAP                             = 0xA0
	CHARGEPUMP                           = 0x8D
	ACTIVATE_SCROLL                      = 0x2F
	DEACTIVATE_SCROLL                    = 0x2E
	SET_VERTICAL_SCROLL_AREA             = 0xA3
	RIGHT_HORIZONTAL_SCROLL              = 0x26
	LEFT_HORIZONTAL_SCROLL               = 0x27
	VERTICAL_AND_RIGHT_HORIZONTAL_SCROLL = 0x29
	VERTICAL_AND_LEFT_HORIZONTAL_SCROLL  = 0x2A

	ExternalVCC  VccMode = 0x1
	SwitchCAPVCC VccMode = 0x2
)

type DataBuilder struct {
	data [][]byte
}

func (b *DataBuilder) WriteCmd(cmd ...byte) *DataBuilder {
	for _, c := range cmd {
		b.data = append(b.data, []byte{0x00, c})
	}
	return b
}

func (b *DataBuilder) WriteData(data []byte) *DataBuilder {
	data = append([]byte{0x40}, data...)
	b.data = append(b.data, data)
	return b
}

// Device wraps an SPI connection.
type Device struct {
	bus          *i2c.I2C
	buffer       []byte
	width        int16
	height       int16
	bufferSize   int16
	vccState     VccMode
	lock         sync.Mutex
	updatedPages int64
}

// Config is the configuration for the display
type Config struct {
	Width    int16
	Height   int16
	VccState VccMode
}

type VccMode uint8

// NewI2C creates a new SSD1306 connection. The I2C wire must already be configured.
func NewI2C(bus *i2c.I2C, cfg Config) (d *Device, err error) {
	d = new(Device)
	if cfg.Width != 0 {
		d.width = cfg.Width
	} else {
		d.width = 128
	}
	if cfg.Height != 0 {
		d.height = cfg.Height
	} else {
		d.height = 64
	}
	d.bus = bus
	if cfg.VccState != 0 {
		d.vccState = cfg.VccState
	} else {
		d.vccState = SwitchCAPVCC
	}
	d.bufferSize = d.width * d.height / 8
	d.buffer = make([]byte, d.bufferSize)

	err = d.Reset()
	if err != nil {
		return nil, fmt.Errorf("reset SH1106 occur error: %w", err)
	}
	return d, nil
}

func (d *Device) Reset() error {
	err := d.tx(func(builder *DataBuilder) error {
		builder.WriteCmd(DISPLAYOFF)
		builder.WriteCmd(SETDISPLAYCLOCKDIV)
		builder.WriteCmd(0x80)
		builder.WriteCmd(SETMULTIPLEX)
		builder.WriteCmd(uint8(d.height - 1))
		builder.WriteCmd(SETDISPLAYOFFSET)
		builder.WriteCmd(0x0)
		builder.WriteCmd(SETSTARTLINE | 0x0)
		builder.WriteCmd(CHARGEPUMP)
		if d.vccState == ExternalVCC {
			builder.WriteCmd(0x10)
		} else {
			builder.WriteCmd(0x14)
		}
		builder.WriteCmd(MEMORYMODE)
		builder.WriteCmd(0x00)
		builder.WriteCmd(SEGREMAP | 0x1)
		builder.WriteCmd(COMSCANDEC)

		if (d.width == 128 && d.height == 64) || (d.width == 64 && d.height == 48) { // 128x64 or 64x48
			builder.WriteCmd(SETCOMPINS)
			builder.WriteCmd(0x12)
			builder.WriteCmd(SETCONTRAST)
			if d.vccState == ExternalVCC {
				builder.WriteCmd(0x9F)
			} else {
				builder.WriteCmd(0xCF)
			}
		} else if d.width == 128 && d.height == 32 { // 128x32
			builder.WriteCmd(SETCOMPINS)
			builder.WriteCmd(0x02)
			builder.WriteCmd(SETCONTRAST)
			builder.WriteCmd(0x8F)
		} else if d.width == 96 && d.height == 16 { // 96x16
			builder.WriteCmd(SETCOMPINS)
			builder.WriteCmd(0x2)
			builder.WriteCmd(SETCONTRAST)
			if d.vccState == ExternalVCC {
				builder.WriteCmd(0x10)
			} else {
				builder.WriteCmd(0xAF)
			}
		} else {
			// fail silently, it might work
			return errors.New("there's no configuration for this display's size")
		}

		builder.WriteCmd(SETPRECHARGE)
		if d.vccState == ExternalVCC {
			builder.WriteCmd(0x22)
		} else {
			builder.WriteCmd(0xF1)
		}
		builder.WriteCmd(SETVCOMDETECT)
		builder.WriteCmd(0x40)
		builder.WriteCmd(DISPLAYALLON_RESUME)
		builder.WriteCmd(NORMALDISPLAY)
		builder.WriteCmd(DEACTIVATE_SCROLL)
		builder.WriteCmd(DISPLAYON)
		return nil
	})
	time.Sleep(50 * time.Millisecond)
	if err == nil {
		err = d.Display(true)
	}
	return err
}

// ClearBuffer clears the image buffer
func (d *Device) ClearBuffer() {
	for i := int16(0); i < d.bufferSize; i++ {
		d.buffer[i] = 0
	}
}

// ClearDisplay clears the image buffer and clear the display
func (d *Device) ClearDisplay() error {
	d.ClearBuffer()
	return d.Display(false)
}

// Display sends the whole buffer to the screen
func (d *Device) Display(full bool) (err error) {
	return d.tx(func(builder *DataBuilder) error {
		// In the 128x64 (SPI) screen resetting to 0x0 after 128 times corrupt the buffer
		// Since we're printing the whole buffer, avoid resetting it
		if d.width != 128 || d.height != 64 {
			builder.WriteCmd(COLUMNADDR)
			builder.WriteCmd(0)
			builder.WriteCmd(uint8(d.width - 1))
			builder.WriteCmd(PAGEADDR)
			builder.WriteCmd(0)
			builder.WriteCmd(uint8(d.height/8) - 1)
		}
		for pg := uint8(0); pg < uint8(d.height/8); pg++ {
			if d.updatedPages&(1<<pg) == 0 && !full {
				continue
			}
			builder.WriteCmd(0xB0 | (pg & 0x07)) // SET_PAGE_ADDR
			builder.WriteCmd(SETLOWCOLUMN | 2)
			builder.WriteCmd(SETHIGHCOLUMN | 0)
			builder.WriteData(d.buffer[uint16(pg)*0x80 : uint16(pg+1)*0x80])
		}
		d.updatedPages = 0
		return nil
	})
}

func (d *Device) DisplayImage(img *image.Gray) error {
	size := img.Bounds().Size()
	width := d.width
	if size.X < int(width) {
		width = int16(size.X)
	}
	height := d.height
	if size.Y < int(height) {
		height = int16(size.Y)
	}
	for x := int16(0); x < width; x++ {
		for y := int16(0); y < height; y++ {
			c := img.GrayAt(int(x), int(y))
			d.SetPixel(x, y, c)
		}
	}
	return d.Display(false)
}

// SetPixel enables or disables a pixel in the buffer
func (d *Device) SetPixel(x int16, y int16, c color.Gray) {
	if x < 0 || x >= d.width || y < 0 || y >= d.height {
		return
	}
	byteIndex := x + (y/8)*d.width
	pix := uint8(1) << uint8(y%8)
	oldPix := d.buffer[byteIndex] & pix
	if c.Y > 70 {
		d.buffer[byteIndex] |= pix
	} else {
		d.buffer[byteIndex] &^= pix
	}
	if oldPix != (d.buffer[byteIndex] & pix) {
		d.updatedPages |= 1 << uint8(y/8)
	}
}

// GetPixel returns if the specified pixel is on (true) or off (false)
func (d *Device) GetPixel(x int16, y int16) bool {
	if x < 0 || x >= d.width || y < 0 || y >= d.height {
		return false
	}
	byteIndex := x + (y/8)*d.width
	return (d.buffer[byteIndex] >> uint8(y%8) & 0x1) == 1
}

// SetBuffer changes the whole buffer at once
func (d *Device) SetBuffer(buffer []byte) error {
	if int16(len(buffer)) != d.bufferSize {
		//return ErrBuffer
		return errors.New("wrong size buffer")
	}
	for i := int16(0); i < d.bufferSize; i++ {
		d.buffer[i] = buffer[i]
	}
	return nil
}

// Tx sends data to the display
func (d *Device) tx(call func(builder *DataBuilder) error) (err error) {
	builder := new(DataBuilder)
	err = call(builder)
	if err == nil {
		d.lock.Lock()
		defer d.lock.Unlock()
		for _, data := range builder.data {
			_, err = d.bus.WriteBytes(data)
			if err != nil {
				return err
			}
		}
	}
	return
}

func (d *Device) Close() error {
	if d.bus == nil {
		return nil
	}
	d.ClearBuffer()
	_ = d.Display(true)
	return d.bus.Close()
}

func (d *Device) GetWidth() int {
	return int(d.width)
}

func (d *Device) GetHeight() int {
	return int(d.height)
}

// Size returns the current size of the display.
func (d *Device) Size() (w, h int16) {
	return d.width, d.height
}
