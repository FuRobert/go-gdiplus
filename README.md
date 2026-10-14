# go-gdiplus

**English Documentation** | [中文文档](./README_zh-CN.md)

`go-gdiplus` is a lightweight, high-performance Windows GDI+ wrapper library for Go. It provides both the low-level flat `gdiplus.dll` API and an idiomatic object-oriented wrapper designed for high-concurrency image generation, data visualization cards, automated report rendering, and server-side graphic card generation.

Forked from `shahfarhadreza/go-gdiplus` and extensively enhanced by `furobert` with an **Automatic Multi-Font & Emoji/Symbol Fallback Engine**, **One-Line Linear & Multi-Stop Gradients**, **Rounded Rectangles & Circular Avatar Clipping**, **Private Font Collection (`.ttf`/`.otf`) Loading**, **Zero-Temp-File In-Memory COM `IStream` Encoding**, and **Idempotent Resource Disposal**.

---

## Key Features

1. **Automatic Multi-Font & Emoji / Symbol Fallback Engine**
   * Solves the classic GDI+ limitation where a single decorative or custom font displays blank boxes (`□`) for Emojis (`😀`, `🔥`, `📊`), special symbols (`★`, `♪`, `♡`), or characters outside the font's character set.
   * Parses OpenType `cmap` tables (Format 4 BMP and Format 12 full Unicode astral planes `U+0000..U+10FFFF`) directly via `gdi32.dll` (`GetFontData`) and caches glyph coverage in memory (`< 20ns` per rune lookup).
   * Automatically segments text into contiguous font runs and falls back to system fonts (`Segoe UI Emoji`, `Segoe UI Symbol`, `Microsoft YaHei`, `SimSun`, etc.) while preserving `StringFormat` horizontal/vertical alignments and line wrapping.
   * Supports CSS-style comma-separated font stacks (e.g., `"CustomFont, Microsoft YaHei, Segoe UI Emoji"`). If the primary font is not installed on the host machine, it gracefully degrades to available system fonts instead of returning a nil pointer.
2. **Native Gradient Brushes & One-Line Gradient Helpers**
   * Eliminates manual row-by-row `NewPen` + `DrawLineI` loops. Includes full support for `LinearGradientBrush`, `PathGradientBrush`, `HatchBrush`, and `TextureBrush`.
   * Provides `g.FillLinearGradientRect`, `g.FillLinearGradientRectAngle`, `g.FillMultiGradientRect`, and `g.FillRoundedLinearGradientRect` for horizontal, vertical, diagonal, arbitrary-angle, and multi-stop color gradients in a single call.
3. **Rounded Shapes, Circular Image Clipping & Outlined Text**
   * Built-in `FillRoundedRectangle`, `DrawRoundedRectangle`, `FillCircle`, and `DrawCircle`.
   * Built-in `DrawImageCircle` (renders circular avatars/thumbnails in one line) and `DrawImageRounded` (rounded-corner image clipping).
   * Built-in `DrawOutlinedString` using `GraphicsPath.AddString` for stroked/outlined typography.
4. **In-Memory Image Encoding & Go `image.Image` Interoperability**
   * `ToPngBytes()` and `ToJpegBytes(quality)` encode directly in memory via Windows COM `IStream` (`CreateStreamOnHGlobal`), avoiding disk I/O and temporary files.
   * Includes `NewBitmapFromBytes(data []byte)` and seamless conversion to/from Go's standard `image.Image` (`NewBitmapFromGoImage` / `ToGoImage`).
5. **Private Font Loading (`PrivateFontCollection`)**
   * Load `.ttf` / `.otf` files or raw byte slices in memory without installing fonts into `C:\Windows\Fonts`.
6. **Idempotent, Crash-Safe `Dispose()`**
   * Every `Dispose()` method checks for `nil` and clears the underlying native pointer, preventing double-free crashes.

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

	// 3. Draw a diagonal linear gradient background in one line
	cStart, _ := gdi.NewColorFromHex("#FFDEE9")
	cEnd, _ := gdi.NewColorFromHex("#B5FFFC")
	g.FillLinearGradientRect(0, 0, 800, 450, cStart, cEnd, gdi.LinearGradientModeForwardDiagonal)

	// 4. Draw a semi-transparent rounded card with border
	cardBrush := gdi.NewSolidBrush(gdi.NewColor(255, 255, 255, 220))
	defer cardBrush.Dispose()
	g.FillRoundedRectangle(cardBrush.AsBrush(), 40, 40, 720, 370, 24)

	borderPen := gdi.NewPen(gdi.NewColor(255, 182, 193, 255), 3)
	defer borderPen.Dispose()
	g.DrawRoundedRectangle(borderPen, 40, 40, 720, 370, 24)

	// 5. Draw text containing Emojis, symbols, and multilingual characters (automatic font fallback!)
	font := gdi.NewFont("CustomDisplayFont, Microsoft YaHei", 32, gdi.FontStyleBold, gdi.UnitPixel, nil)
	defer font.Dispose()

	textBrush := gdi.NewSolidBrush(gdi.NewColorRGB(50, 50, 70))
	defer textBrush.Dispose()

	g.DrawString("📊 System Performance Metrics 🔥 ★ 100% Complete!", font, gdi.NewPointF(80, 90), textBrush.AsBrush())

	// 6. Export to JPEG bytes in memory or save to disk
	jpegBytes, err := bmp.ToJpegBytes(90)
	if err == nil {
		_ = os.WriteFile("card.jpg", jpegBytes, 0644)
	}
}
```

---

## Feature Guide

### 1. Linear & Multi-Stop Gradients (`LinearGradientBrush`)

```go
c1 := gdi.NewColorRGB(255, 210, 230)
c2 := gdi.NewColorRGB(190, 225, 255)

// Option A: Fill a rectangle with a 2-color gradient
// Modes:
// - gdi.LinearGradientModeHorizontal       (Left -> Right)
// - gdi.LinearGradientModeVertical         (Top -> Bottom)
// - gdi.LinearGradientModeForwardDiagonal  (Top-Left -> Bottom-Right)
// - gdi.LinearGradientModeBackwardDiagonal (Top-Right -> Bottom-Left)
g.FillLinearGradientRect(0, 0, 1000, 600, c1, c2, gdi.LinearGradientModeVertical)

// Option B: Fill a rectangle with a gradient at any custom angle (e.g., 135 degrees)
g.FillLinearGradientRectAngle(0, 0, 1000, 600, c1, c2, 135.0)

// Option C: Fill a rectangle with a multi-stop gradient (3+ colors)
c3 := gdi.NewColorRGB(240, 255, 220)
colors := []*gdi.Color{c1, c2, c3}
positions := []float32{0.0, 0.5, 1.0} // Must start at 0.0 and end at 1.0
g.FillMultiGradientRect(0, 0, 1000, 600, colors, positions, gdi.LinearGradientModeVertical)

// Option D: Fill a rounded rectangle with a linear gradient
g.FillRoundedLinearGradientRect(50, 50, 300, 80, 20, c1, c2, gdi.LinearGradientModeHorizontal)

// Option E: Create a LinearGradientBrush directly for custom shapes/paths
gradBrush := gdi.NewLinearGradientBrush(gdi.NewPointF(0, 0), gdi.NewPointF(0, 600), c1, c2)
defer gradBrush.Dispose()
g.FillEllipse(gradBrush.AsBrush(), 100, 100, 200, 200)
```

### 2. Automatic Font Fallback & Emoji Rendering

Automatic font fallback is enabled by default (`AutoFontFallback = true`). When calling `g.DrawString`, `g.DrawStringEx`, `g.MeasureString`, or `g.MeasureStringEx`:
* If all characters in `text` are supported by the primary `Font`, `go-gdiplus` invokes `GdipDrawString` directly with zero overhead.
* If any character (such as `😀`, `⏳`, `★`, or unmapped glyphs) is missing from the primary font, `go-gdiplus` automatically resolves a matching font from the fallback chain using the exact same `emSize`, `style`, and `unit`.

```go
// 1. Pass a comma-separated font stack in NewFont
font := gdi.NewFont("CustomFont, Microsoft YaHei, Segoe UI Emoji", 28, gdi.FontStyleRegular, gdi.UnitPixel, nil)
defer font.Dispose()

// 2. Attach per-font fallback families
font.WithFallbacks("Segoe UI Emoji", "SimSun")

// 3. Customize the global fallback chain
gdi.AddFallbackFont("HarmonyOS Sans SC", true) // prepend = true

// 4. Query whether a font family supports a specific Unicode rune
supported := gdi.FontSupportsRune("CustomFont", '😀')
```

### 3. Loading Uninstalled `.ttf` / `.otf` Font Files (`PrivateFontCollection`)

```go
pfc := gdi.NewPrivateFontCollection()
defer pfc.Dispose()

// Load from disk (or use pfc.AddMemoryFont(ttfBytes) to load from memory)
_ = pfc.AddFontFile("./assets/fonts/CustomFont.ttf")

families := pfc.Families()
font := gdi.NewFont(families[0], 36, gdi.FontStyleRegular, gdi.UnitPixel, pfc.AsFontCollection())
defer font.Dispose()
```

### 4. Text Alignment, Trimming & Outlined Typography (`StringFormat` & `DrawOutlinedString`)

```go
// Create a horizontally and vertically centered StringFormat
centerFmt := gdi.NewCenteredStringFormat()
defer centerFmt.Dispose()

// Custom alignment and ellipsis trimming
fmtCustom := gdi.NewStringFormat()
defer fmtCustom.Dispose()
fmtCustom.SetAlignment(gdi.StringAlignmentNear)
fmtCustom.SetLineAlignment(gdi.StringAlignmentCenter)
fmtCustom.SetFormatFlags(gdi.StringFormatFlagsNoWrap)
fmtCustom.SetTrimming(gdi.StringTrimmingEllipsisCharacter)

rect := gdi.NewRectF(50, 50, 400, 60)
g.DrawStringEx("A very long string that will be trimmed with an ellipsis...", font, rect, fmtCustom, textBrush.AsBrush())

// Outlined/Stroked text
fillBrush := gdi.NewSolidBrush(gdi.NewColorRGB(255, 255, 255))
defer fillBrush.Dispose()
outlinePen := gdi.NewPen(gdi.NewColorRGB(255, 105, 180), 4)
outlinePen.SetLineJoin(gdi.LineJoinRound)
defer outlinePen.Dispose()

g.DrawOutlinedString("GRAPHICS ENGINE", "Microsoft YaHei", gdi.FontStyleBold, 48, gdi.NewRectF(80, 200, 600, 80), nil, fillBrush.AsBrush(), outlinePen)
```

### 5. Circular Avatars, Rounded Images & Transforms

```go
avatar := gdi.NewImageFromFile("avatar.jpg")
defer avatar.Dispose()

// Draw a circular image at (60, 60) with diameter 120
g.DrawImageCircle(avatar, 60, 60, 120)

// Draw a rounded-corner image with corner radius 16
g.DrawImageRounded(avatar, 220, 60, 200, 120, 16)

// Rotate around a specific pivot point
state := g.Save()
g.RotateAt(15.0, 320, 120)
g.DrawImageRect(avatar, 220, 60, 200, 120)
g.Restore(state)
```

### 6. In-Memory Streams & Go `image.Image` Interop

```go
// 1. Decode a Bitmap directly from []byte in memory
bmpFromMem, err := gdi.NewBitmapFromBytes(rawBytes)
if err == nil {
	defer bmpFromMem.Dispose()
}

// 2. Convert between *gdi.Bitmap and Go standard image.Image
goImg, err := bmp.ToGoImage()
bmp2 := gdi.NewBitmapFromGoImage(goImg)
defer bmp2.Dispose()

// 3. Encode to PNG or JPEG bytes in memory via COM IStream
pngBytes, err := bmp.ToPngBytes()
jpgBytes, err := bmp.ToJpegBytes(85)
```

---

## Resource Management Best Practices

Because GDI+ objects are allocated in unmanaged Windows memory outside Go's garbage collector:
1. Always pair every `New*` call (`NewBitmap`, `NewImageFromFile`, `NewGraphicsFromBitmap`, `NewFont`, `NewPen`, `NewSolidBrush`, `NewLinearGradientBrush`, `NewPath`, `NewStringFormat`, `NewPrivateFontCollection`) with `defer obj.Dispose()`.
2. Avoid creating undisposed GDI+ objects inside tight loops; if you create objects inside a loop, call `.Dispose()` at the end of each iteration. All `.Dispose()` methods in this library are idempotent and nil-safe.