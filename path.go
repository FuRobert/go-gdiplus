package gdiplus

import "syscall"

type GraphicsPath struct {
	nativePath *GpPath
}

func NewPath(fillMode int32) *GraphicsPath {
	p := &GraphicsPath{}
	GdipCreatePath(fillMode, &p.nativePath)
	return p
}

func (p *GraphicsPath) Reset() {
	if p != nil && p.nativePath != nil {
		GdipResetPath(p.nativePath)
	}
}

func (p *GraphicsPath) StartFigure() {
	if p != nil && p.nativePath != nil {
		GdipStartPathFigure(p.nativePath)
	}
}

func (p *GraphicsPath) AddArcRect(rect *Rect, startAngle, sweepAngle float32) {
	if p == nil || p.nativePath == nil || rect == nil {
		return
	}
	GdipAddPathArcI(p.nativePath, rect.X, rect.Y, rect.Width, rect.Height, startAngle, sweepAngle)
}

func (p *GraphicsPath) AddArcRectF(rect *RectF, startAngle, sweepAngle float32) {
	if p == nil || p.nativePath == nil || rect == nil {
		return
	}
	GdipAddPathArc(p.nativePath, rect.X, rect.Y, rect.Width, rect.Height, startAngle, sweepAngle)
}

func (p *GraphicsPath) AddArc(x, y, width, height, startAngle, sweepAngle float32) {
	if p == nil || p.nativePath == nil {
		return
	}
	GdipAddPathArc(p.nativePath, x, y, width, height, startAngle, sweepAngle)
}

func (p *GraphicsPath) AddArcI(x, y, width, height int32, startAngle, sweepAngle float32) {
	if p == nil || p.nativePath == nil {
		return
	}
	GdipAddPathArcI(p.nativePath, x, y, width, height, startAngle, sweepAngle)
}

func (p *GraphicsPath) AddLine(x1, y1, x2, y2 float32) {
	if p == nil || p.nativePath == nil {
		return
	}
	GdipAddPathLine(p.nativePath, x1, y1, x2, y2)
}

func (p *GraphicsPath) AddLineI(x1, y1, x2, y2 int32) {
	if p == nil || p.nativePath == nil {
		return
	}
	GdipAddPathLineI(p.nativePath, x1, y1, x2, y2)
}

func (p *GraphicsPath) AddRectangle(x, y, width, height float32) {
	if p == nil || p.nativePath == nil {
		return
	}
	GdipAddPathRectangle(p.nativePath, x, y, width, height)
}

func (p *GraphicsPath) AddRectangleI(x, y, width, height int32) {
	if p == nil || p.nativePath == nil {
		return
	}
	GdipAddPathRectangleI(p.nativePath, x, y, width, height)
}

func (p *GraphicsPath) AddEllipse(x, y, width, height float32) {
	if p == nil || p.nativePath == nil {
		return
	}
	GdipAddPathEllipse(p.nativePath, x, y, width, height)
}

func (p *GraphicsPath) AddEllipseI(x, y, width, height int32) {
	if p == nil || p.nativePath == nil {
		return
	}
	GdipAddPathEllipseI(p.nativePath, x, y, width, height)
}

func (p *GraphicsPath) AddPie(x, y, width, height, startAngle, sweepAngle float32) {
	if p == nil || p.nativePath == nil {
		return
	}
	GdipAddPathPie(p.nativePath, x, y, width, height, startAngle, sweepAngle)
}

func (p *GraphicsPath) AddPolygon(points []PointF) {
	if p == nil || p.nativePath == nil || len(points) == 0 {
		return
	}
	GdipAddPathPolygon(p.nativePath, &points[0], int32(len(points)))
}

func (p *GraphicsPath) AddBezier(x1, y1, x2, y2, x3, y3, x4, y4 float32) {
	if p == nil || p.nativePath == nil {
		return
	}
	GdipAddPathBezier(p.nativePath, x1, y1, x2, y2, x3, y3, x4, y4)
}

// AddRoundedRectangle adds a rounded rectangle figure with uniform corner radius to the path.
func (p *GraphicsPath) AddRoundedRectangle(x, y, width, height, radius float32) {
	if p == nil || p.nativePath == nil {
		return
	}
	if radius <= 0 {
		p.AddRectangle(x, y, width, height)
		return
	}
	maxRadius := width * 0.5
	if height*0.5 < maxRadius {
		maxRadius = height * 0.5
	}
	if radius > maxRadius {
		radius = maxRadius
	}
	d := radius * 2
	p.StartFigure()
	p.AddArc(x, y, d, d, 180, 90)
	p.AddArc(x+width-d, y, d, d, 270, 90)
	p.AddArc(x+width-d, y+height-d, d, d, 0, 90)
	p.AddArc(x, y+height-d, d, d, 90, 90)
	p.CloseFigure()
}

// AddString adds a text outline to the path (useful for drawing stroked/outlined text).
func (p *GraphicsPath) AddString(text string, familyName string, style int32, emSize float32, layoutRect *RectF, format *StringFormat) {
	if p == nil || p.nativePath == nil || text == "" {
		return
	}
	fm := NewFontFamily(familyName, nil)
	defer fm.Dispose()
	if fm.nativeFontFamily == nil {
		return
	}
	textUTF16, _ := syscall.UTF16FromString(text)
	var nativeFormat *GpStringFormat
	if format != nil {
		nativeFormat = format.nativeFormat
	}
	rect := layoutRect
	if rect == nil {
		rect = &RectF{}
	}
	GdipAddPathString(p.nativePath, &textUTF16[0], int32(len(textUTF16)-1), fm.nativeFontFamily, style, emSize, rect, nativeFormat)
}

func (p *GraphicsPath) CloseAllFigures() {
	if p != nil && p.nativePath != nil {
		GdipClosePathFigures(p.nativePath)
	}
}

func (p *GraphicsPath) CloseFigure() {
	if p != nil && p.nativePath != nil {
		GdipClosePathFigure(p.nativePath)
	}
}

func (p *GraphicsPath) Dispose() {
	if p != nil && p.nativePath != nil {
		GdipDeletePath(p.nativePath)
		p.nativePath = nil
	}
}
