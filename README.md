# go-gdiplus

**English** | [简体中文](./README_zh-CN.md)

`go-gdiplus` is a Windows GDI+ (`gdiplus.dll`) wrapper library for the Go programming language. It provides both the low-level flat GDI+ API and a high-level object-oriented interface for 2D vector graphics, typography, and image processing on Windows.

Forked from `shahfarhadreza/go-gdiplus` and expanded with gradient brushes, rounded shapes, clipping helpers, private font collections, automatic multi-font/emoji glyph fallback, in-memory stream encoding, and Go `image.Image` interoperability.

---

## Features

- **Automatic Multi-Font & Emoji Fallback**: Parses OpenType `cmap` tables (Format 4 BMP and Format 12 full Unicode `U+0000..U+10FFFF`) via `gdi32.dll` and caches glyph coverage in memory. Automatically falls back to system fonts (`Segoe UI Emoji`, `Segoe UI Symbol`, `Microsoft YaHei`, `SimSun`, etc.) when characters or emojis are missing from the primary font.
- **Gradient & Pattern Brushes**: Full support for `LinearGradientBrush` (2-point, rectangle modes, custom angles, and multi-stop color interpolation), `PathGradientBrush`, `HatchBrush`, and `TextureBrush`, plus one-line gradient fill helpers on `Graphics`.
- **Vector Paths, Rounded Shapes & Clipping**: Built-in `GraphicsPath` methods (`AddRoundedRectangle`, `AddString`, `AddPolygon`, `AddBezier`, `AddEllipse`), `DrawRoundedRectangle`, `FillRoundedRectangle`, `DrawOutlinedString`, `DrawImageCircle`, and `DrawImageRounded`.
- **Private Font Loading (`PrivateFontCollection`)**: Load `.ttf` / `.otf` files or raw byte slices directly from disk or memory without installing fonts into the operating system.
- **In-Memory Image Encoding & `image.Image` Interop**: `ToPngBytes()` and `ToJpegBytes(quality)` encode directly in memory via Windows COM `IStream` (`CreateStreamOnHGlobal`) without temporary files. Seamless conversion between `*gdiplus.Bitmap` and Go's standard `image.Image`.
- **Safe Resource Management**: All `Dispose()` methods are nil-safe and idempotent, preventing double-free errors.

---

## Installation

```bash
go get -u github.com/furobert/go-gdiplus
```

Import the package:

```go
import gdi "github.com/furobert/go-gdiplus"
```

---

## Quick Start

```go
package main

import (
	"os"
	gdi "github.com/furobert/go-gdiplus"
)

func main() {
	// 1. Initialize GDI+ once at application startup
	input := &gdi.GdiplusStartupInput{GdiplusVersion: 1}
	output := &gdi.GdiplusStartupOutput{}
	if gdi.GdiplusStartup(input, output) != gdi.Ok {
		panic("failed to initialize GDI+")
	}
	defer gdi.GdiplusShutdown()

	// 2. Create an 800x450 Bitmap and Graphics context
	bmp := gdi.NewBitmap(800, 450, gdi.PixelFormat32bppARGB)
	defer bmp.Dispose()

	g := gdi.NewGraphicsFromBitmap(bmp)
	defer g.Dispose()
	g.EnableHighQuality() // Enables anti-aliasing, bicubic interpolation, and ClearType/AntiAlias text

	// 3. Fill background with a diagonal linear gradient
	cStart, _ := gdi.NewColorFromHex("#4FACFE")
	cEnd, _ := gdi.NewColorFromHex("#00F2FE")
	g.FillLinearGradientRect(0, 0, 800, 450, cStart, cEnd, gdi.LinearGradientModeForwardDiagonal)

	// 4. Draw a semi-transparent rounded rectangle card
	cardBrush := gdi.NewSolidBrush(gdi.NewColor(255, 255, 255, 230))
	defer cardBrush.Dispose()
	g.FillRoundedRectangle(cardBrush.AsBrush(), 40, 40, 720, 370, 20)

	borderPen := gdi.NewPen(gdi.NewColor(255, 255, 255, 255), 2)
	defer borderPen.Dispose()
	g.DrawRoundedRectangle(borderPen, 40, 40, 720, 370, 20)

	// 5. Draw text with automatic multi-font and Emoji fallback
	font := gdi.NewFont("Segoe UI, Microsoft YaHei", 32, gdi.FontStyleBold, gdi.UnitPixel, nil)
	defer font.Dispose()

	textBrush := gdi.NewSolidBrush(gdi.NewColorRGB(30, 41, 59))
	defer textBrush.Dispose()

	g.DrawString("Hello GDI+ 🚀 ★ Unicode & Emoji Support", font, gdi.NewPointF(80, 90), textBrush.AsBrush())

	// 6. Encode to PNG or JPEG bytes in memory
	pngBytes, err := bmp.ToPngBytes()
	if err == nil {
		_ = os.WriteFile("output.png", pngBytes, 0644)
	}
}
```

---

## Usage Guide

### 1. Gradients (`LinearGradientBrush`)

```go
c1 := gdi.NewColorRGB(67, 233, 123)
c2 := gdi.NewColorRGB(56, 249, 215)

// Fill a rectangle with a 2-color linear gradient
// Modes:
// - gdi.LinearGradientModeHorizontal       (Left -> Right)
// - gdi.LinearGradientModeVertical         (Top -> Bottom)
// - gdi.LinearGradientModeForwardDiagonal  (Top-Left -> Bottom-Right)
// - gdi.LinearGradientModeBackwardDiagonal (Top-Right -> Bottom-Left)
g.FillLinearGradientRect(0, 0, 800, 400, c1, c2, gdi.LinearGradientModeVertical)

// Fill a rectangle with a linear gradient at a custom angle (in degrees)
g.FillLinearGradientRectAngle(0, 0, 800, 400, c1, c2, 135.0)

// Fill a rectangle with a multi-stop gradient (3 or more colors)
c3 := gdi.NewColorRGB(79, 172, 254)
colors := []*gdi.Color{c1, c2, c3}
positions := []float32{0.0, 0.5, 1.0} // Must start at 0.0 and end at 1.0
g.FillMultiGradientRect(0, 0, 800, 400, colors, positions, gdi.LinearGradientModeHorizontal)

// Fill a rounded rectangle with a linear gradient
g.FillRoundedLinearGradientRect(50, 50, 300, 80, 16, c1, c2, gdi.LinearGradientModeHorizontal)

// Create a standalone LinearGradientBrush for arbitrary shapes or paths
gradBrush := gdi.NewLinearGradientBrush(gdi.NewPointF(0, 0), gdi.NewPointF(0, 400), c1, c2)
defer gradBrush.Dispose()
g.FillEllipse(gradBrush.AsBrush(), 100, 100, 200, 200)
```

### 2. Automatic Font Fallback & Emoji Rendering

Automatic font fallback is enabled by default (`AutoFontFallback = true`). When calling `g.DrawString`, `g.DrawStringEx`, `g.MeasureString`, or `g.MeasureStringEx`:
- If all characters in `text` are supported by the primary `Font`, `go-gdiplus` invokes `GdipDrawString` directly.
- If any character (such as emojis `😀`, `🔥`, symbols `★`, or CJK characters) is missing from the primary font, `go-gdiplus` segments the text into contiguous runs and renders missing glyphs using matching fonts from the fallback chain while preserving `emSize`, `style`, `unit`, alignment, and word wrapping.

```go
// Pass a CSS-style comma-separated font stack in NewFont
font := gdi.NewFont("Arial, Microsoft YaHei, Segoe UI Emoji", 24, gdi.FontStyleRegular, gdi.UnitPixel, nil)
defer font.Dispose()

// Or attach custom fallback families to a specific Font instance
font.WithFallbacks("Segoe UI Emoji", "SimSun")

// Customize the global fallback chain
gdi.AddFallbackFont("Segoe UI Symbol", true) // prepend = true
gdi.SetFallbackFonts([]string{"Segoe UI Emoji", "Microsoft YaHei", "Arial"})

// Disable or enable automatic font fallback globally
gdi.SetAutoFontFallback(true)

// Check whether a font family supports a specific Unicode rune
supported := gdi.FontSupportsRune("Arial", '😀')
```

### 3. Private Font Collections (`PrivateFontCollection`)

Load custom `.ttf` or `.otf` font files from disk or memory without installing them on the system:

```go
pfc := gdi.NewPrivateFontCollection()
defer pfc.Dispose()

// Load from file path (or use pfc.AddMemoryFont(fontBytes) to load from a []byte slice)
_ = pfc.AddFontFile("./fonts/OpenSans-Regular.ttf")

families := pfc.Families()
if len(families) > 0 {
	font := gdi.NewFont(families[0], 28, gdi.FontStyleRegular, gdi.UnitPixel, pfc.AsFontCollection())
	defer font.Dispose()
}
```

### 4. Text Formatting & Outlined Text (`StringFormat` & `DrawOutlinedString`)

```go
// Create a horizontally and vertically centered StringFormat
centerFmt := gdi.NewCenteredStringFormat()
defer centerFmt.Dispose()

// Configure alignment, wrapping, and ellipsis trimming
fmtCustom := gdi.NewStringFormat()
defer fmtCustom.Dispose()
fmtCustom.SetAlignment(gdi.StringAlignmentNear)            // Near (Left), Center, Far (Right)
fmtCustom.SetLineAlignment(gdi.StringAlignmentCenter)      // Near (Top), Center (Middle), Far (Bottom)
fmtCustom.SetFormatFlags(gdi.StringFormatFlagsNoWrap)      // Disable wrapping
fmtCustom.SetTrimming(gdi.StringTrimmingEllipsisCharacter) // Append "..." when clipped

rect := gdi.NewRectF(50, 50, 400, 60)
g.DrawStringEx("Long text content that will be trimmed with an ellipsis...", font, rect, fmtCustom, textBrush.AsBrush())

// Draw outlined/stroked text using GraphicsPath
fillBrush := gdi.NewSolidBrush(gdi.NewColorRGB(255, 255, 255))
defer fillBrush.Dispose()
outlinePen := gdi.NewPen(gdi.NewColorRGB(30, 41, 59), 4)
outlinePen.SetLineJoin(gdi.LineJoinRound)
defer outlinePen.Dispose()

g.DrawOutlinedString("Outlined Title", "Segoe UI", gdi.FontStyleBold, 48, gdi.NewRectF(50, 150, 500, 80), nil, fillBrush.AsBrush(), outlinePen)
```

### 5. Rounded Shapes, Image Clipping & World Transforms

```go
img := gdi.NewImageFromFile("photo.jpg")
defer img.Dispose()

// Draw an image clipped to a circle at (60, 60) with diameter 120
g.DrawImageCircle(img, 60, 60, 120)

// Draw an image clipped to a rounded rectangle with corner radius 16
g.DrawImageRounded(img, 220, 60, 200, 120, 16)

// Save graphics state, apply rotation around (320, 120), and restore state
state := g.Save()
g.RotateAt(15.0, 320, 120)
g.DrawImageRect(img, 220, 60, 200, 120)
g.Restore(state)
```

### 6. In-Memory Streams & `image.Image` Conversion

```go
// Decode a Bitmap directly from encoded []byte in memory (PNG, JPEG, BMP, GIF)
bmpFromMem, err := gdi.NewBitmapFromBytes(encodedBytes)
if err == nil {
	defer bmpFromMem.Dispose()
}

// Convert between *gdi.Bitmap and Go standard library image.Image
stdImg, err := bmp.ToGoImage()            // *gdi.Bitmap -> *image.NRGBA
bmp2 := gdi.NewBitmapFromGoImage(stdImg)  // image.Image -> *gdi.Bitmap
defer bmp2.Dispose()

// Encode to PNG or JPEG byte slices in memory
pngData, err := bmp.ToPngBytes()
jpgData, err := bmp.ToJpegBytes(85)

// Save directly to disk
_ = bmp.SaveAsPng("output.png")
_ = bmp.SaveAsJpeg("output.jpg", 90)
```

---

## Resource Management

GDI+ objects allocate unmanaged system resources outside Go's garbage collector. Always pair every constructor (`NewBitmap`, `NewImageFromFile`, `NewGraphicsFromBitmap`, `NewFont`, `NewPen`, `NewSolidBrush`, `NewLinearGradientBrush`, `NewPath`, `NewStringFormat`, `NewPrivateFontCollection`) with `defer obj.Dispose()`, or call `.Dispose()` explicitly when creating objects inside loops.
