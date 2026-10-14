package gdiplus

import (
	"fmt"
	"strconv"
	"strings"
)

type Color struct {
	Argb ARGB
}

func MakeARGB(a, r, g, b byte) ARGB {
	return ((ARGB(b) << BlueShift) | (ARGB(g) << GreenShift) | (ARGB(r) << RedShift) | (ARGB(a) << AlphaShift))
}

func NewColor(r, g, b, a byte) *Color {
	return &Color{
		Argb: MakeARGB(a, r, g, b),
	}
}

// NewColorRGB creates an opaque Color (alpha = 255).
func NewColorRGB(r, g, b byte) *Color {
	return &Color{
		Argb: MakeARGB(255, r, g, b),
	}
}

// NewColorFromHex parses "#RRGGBB", "RRGGBB", "#AARRGGBB", or "AARRGGBB".
// If parsing fails, returns opaque black and an error.
func NewColorFromHex(hexStr string) (*Color, error) {
	s := strings.TrimSpace(strings.TrimPrefix(hexStr, "#"))
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	switch len(s) {
	case 6:
		v, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			return NewColor(0, 0, 0, 255), err
		}
		r := byte((v >> 16) & 0xFF)
		g := byte((v >> 8) & 0xFF)
		b := byte(v & 0xFF)
		return NewColor(r, g, b, 255), nil
	case 8:
		v, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			return NewColor(0, 0, 0, 255), err
		}
		a := byte((v >> 24) & 0xFF)
		r := byte((v >> 16) & 0xFF)
		g := byte((v >> 8) & 0xFF)
		b := byte(v & 0xFF)
		return NewColor(r, g, b, a), nil
	default:
		return NewColor(0, 0, 0, 255), fmt.Errorf("invalid hex color string: %q", hexStr)
	}
}

// LerpColor linearly interpolates between c1 and c2 by t in [0.0, 1.0].
func LerpColor(c1, c2 *Color, t float32) *Color {
	if c1 == nil && c2 == nil {
		return NewColor(0, 0, 0, 255)
	}
	if c1 == nil {
		return NewColor(c2.GetR(), c2.GetG(), c2.GetB(), c2.GetA())
	}
	if c2 == nil {
		return NewColor(c1.GetR(), c1.GetG(), c1.GetB(), c1.GetA())
	}
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	r := byte(float32(c1.GetR()) + (float32(c2.GetR())-float32(c1.GetR()))*t)
	g := byte(float32(c1.GetG()) + (float32(c2.GetG())-float32(c1.GetG()))*t)
	b := byte(float32(c1.GetB()) + (float32(c2.GetB())-float32(c1.GetB()))*t)
	a := byte(float32(c1.GetA()) + (float32(c2.GetA())-float32(c1.GetA()))*t)
	return NewColor(r, g, b, a)
}

func (c *Color) GetAlpha() byte {
	if c == nil {
		return 0
	}
	return byte(c.Argb >> AlphaShift)
}

func (c *Color) GetA() byte {
	return c.GetAlpha()
}

func (c *Color) GetRed() byte {
	if c == nil {
		return 0
	}
	return byte(c.Argb >> RedShift)
}

func (c *Color) GetR() byte {
	return c.GetRed()
}

func (c *Color) GetGreen() byte {
	if c == nil {
		return 0
	}
	return byte(c.Argb >> GreenShift)
}

func (c *Color) GetG() byte {
	return c.GetGreen()
}

func (c *Color) GetBlue() byte {
	if c == nil {
		return 0
	}
	return byte(c.Argb >> BlueShift)
}

func (c *Color) GetB() byte {
	return c.GetBlue()
}

func (c *Color) GetValue() ARGB {
	if c == nil {
		return 0
	}
	return c.Argb
}

// WithAlpha returns a copy of the color with a new alpha channel (0-255).
func (c *Color) WithAlpha(a byte) *Color {
	if c == nil {
		return NewColor(0, 0, 0, a)
	}
	return NewColor(c.GetR(), c.GetG(), c.GetB(), a)
}
