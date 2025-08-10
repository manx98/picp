package driver

import (
	_ "embed"
	"errors"
	"github.com/golang/freetype/truetype"
	"go.uber.org/zap"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"picp/config"
	"picp/logger"
	"picp/sh1106"
)

//go:embed fonts/BoutiqueBitmap9x9_1.6.ttf
var boutiqueBitmap9x9 []byte

var BoutiqueBitmap9x9FontFace font.Face

const FontSize = 10

func init() {
	f, err := truetype.Parse(boutiqueBitmap9x9)
	if err != nil {
		panic(err)
	}
	BoutiqueBitmap9x9FontFace = truetype.NewFace(f, &truetype.Options{
		Size: FontSize,
	})
}

var display *sh1106.Device

func sh1106Init() {
	bus, err := config.SH1106.Create()
	if err != nil {
		if !errors.Is(err, config.ErrorSensorDisabled) {
			logger.Fatal("create sh1106 i2c bus failed", zap.Error(err))
		}
	} else {
		display, err = sh1106.NewI2C(bus, sh1106.Config{
			Height:   int16(config.SH1106.Height),
			VccState: config.SH1106.GetMode(),
			Width:    int16(config.SH1106.Width),
		})
		if err != nil {
			logger.Fatal("create sh1106 device failed", zap.Error(err))
		}
	}
}

type DrawOptions struct {
	VerticalAlign   bool
	HorizontalAlign bool
	MarginLeft      int
	MarginRight     int
	MarginTop       int
	MarginBottom    int
}

func (o *DrawOptions) verticalAlign() bool {
	if o != nil {
		return o.VerticalAlign
	}
	return false
}

func (o *DrawOptions) horizontalAlign() bool {
	if o != nil {
		return o.HorizontalAlign
	}
	return false
}

func (o *DrawOptions) marginLeft() int {
	if o != nil {
		return o.MarginLeft
	}
	return 0
}

func (o *DrawOptions) marginRight() int {
	if o != nil {
		return o.MarginRight
	}
	return 0
}

func (o *DrawOptions) marginTop() int {
	if o != nil {
		return o.MarginTop
	}
	return 0
}

func (o *DrawOptions) marginBottom() int {
	if o != nil {
		return o.MarginBottom
	}
	return 0
}

func drawText(width, height int, opt *DrawOptions, lines ...string) *image.Gray {
	dst := image.NewGray(image.Rect(0, 0, width, height))
	draw.Draw(dst, dst.Bounds(), image.Black, image.Point{}, draw.Src)
	d := &font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: BoutiqueBitmap9x9FontFace,
	}
	var yOffset int
	if opt.verticalAlign() {
		yOffset = (height - len(lines)*FontSize - opt.marginTop() - opt.marginBottom()) / 2
	}
	for i, line := range lines {
		if opt.horizontalAlign() {
			strWidth := d.MeasureString(line).Round() + opt.marginRight() + opt.marginLeft()
			d.Dot = fixed.P(opt.marginLeft()+(width-strWidth)/2, yOffset+FontSize*(i+1)+opt.marginTop())
		} else {
			d.Dot = fixed.P(opt.marginLeft(), yOffset+FontSize*(i+1)+opt.marginTop())
		}
		d.DrawString(line)
	}
	return dst
}

func Display(opt *DrawOptions, lines ...string) error {
	if display == nil {
		return nil
	}
	return display.DisplayImage(drawText(display.GetWidth(), display.GetHeight(), opt, lines...))
}
