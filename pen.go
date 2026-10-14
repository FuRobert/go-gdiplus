package gdiplus

type Pen struct {
	nativePen *GpPen
}

func NewPen(color *Color, width float32) *Pen {
	p := &Pen{}
	var argb ARGB
	if color != nil {
		argb = color.GetValue()
	}
	GdipCreatePen1(argb, width, UnitWorld, &p.nativePen)
	return p
}

func NewPenFromBrush(brush *Brush, width float32) *Pen {
	p := &Pen{}
	if brush != nil && brush.nativeBrush != nil {
		GdipCreatePen2(brush.nativeBrush, width, UnitWorld, &p.nativePen)
	}
	return p
}

func (p *Pen) Dispose() {
	if p != nil && p.nativePen != nil {
		GdipDeletePen(p.nativePen)
		p.nativePen = nil
	}
}

func (p *Pen) Clone() *Pen {
	if p == nil || p.nativePen == nil {
		return &Pen{}
	}
	clone := &Pen{}
	GdipClonePen(p.nativePen, &clone.nativePen)
	return clone
}

func (p *Pen) SetWidth(width float32) {
	if p != nil && p.nativePen != nil {
		GdipSetPenWidth(p.nativePen, width)
	}
}

func (p *Pen) GetWidth() (width float32) {
	if p != nil && p.nativePen != nil {
		GdipGetPenWidth(p.nativePen, &width)
	}
	return
}

func (p *Pen) SetLineCap(startCap, endCap LineCap, dashCap DashCap) {
	if p != nil && p.nativePen != nil {
		GdipSetPenLineCap197819(p.nativePen, GpLineCap(startCap), GpLineCap(endCap), GpDashCap(dashCap))
	}
}

func (p *Pen) SetStartCap(startCap LineCap) {
	if p != nil && p.nativePen != nil {
		GdipSetPenStartCap(p.nativePen, GpLineCap(startCap))
	}
}

func (p *Pen) SetEndCap(endCap LineCap) {
	if p != nil && p.nativePen != nil {
		GdipSetPenEndCap(p.nativePen, GpLineCap(endCap))
	}
}

func (p *Pen) SetDashCap(dashCap DashCap) {
	if p != nil && p.nativePen != nil {
		GdipSetPenDashCap197819(p.nativePen, GpDashCap(dashCap))
	}
}

func (p *Pen) GetStartCap() (startCap LineCap) {
	if p != nil && p.nativePen != nil {
		GdipGetPenStartCap(p.nativePen, (*GpLineCap)(&startCap))
	}
	return
}

func (p *Pen) GetEndCap() (endCap LineCap) {
	if p != nil && p.nativePen != nil {
		GdipGetPenEndCap(p.nativePen, (*GpLineCap)(&endCap))
	}
	return
}

func (p *Pen) GetDashCap() (dashCap DashCap) {
	if p != nil && p.nativePen != nil {
		GdipGetPenDashCap197819(p.nativePen, (*GpDashCap)(&dashCap))
	}
	return
}

func (p *Pen) SetLineJoin(lineJoin LineJoin) {
	if p != nil && p.nativePen != nil {
		GdipSetPenLineJoin(p.nativePen, GpLineJoin(lineJoin))
	}
}

func (p *Pen) GetLineJoin() (lineJoin LineJoin) {
	if p != nil && p.nativePen != nil {
		GdipGetPenLineJoin(p.nativePen, (*GpLineJoin)(&lineJoin))
	}
	return
}

func (p *Pen) SetCustomStartCap(customCap *GpCustomLineCap) {
	if p != nil && p.nativePen != nil {
		GdipSetPenCustomStartCap(p.nativePen, customCap)
	}
}

func (p *Pen) GetCustomStartCap() (customCap *GpCustomLineCap) {
	if p != nil && p.nativePen != nil {
		GdipGetPenCustomStartCap(p.nativePen, &customCap)
	}
	return
}

func (p *Pen) SetCustomEndCap(customCap *GpCustomLineCap) {
	if p != nil && p.nativePen != nil {
		GdipSetPenCustomEndCap(p.nativePen, customCap)
	}
}

func (p *Pen) GetCustomEndCap() (customCap *GpCustomLineCap) {
	if p != nil && p.nativePen != nil {
		GdipGetPenCustomEndCap(p.nativePen, &customCap)
	}
	return
}

func (p *Pen) SetMiterLimit(miterLimit float32) {
	if p != nil && p.nativePen != nil {
		GdipSetPenMiterLimit(p.nativePen, miterLimit)
	}
}

func (p *Pen) GetMiterLimit() (miterLimit float32) {
	if p != nil && p.nativePen != nil {
		GdipGetPenMiterLimit(p.nativePen, &miterLimit)
	}
	return
}

func (p *Pen) SetMode(penMode PenAlignment) {
	if p != nil && p.nativePen != nil {
		GdipSetPenMode(p.nativePen, GpPenAlignment(penMode))
	}
}

func (p *Pen) GetMode() (penMode PenAlignment) {
	if p != nil && p.nativePen != nil {
		GdipGetPenMode(p.nativePen, (*GpPenAlignment)(&penMode))
	}
	return
}

func (p *Pen) SetTransform(matrix *GpMatrix) {
	if p != nil && p.nativePen != nil {
		GdipSetPenTransform(p.nativePen, matrix)
	}
}

func (p *Pen) GetTransform(matrix *GpMatrix) {
	if p != nil && p.nativePen != nil {
		GdipGetPenTransform(p.nativePen, matrix)
	}
}

func (p *Pen) ResetTransform() {
	if p != nil && p.nativePen != nil {
		GdipResetPenTransform(p.nativePen)
	}
}

func (p *Pen) MultiplyTransform(matrix *GpMatrix, order MatrixOrder) {
	if p != nil && p.nativePen != nil {
		GdipMultiplyPenTransform(p.nativePen, matrix, GpMatrixOrder(order))
	}
}

func (p *Pen) TranslateTransform(dx, dy float32, order MatrixOrder) {
	if p != nil && p.nativePen != nil {
		GdipTranslatePenTransform(p.nativePen, dx, dy, GpMatrixOrder(order))
	}
}

func (p *Pen) ScaleTransform(sx, sy float32, order MatrixOrder) {
	if p != nil && p.nativePen != nil {
		GdipScalePenTransform(p.nativePen, sx, sy, GpMatrixOrder(order))
	}
}

func (p *Pen) RotateTransform(angle float32, order MatrixOrder) {
	if p != nil && p.nativePen != nil {
		GdipRotatePenTransform(p.nativePen, angle, GpMatrixOrder(order))
	}
}

func (p *Pen) SetColor(color *Color) {
	if p != nil && p.nativePen != nil && color != nil {
		GdipSetPenColor(p.nativePen, color.GetValue())
	}
}

func (p *Pen) GetColor() (color Color) {
	if p != nil && p.nativePen != nil {
		GdipGetPenColor(p.nativePen, &color.Argb)
	}
	return
}

func (p *Pen) SetBrush(brush *Brush) {
	if p != nil && p.nativePen != nil && brush != nil {
		GdipSetPenBrushFill(p.nativePen, brush.nativeBrush)
	}
}

func (p *Pen) GetBrush() *Brush {
	brush := &Brush{}
	if p != nil && p.nativePen != nil {
		GdipGetPenBrushFill(p.nativePen, &brush.nativeBrush)
	}
	return brush
}

func (p *Pen) GetPenType() (penType PenType) {
	if p != nil && p.nativePen != nil {
		GdipGetPenFillType(p.nativePen, (*GpPenType)(&penType))
	}
	return
}

func (p *Pen) GetDashStyle() (dashStyle DashStyle) {
	if p != nil && p.nativePen != nil {
		GdipGetPenDashStyle(p.nativePen, (*GpDashStyle)(&dashStyle))
	}
	return
}

func (p *Pen) SetDashStyle(dashStyle DashStyle) {
	if p != nil && p.nativePen != nil {
		GdipSetPenDashStyle(p.nativePen, GpDashStyle(dashStyle))
	}
}

func (p *Pen) GetDashOffset() (offset float32) {
	if p != nil && p.nativePen != nil {
		GdipGetPenDashOffset(p.nativePen, &offset)
	}
	return
}

func (p *Pen) SetDashOffset(offset float32) {
	if p != nil && p.nativePen != nil {
		GdipSetPenDashOffset(p.nativePen, offset)
	}
}

func (p *Pen) GetDashCount() (count int32) {
	if p != nil && p.nativePen != nil {
		GdipGetPenDashCount(p.nativePen, &count)
	}
	return
}

func (p *Pen) SetDashArray(dash []float32) {
	if p != nil && p.nativePen != nil && len(dash) > 0 {
		GdipSetPenDashArray(p.nativePen, &dash[0], int32(len(dash)))
	}
}

func (p *Pen) GetDashArray(dash *float32, count int32) {
	if p != nil && p.nativePen != nil {
		GdipGetPenDashArray(p.nativePen, dash, count)
	}
}

func (p *Pen) GetCompoundCount() (count int32) {
	if p != nil && p.nativePen != nil {
		GdipGetPenCompoundCount(p.nativePen, &count)
	}
	return
}

func (p *Pen) SetCompoundArray(dash []float32) {
	if p != nil && p.nativePen != nil && len(dash) > 0 {
		GdipSetPenCompoundArray(p.nativePen, &dash[0], int32(len(dash)))
	}
}

func (p *Pen) GetCompoundArray(dash *float32, count int32) {
	if p != nil && p.nativePen != nil {
		GdipGetPenCompoundArray(p.nativePen, dash, count)
	}
}
