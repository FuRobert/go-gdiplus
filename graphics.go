package gdiplus

import (
	"log"
	"syscall"
)

type Graphics struct {
	nativeGraphics *GpGraphics
}

func NewGraphicsFromHDC(hdc HDC) *Graphics {
	g := &Graphics{}
	GdipCreateFromHDC(hdc, &g.nativeGraphics)
	return g
}

func NewGraphicsFromHDC2(hdc HDC, hDevice HANDLE) *Graphics {
	g := &Graphics{}
	GdipCreateFromHDC2(hdc, hDevice, &g.nativeGraphics)
	return g
}

func NewGraphicsFromHWND(hwnd HWND) *Graphics {
	g := &Graphics{}
	GdipCreateFromHWND(hwnd, &g.nativeGraphics)
	return g
}

func NewGraphicsFromHWNDICM(hwnd HWND) *Graphics {
	g := &Graphics{}
	GdipCreateFromHWNDICM(hwnd, &g.nativeGraphics)
	return g
}

func (g *Graphics) GetDC() (hdc HDC) {
	if g != nil && g.nativeGraphics != nil {
		GdipGetDC(g.nativeGraphics, &hdc)
	}
	return
}

func (g *Graphics) ReleaseDC(hdc HDC) {
	if g != nil && g.nativeGraphics != nil {
		GdipReleaseDC(g.nativeGraphics, hdc)
	}
}

func NewGraphicsFromBitmap(bitmap *Bitmap) *Graphics {
	g := &Graphics{}
	if bitmap != nil && bitmap.nativeImage != nil {
		GdipGetImageGraphicsContext(bitmap.nativeImage, &g.nativeGraphics)
	}
	return g
}

// EnableHighQuality configures anti-aliasing, high-quality bicubic interpolation,
// high-quality compositing, and ClearType/AntiAlias text rendering in a single call.
func (g *Graphics) EnableHighQuality() *Graphics {
	if g == nil || g.nativeGraphics == nil {
		return g
	}
	g.SetSmoothingMode(SmoothingModeAntiAlias)
	g.SetInterpolationMode(InterpolationModeHighQualityBicubic)
	g.SetCompositingQuality(CompositingQualityHighQuality)
	g.SetPixelOffsetMode(PixelOffsetModeHighQuality)
	g.SetTextRenderingHint(TextRenderingHintAntiAliasGridFit)
	return g
}

func (g *Graphics) SetCompositingMode(mode int32) {
	if g != nil && g.nativeGraphics != nil {
		GdipSetCompositingMode(g.nativeGraphics, mode)
	}
}

func (g *Graphics) SetRenderingOrigin(x, y int32) {
	if g != nil && g.nativeGraphics != nil {
		GdipSetRenderingOrigin(g.nativeGraphics, x, y)
	}
}

func (g *Graphics) SetCompositingQuality(quality int32) {
	if g != nil && g.nativeGraphics != nil {
		GdipSetCompositingQuality(g.nativeGraphics, quality)
	}
}

func (g *Graphics) SetInterpolationMode(mode int32) {
	if g != nil && g.nativeGraphics != nil {
		GdipSetInterpolationMode(g.nativeGraphics, mode)
	}
}

func (g *Graphics) SetPixelOffsetMode(mode int32) {
	if g != nil && g.nativeGraphics != nil {
		GdipSetPixelOffsetMode(g.nativeGraphics, mode)
	}
}

func (g *Graphics) SetSmoothingMode(mode int32) {
	if g != nil && g.nativeGraphics != nil {
		GdipSetSmoothingMode(g.nativeGraphics, mode)
	}
}

func (g *Graphics) SetTextRenderingHint(hint int32) {
	if g != nil && g.nativeGraphics != nil {
		GdipSetTextRenderingHint(g.nativeGraphics, hint)
	}
}

// ==================== State, Transform & Clipping ====================

// Save saves the current state (transforms, clip region, quality settings) of the Graphics object.
func (g *Graphics) Save() GraphicsState {
	if g == nil || g.nativeGraphics == nil {
		return 0
	}
	var state uint32
	GdipSaveGraphics(g.nativeGraphics, &state)
	return GraphicsState(state)
}

// Restore restores the state of the Graphics object to a previously saved GraphicsState.
func (g *Graphics) Restore(state GraphicsState) {
	if g != nil && g.nativeGraphics != nil {
		GdipRestoreGraphics(g.nativeGraphics, uint32(state))
	}
}

func (g *Graphics) Flush() {
	if g != nil && g.nativeGraphics != nil {
		GdipFlush(g.nativeGraphics, FlushIntentionSync)
	}
}

func (g *Graphics) TranslateTransform(dx, dy float32) {
	if g != nil && g.nativeGraphics != nil {
		GdipTranslateWorldTransform(g.nativeGraphics, dx, dy, MatrixOrderPrepend)
	}
}

func (g *Graphics) ScaleTransform(sx, sy float32) {
	if g != nil && g.nativeGraphics != nil {
		GdipScaleWorldTransform(g.nativeGraphics, sx, sy, MatrixOrderPrepend)
	}
}

func (g *Graphics) RotateTransform(angleDegrees float32) {
	if g != nil && g.nativeGraphics != nil {
		GdipRotateWorldTransform(g.nativeGraphics, angleDegrees, MatrixOrderPrepend)
	}
}

// RotateAt rotates the coordinate system around a specific pivot point (centerX, centerY).
func (g *Graphics) RotateAt(angleDegrees, centerX, centerY float32) {
	if g == nil || g.nativeGraphics == nil {
		return
	}
	GdipTranslateWorldTransform(g.nativeGraphics, centerX, centerY, MatrixOrderPrepend)
	GdipRotateWorldTransform(g.nativeGraphics, angleDegrees, MatrixOrderPrepend)
	GdipTranslateWorldTransform(g.nativeGraphics, -centerX, -centerY, MatrixOrderPrepend)
}

func (g *Graphics) ResetTransform() {
	if g != nil && g.nativeGraphics != nil {
		GdipResetWorldTransform(g.nativeGraphics)
	}
}

func (g *Graphics) SetClipRect(x, y, width, height float32, combineMode CombineMode) {
	if g != nil && g.nativeGraphics != nil {
		GdipSetClipRect(g.nativeGraphics, x, y, width, height, combineMode)
	}
}

func (g *Graphics) SetClipRectI(x, y, width, height int32, combineMode CombineMode) {
	if g != nil && g.nativeGraphics != nil {
		GdipSetClipRectI(g.nativeGraphics, x, y, width, height, combineMode)
	}
}

func (g *Graphics) SetClipPath(path *GraphicsPath, combineMode CombineMode) {
	if g != nil && g.nativeGraphics != nil && path != nil && path.nativePath != nil {
		GdipSetClipPath(g.nativeGraphics, path.nativePath, combineMode)
	}
}

func (g *Graphics) ResetClip() {
	if g != nil && g.nativeGraphics != nil {
		GdipResetClip(g.nativeGraphics)
	}
}

// ==================== Drawing & Filling Shapes ====================

func (g *Graphics) Clear(color *Color) {
	if g != nil && g.nativeGraphics != nil && color != nil {
		GdipGraphicsClear(g.nativeGraphics, color.GetValue())
	}
}

func (g *Graphics) DrawLine(pen *Pen, x1, y1, x2, y2 float32) {
	if g != nil && g.nativeGraphics != nil && pen != nil {
		GdipDrawLine(g.nativeGraphics, pen.nativePen, x1, y1, x2, y2)
	}
}

func (g *Graphics) DrawLineI(pen *Pen, x1, y1, x2, y2 int32) {
	if g != nil && g.nativeGraphics != nil && pen != nil {
		GdipDrawLineI(g.nativeGraphics, pen.nativePen, x1, y1, x2, y2)
	}
}

func (g *Graphics) DrawArc(pen *Pen, x, y, width, height, startAngle, sweepAngle float32) {
	if g != nil && g.nativeGraphics != nil && pen != nil {
		GdipDrawArc(g.nativeGraphics, pen.nativePen, x, y, width, height, startAngle, sweepAngle)
	}
}

func (g *Graphics) DrawArcI(pen *Pen, x, y, width, height int32, startAngle, sweepAngle float32) {
	if g != nil && g.nativeGraphics != nil && pen != nil {
		GdipDrawArcI(g.nativeGraphics, pen.nativePen, x, y, width, height, startAngle, sweepAngle)
	}
}

func (g *Graphics) DrawBezier(pen *Pen, x1, y1, x2, y2, x3, y3, x4, y4 float32) {
	if g != nil && g.nativeGraphics != nil && pen != nil {
		GdipDrawBezier(g.nativeGraphics, pen.nativePen, x1, y1, x2, y2, x3, y3, x4, y4)
	}
}

func (g *Graphics) DrawBezierI(pen *Pen, x1, y1, x2, y2, x3, y3, x4, y4 int32) {
	if g != nil && g.nativeGraphics != nil && pen != nil {
		GdipDrawBezierI(g.nativeGraphics, pen.nativePen, x1, y1, x2, y2, x3, y3, x4, y4)
	}
}

func (g *Graphics) DrawRectangle(pen *Pen, x, y, width, height float32) {
	if g != nil && g.nativeGraphics != nil && pen != nil {
		GdipDrawRectangle(g.nativeGraphics, pen.nativePen, x, y, width, height)
	}
}

func (g *Graphics) DrawRectangleI(pen *Pen, x, y, width, height int32) {
	if g != nil && g.nativeGraphics != nil && pen != nil {
		GdipDrawRectangleI(g.nativeGraphics, pen.nativePen, x, y, width, height)
	}
}

// DrawRoundedRectangle draws a rounded rectangle outline with corner radius.
func (g *Graphics) DrawRoundedRectangle(pen *Pen, x, y, width, height, radius float32) {
	if g == nil || g.nativeGraphics == nil || pen == nil {
		return
	}
	p := NewPath(FillModeAlternate)
	defer p.Dispose()
	p.AddRoundedRectangle(x, y, width, height, radius)
	g.DrawPath(pen, p)
}

func (g *Graphics) DrawEllipse(pen *Pen, x, y, width, height float32) {
	if g != nil && g.nativeGraphics != nil && pen != nil {
		GdipDrawEllipse(g.nativeGraphics, pen.nativePen, x, y, width, height)
	}
}

func (g *Graphics) DrawEllipseI(pen *Pen, x, y, width, height int32) {
	if g != nil && g.nativeGraphics != nil && pen != nil {
		GdipDrawEllipseI(g.nativeGraphics, pen.nativePen, x, y, width, height)
	}
}

// DrawCircle draws a circle centered at (centerX, centerY) with the specified radius.
func (g *Graphics) DrawCircle(pen *Pen, centerX, centerY, radius float32) {
	g.DrawEllipse(pen, centerX-radius, centerY-radius, radius*2, radius*2)
}

func (g *Graphics) DrawPie(pen *Pen, x, y, width, height, startAngle, sweepAngle float32) {
	if g != nil && g.nativeGraphics != nil && pen != nil {
		GdipDrawPie(g.nativeGraphics, pen.nativePen, x, y, width, height, startAngle, sweepAngle)
	}
}

func (g *Graphics) DrawPieI(pen *Pen, x, y, width, height int32, startAngle, sweepAngle float32) {
	if g != nil && g.nativeGraphics != nil && pen != nil {
		GdipDrawPieI(g.nativeGraphics, pen.nativePen, x, y, width, height, startAngle, sweepAngle)
	}
}

func (g *Graphics) DrawPolygon(pen *Pen, points []PointF) {
	if g != nil && g.nativeGraphics != nil && pen != nil && len(points) > 0 {
		GdipDrawPolygon(g.nativeGraphics, pen.nativePen, &points[0], int32(len(points)))
	}
}

func (g *Graphics) DrawPolygonI(pen *Pen, points []Point) {
	if g != nil && g.nativeGraphics != nil && pen != nil && len(points) > 0 {
		GdipDrawPolygonI(g.nativeGraphics, pen.nativePen, &points[0], int32(len(points)))
	}
}

func (g *Graphics) DrawCurve(pen *Pen, points []PointF) {
	if g != nil && g.nativeGraphics != nil && pen != nil && len(points) > 0 {
		GdipDrawCurve(g.nativeGraphics, pen.nativePen, &points[0], int32(len(points)))
	}
}

func (g *Graphics) DrawClosedCurve(pen *Pen, points []PointF) {
	if g != nil && g.nativeGraphics != nil && pen != nil && len(points) > 0 {
		GdipDrawClosedCurve(g.nativeGraphics, pen.nativePen, &points[0], int32(len(points)))
	}
}

func (g *Graphics) DrawPath(pen *Pen, path *GraphicsPath) {
	if g != nil && g.nativeGraphics != nil && pen != nil && path != nil {
		GdipDrawPath(g.nativeGraphics, pen.nativePen, path.nativePath)
	}
}

func (g *Graphics) FillRectangle(brush *Brush, x, y, width, height float32) {
	if g != nil && g.nativeGraphics != nil && brush != nil {
		GdipFillRectangle(g.nativeGraphics, brush.nativeBrush, x, y, width, height)
	}
}

func (g *Graphics) FillRectangleI(brush *Brush, x, y, width, height int32) {
	if g != nil && g.nativeGraphics != nil && brush != nil {
		GdipFillRectangleI(g.nativeGraphics, brush.nativePenOrBrush(), x, y, width, height)
	}
}

func (b *Brush) nativePenOrBrush() *GpBrush {
	if b == nil {
		return nil
	}
	return b.nativeBrush
}

// FillRoundedRectangle fills a rounded rectangle with the given Brush and corner radius.
func (g *Graphics) FillRoundedRectangle(brush *Brush, x, y, width, height, radius float32) {
	if g == nil || g.nativeGraphics == nil || brush == nil {
		return
	}
	p := NewPath(FillModeAlternate)
	defer p.Dispose()
	p.AddRoundedRectangle(x, y, width, height, radius)
	g.FillPath(brush, p)
}

// FillLinearGradientRect fills a rectangle with a smooth linear gradient in one call (no manual line loops needed).
func (g *Graphics) FillLinearGradientRect(x, y, width, height float32, color1, color2 *Color, mode LinearGradientMode) {
	if g == nil || g.nativeGraphics == nil || width <= 0 || height <= 0 {
		return
	}
	rect := &RectF{X: x, Y: y, Width: width, Height: height}
	grad := NewLinearGradientBrushFromRect(rect, color1, color2, mode)
	defer grad.Dispose()
	g.FillRectangle(grad.AsBrush(), x, y, width, height)
}

// FillLinearGradientRectAngle fills a rectangle with a linear gradient at any angle (in degrees).
func (g *Graphics) FillLinearGradientRectAngle(x, y, width, height float32, color1, color2 *Color, angle float32) {
	if g == nil || g.nativeGraphics == nil || width <= 0 || height <= 0 {
		return
	}
	rect := &RectF{X: x, Y: y, Width: width, Height: height}
	grad := NewLinearGradientBrushWithAngle(rect, color1, color2, angle, true)
	defer grad.Dispose()
	g.FillRectangle(grad.AsBrush(), x, y, width, height)
}

// FillMultiGradientRect fills a rectangle with a multi-stop linear gradient (e.g. 3 or more colors).
// positions must start at 0.0 and end at 1.0.
func (g *Graphics) FillMultiGradientRect(x, y, width, height float32, colors []*Color, positions []float32, mode LinearGradientMode) {
	if g == nil || g.nativeGraphics == nil || len(colors) < 2 || len(colors) != len(positions) {
		return
	}
	rect := &RectF{X: x, Y: y, Width: width, Height: height}
	grad := NewLinearGradientBrushFromRect(rect, colors[0], colors[len(colors)-1], mode)
	defer grad.Dispose()
	_ = grad.SetInterpolationColors(colors, positions)
	g.FillRectangle(grad.AsBrush(), x, y, width, height)
}

// FillRoundedLinearGradientRect fills a rounded rectangle with a linear gradient in one call.
func (g *Graphics) FillRoundedLinearGradientRect(x, y, width, height, radius float32, color1, color2 *Color, mode LinearGradientMode) {
	if g == nil || g.nativeGraphics == nil || width <= 0 || height <= 0 {
		return
	}
	rect := &RectF{X: x, Y: y, Width: width, Height: height}
	grad := NewLinearGradientBrushFromRect(rect, color1, color2, mode)
	defer grad.Dispose()
	g.FillRoundedRectangle(grad.AsBrush(), x, y, width, height, radius)
}

func (g *Graphics) FillEllipse(brush *Brush, x, y, width, height float32) {
	if g != nil && g.nativeGraphics != nil && brush != nil {
		GdipFillEllipse(g.nativeGraphics, brush.nativeBrush, x, y, width, height)
	}
}

func (g *Graphics) FillEllipseI(brush *Brush, x, y, width, height int32) {
	if g != nil && g.nativeGraphics != nil && brush != nil {
		GdipFillEllipseI(g.nativeGraphics, brush.nativeBrush, x, y, width, height)
	}
}

// FillCircle fills a circle centered at (centerX, centerY) with the specified radius.
func (g *Graphics) FillCircle(brush *Brush, centerX, centerY, radius float32) {
	g.FillEllipse(brush, centerX-radius, centerY-radius, radius*2, radius*2)
}

func (g *Graphics) FillPie(brush *Brush, x, y, width, height, startAngle, sweepAngle float32) {
	if g != nil && g.nativeGraphics != nil && brush != nil {
		GdipFillPie(g.nativeGraphics, brush.nativeBrush, x, y, width, height, startAngle, sweepAngle)
	}
}

func (g *Graphics) FillPieI(brush *Brush, x, y, width, height int32, startAngle, sweepAngle float32) {
	if g != nil && g.nativeGraphics != nil && brush != nil {
		GdipFillPieI(g.nativeGraphics, brush.nativeBrush, x, y, width, height, startAngle, sweepAngle)
	}
}

func (g *Graphics) FillPolygon(brush *Brush, points []PointF, fillMode int32) {
	if g != nil && g.nativeGraphics != nil && brush != nil && len(points) > 0 {
		GdipFillPolygon(g.nativeGraphics, brush.nativeBrush, &points[0], int32(len(points)), fillMode)
	}
}

func (g *Graphics) FillPolygonI(brush *Brush, points []Point, fillMode int32) {
	if g != nil && g.nativeGraphics != nil && brush != nil && len(points) > 0 {
		GdipFillPolygonI(g.nativeGraphics, brush.nativeBrush, &points[0], int32(len(points)), fillMode)
	}
}

func (g *Graphics) FillClosedCurve(brush *Brush, points []PointF) {
	if g != nil && g.nativeGraphics != nil && brush != nil && len(points) > 0 {
		GdipFillClosedCurve(g.nativeGraphics, brush.nativeBrush, &points[0], int32(len(points)))
	}
}

func (g *Graphics) FillPath(brush *Brush, path *GraphicsPath) {
	if g != nil && g.nativeGraphics != nil && brush != nil && path != nil {
		GdipFillPath(g.nativeGraphics, brush.nativeBrush, path.nativePath)
	}
}

// ==================== Image Rendering ====================

func (g *Graphics) DrawImage(image *Image, x, y float32) {
	if g == nil || g.nativeGraphics == nil || image == nil || image.nativeImage == nil {
		return
	}
	status := GdipDrawImage(g.nativeGraphics, image.nativeImage, x, y)
	if status != Ok {
		log.Panicln(status.String())
	}
}

func (g *Graphics) DrawImageI(image *Image, x, y int32) {
	if g == nil || g.nativeGraphics == nil || image == nil || image.nativeImage == nil {
		return
	}
	status := GdipDrawImageI(g.nativeGraphics, image.nativeImage, x, y)
	if status != Ok {
		log.Panicln(status.String())
	}
}

func (g *Graphics) DrawImageRect(image *Image, x, y, width, height float32) {
	if g != nil && g.nativeGraphics != nil && image != nil && image.nativeImage != nil {
		GdipDrawImageRect(g.nativeGraphics, image.nativeImage, x, y, width, height)
	}
}

func (g *Graphics) DrawImageRectI(image *Image, x, y, width, height int32) {
	if g != nil && g.nativeGraphics != nil && image != nil && image.nativeImage != nil {
		GdipDrawImageRectI(g.nativeGraphics, image.nativeImage, x, y, width, height)
	}
}

// DrawImageCircle draws an image clipped into a circle at (x, y) with the given diameter.
// Ideal for rendering circular user avatars.
func (g *Graphics) DrawImageCircle(image *Image, x, y, diameter float32) {
	if g == nil || g.nativeGraphics == nil || image == nil || image.nativeImage == nil || diameter <= 0 {
		return
	}
	state := g.Save()
	defer g.Restore(state)

	path := NewPath(FillModeAlternate)
	defer path.Dispose()
	path.AddEllipse(x, y, diameter, diameter)

	g.SetClipPath(path, CombineModeIntersect)
	g.DrawImageRect(image, x, y, diameter, diameter)
}

// DrawImageRounded draws an image clipped to a rounded rectangle at (x, y, width, height) with corner radius.
func (g *Graphics) DrawImageRounded(image *Image, x, y, width, height, radius float32) {
	if g == nil || g.nativeGraphics == nil || image == nil || image.nativeImage == nil || width <= 0 || height <= 0 {
		return
	}
	state := g.Save()
	defer g.Restore(state)

	path := NewPath(FillModeAlternate)
	defer path.Dispose()
	path.AddRoundedRectangle(x, y, width, height, radius)

	g.SetClipPath(path, CombineModeIntersect)
	g.DrawImageRect(image, x, y, width, height)
}

// ==================== Text Rendering (With Automatic Multi-Font & Emoji Fallback) ====================

// DrawString draws text at origin. Automatically falls back to system fonts (such as Segoe UI Emoji,
// Microsoft YaHei, SimSun) for emojis, special symbols, or characters missing from the primary font.
func (g *Graphics) DrawString(text string, font *Font, origin *PointF, brush *Brush) {
	if g == nil || g.nativeGraphics == nil || font == nil || font.nativeFont == nil || origin == nil || brush == nil || text == "" {
		return
	}
	if textNeedsFallback(text, font) {
		g.drawStringFallbackPoint(text, font, origin, brush)
		return
	}
	rect := &RectF{X: origin.X, Y: origin.Y, Width: 0, Height: 0}
	textUTF16, _ := syscall.UTF16FromString(text)
	GdipDrawString(g.nativeGraphics, &textUTF16[0], int32(len(textUTF16)-1), font.nativeFont, rect, nil, brush.nativeBrush)
}

// DrawStringEx draws formatted text inside a layout RectF. Automatically falls back to system fonts
// for emojis, symbols, or missing characters while preserving horizontal/vertical alignment and wrapping.
func (g *Graphics) DrawStringEx(text string, font *Font, rect *RectF, format *StringFormat, brush *Brush) {
	if g == nil || g.nativeGraphics == nil || font == nil || font.nativeFont == nil || rect == nil || brush == nil || text == "" {
		return
	}
	if textNeedsFallback(text, font) {
		g.drawStringFallbackRect(text, font, rect, format, brush)
		return
	}
	textUTF16, _ := syscall.UTF16FromString(text)
	var nativeFormat *GpStringFormat
	if format != nil {
		nativeFormat = format.nativeFormat
	}
	GdipDrawString(g.nativeGraphics, &textUTF16[0], int32(len(textUTF16)-1), font.nativeFont, rect, nativeFormat, brush.nativeBrush)
}

// DrawOutlinedString draws text with both an outer stroke (outlinePen) and an inner fill (fillBrush) using GraphicsPath.
func (g *Graphics) DrawOutlinedString(text string, familyName string, style int32, emSize float32, layoutRect *RectF, format *StringFormat, fillBrush *Brush, outlinePen *Pen) {
	if g == nil || g.nativeGraphics == nil || text == "" {
		return
	}
	path := NewPath(FillModeWinding)
	defer path.Dispose()
	path.AddString(text, familyName, style, emSize, layoutRect, format)
	if outlinePen != nil {
		g.DrawPath(outlinePen, path)
	}
	if fillBrush != nil {
		g.FillPath(fillBrush, path)
	}
}

func (g *Graphics) MeasureStringEx(text string, font *Font, layoutRect *RectF,
	stringFormat *StringFormat, boundingBox *RectF,
	codepointsFitted *int32, linesFilled *int32) {
	if g == nil || g.nativeGraphics == nil || font == nil || font.nativeFont == nil || boundingBox == nil {
		return
	}
	if text == "" {
		boundingBox.Width = 0
		boundingBox.Height = 0
		return
	}
	if textNeedsFallback(text, font) {
		g.measureStringFallbackRect(text, font, layoutRect, stringFormat, boundingBox, codepointsFitted, linesFilled)
		return
	}
	textUTF16, _ := syscall.UTF16FromString(text)
	var nativeFormat *GpStringFormat
	if stringFormat != nil {
		nativeFormat = stringFormat.nativeFormat
	}
	GdipMeasureString(g.nativeGraphics, &textUTF16[0], int32(len(textUTF16)-1), font.nativeFont, layoutRect, nativeFormat, boundingBox, codepointsFitted, linesFilled)
}

func (g *Graphics) MeasureString(text string, font *Font, layoutRect *RectF,
	stringFormat *StringFormat, boundingBox *RectF) {
	g.MeasureStringEx(text, font, layoutRect, stringFormat, boundingBox, nil, nil)
}

func (g *Graphics) MeasureCharacterRanges(text string, font *Font,
	layoutRect *RectF,
	stringFormat *StringFormat, regionCount int32,
	regions **GpRegion) {
	if g == nil || g.nativeGraphics == nil || font == nil || font.nativeFont == nil {
		return
	}
	textUTF16, _ := syscall.UTF16FromString(text)
	var nativeFormat *GpStringFormat
	if stringFormat != nil {
		nativeFormat = stringFormat.nativeFormat
	}
	GdipMeasureCharacterRanges(g.nativeGraphics, &textUTF16[0], int32(len(textUTF16)-1), font.nativeFont, layoutRect, nativeFormat, regionCount, regions)
}

func (g *Graphics) Dispose() {
	if g != nil && g.nativeGraphics != nil {
		GdipDeleteGraphics(g.nativeGraphics)
		g.nativeGraphics = nil
	}
}
