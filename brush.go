package gdiplus

import "fmt"

type Brush struct {
	nativeBrush *GpBrush
}

type SolidBrush struct {
	Brush
}

type LinearGradientBrush struct {
	Brush
}

type HatchBrush struct {
	Brush
}

type TextureBrush struct {
	Brush
}

type PathGradientBrush struct {
	Brush
}

func (b *Brush) Dispose() {
	if b != nil && b.nativeBrush != nil {
		GdipDeleteBrush(b.nativeBrush)
		b.nativeBrush = nil
	}
}

func (b *Brush) GetBrushType() (brushType BrushType) {
	if b == nil || b.nativeBrush == nil {
		return
	}
	GdipGetBrushType(b.nativeBrush, (*GpBrushType)(&brushType))
	return
}

func (b *Brush) Clone() *Brush {
	if b == nil || b.nativeBrush == nil {
		return &Brush{}
	}
	clone := &Brush{}
	GdipCloneBrush(b.nativeBrush, &clone.nativeBrush)
	return clone
}

// ==================== SolidBrush ====================

func NewSolidBrush(color *Color) *SolidBrush {
	b := &SolidBrush{}
	var solidFill *GpSolidFill
	var argb ARGB
	if color != nil {
		argb = color.GetValue()
	}
	GdipCreateSolidFill(argb, &solidFill)
	if solidFill != nil {
		b.nativeBrush = &solidFill.GpBrush
	}
	return b
}

func (b *SolidBrush) AsBrush() *Brush {
	if b == nil {
		return nil
	}
	return &b.Brush
}

func (b *SolidBrush) SetColor(color *Color) {
	if b == nil || b.nativeBrush == nil || color == nil {
		return
	}
	GdipSetSolidFillColor(b.nativeBrush, color.GetValue())
}

func (b *SolidBrush) GetColor() (color Color) {
	if b == nil || b.nativeBrush == nil {
		return
	}
	GdipGetSolidFillColor(b.nativeBrush, &color.Argb)
	return
}

// ==================== LinearGradientBrush ====================

// NewLinearGradientBrush creates a linear gradient brush between two floating-point points.
func NewLinearGradientBrush(point1, point2 *PointF, color1, color2 *Color) *LinearGradientBrush {
	return NewLinearGradientBrushEx(point1, point2, color1, color2, WrapModeTileFlipXY)
}

// NewLinearGradientBrushEx creates a linear gradient brush between two points with an explicit WrapMode.
func NewLinearGradientBrushEx(point1, point2 *PointF, color1, color2 *Color, wrapMode WrapMode) *LinearGradientBrush {
	b := &LinearGradientBrush{}
	var lineGrad *GpLineGradient
	var c1, c2 ARGB
	if color1 != nil {
		c1 = color1.GetValue()
	}
	if color2 != nil {
		c2 = color2.GetValue()
	}
	GdipCreateLineBrush(point1, point2, c1, c2, wrapMode, &lineGrad)
	b.nativeBrush = (*GpBrush)(lineGrad)
	return b
}

// NewLinearGradientBrushI creates a linear gradient brush between two integer points.
func NewLinearGradientBrushI(point1, point2 *Point, color1, color2 *Color) *LinearGradientBrush {
	b := &LinearGradientBrush{}
	var lineGrad *GpLineGradient
	var c1, c2 ARGB
	if color1 != nil {
		c1 = color1.GetValue()
	}
	if color2 != nil {
		c2 = color2.GetValue()
	}
	GdipCreateLineBrushI(point1, point2, c1, c2, WrapModeTileFlipXY, &lineGrad)
	b.nativeBrush = (*GpBrush)(lineGrad)
	return b
}

// NewLinearGradientBrushFromRect creates a linear gradient brush from a RectF and a LinearGradientMode
// (LinearGradientModeHorizontal, LinearGradientModeVertical, LinearGradientModeForwardDiagonal, LinearGradientModeBackwardDiagonal).
func NewLinearGradientBrushFromRect(rect *RectF, color1, color2 *Color, mode LinearGradientMode) *LinearGradientBrush {
	b := &LinearGradientBrush{}
	var lineGrad *GpLineGradient
	var c1, c2 ARGB
	if color1 != nil {
		c1 = color1.GetValue()
	}
	if color2 != nil {
		c2 = color2.GetValue()
	}
	GdipCreateLineBrushFromRect(rect, c1, c2, mode, WrapModeTileFlipXY, &lineGrad)
	b.nativeBrush = (*GpBrush)(lineGrad)
	return b
}

// NewLinearGradientBrushFromRectI creates a linear gradient brush from an integer Rect and a LinearGradientMode.
func NewLinearGradientBrushFromRectI(rect *Rect, color1, color2 *Color, mode LinearGradientMode) *LinearGradientBrush {
	b := &LinearGradientBrush{}
	var lineGrad *GpLineGradient
	var c1, c2 ARGB
	if color1 != nil {
		c1 = color1.GetValue()
	}
	if color2 != nil {
		c2 = color2.GetValue()
	}
	GdipCreateLineBrushFromRectI(rect, c1, c2, mode, WrapModeTileFlipXY, &lineGrad)
	b.nativeBrush = (*GpBrush)(lineGrad)
	return b
}

// NewLinearGradientBrushWithAngle creates a linear gradient brush from a RectF at any custom angle (in degrees).
func NewLinearGradientBrushWithAngle(rect *RectF, color1, color2 *Color, angle float32, isAngleScalable bool) *LinearGradientBrush {
	b := &LinearGradientBrush{}
	var lineGrad *GpLineGradient
	var c1, c2 ARGB
	if color1 != nil {
		c1 = color1.GetValue()
	}
	if color2 != nil {
		c2 = color2.GetValue()
	}
	var scalable int32
	if isAngleScalable {
		scalable = 1
	}
	GdipCreateLineBrushFromRectWithAngle(rect, c1, c2, angle, scalable, WrapModeTileFlipXY, &lineGrad)
	b.nativeBrush = (*GpBrush)(lineGrad)
	return b
}

func (b *LinearGradientBrush) AsBrush() *Brush {
	if b == nil {
		return nil
	}
	return &b.Brush
}

func (b *LinearGradientBrush) SetLinearColors(color1, color2 *Color) {
	if b == nil || b.nativeBrush == nil || color1 == nil || color2 == nil {
		return
	}
	GdipSetLineColors((*GpLineGradient)(b.nativeBrush), color1.GetValue(), color2.GetValue())
}

// SetInterpolationColors configures a multi-stop color gradient.
// colors and positions must have the same length (>= 2), with positions[0] == 0.0 and positions[last] == 1.0.
func (b *LinearGradientBrush) SetInterpolationColors(colors []*Color, positions []float32) error {
	if b == nil || b.nativeBrush == nil {
		return fmt.Errorf("LinearGradientBrush is nil")
	}
	if len(colors) < 2 || len(colors) != len(positions) {
		return fmt.Errorf("colors and positions must have matching length >= 2")
	}
	argbs := make([]ARGB, len(colors))
	for i, c := range colors {
		if c != nil {
			argbs[i] = c.GetValue()
		}
	}
	status := GdipSetLinePresetBlend((*GpLineGradient)(b.nativeBrush), &argbs[0], &positions[0], int32(len(colors)))
	if status != Ok {
		return fmt.Errorf("GdipSetLinePresetBlend failed: %s", status.String())
	}
	return nil
}

func (b *LinearGradientBrush) SetBlendBellShape(focus, scale float32) {
	if b == nil || b.nativeBrush == nil {
		return
	}
	GdipSetLineSigmaBlend((*GpLineGradient)(b.nativeBrush), focus, scale)
}

func (b *LinearGradientBrush) SetBlendTriangularShape(focus, scale float32) {
	if b == nil || b.nativeBrush == nil {
		return
	}
	GdipSetLineLinearBlend((*GpLineGradient)(b.nativeBrush), focus, scale)
}

func (b *LinearGradientBrush) SetWrapMode(wrapMode WrapMode) {
	if b == nil || b.nativeBrush == nil {
		return
	}
	GdipSetLineWrapMode((*GpLineGradient)(b.nativeBrush), wrapMode)
}

func (b *LinearGradientBrush) SetGammaCorrection(useGammaCorrection bool) {
	if b == nil || b.nativeBrush == nil {
		return
	}
	var flag int32
	if useGammaCorrection {
		flag = 1
	}
	GdipSetLineGammaCorrection((*GpLineGradient)(b.nativeBrush), flag)
}

// ==================== HatchBrush ====================

func NewHatchBrush(style HatchStyle, foreColor, backColor *Color) *HatchBrush {
	b := &HatchBrush{}
	var hatch *GpHatch
	var fg, bg ARGB
	if foreColor != nil {
		fg = foreColor.GetValue()
	}
	if backColor != nil {
		bg = backColor.GetValue()
	}
	GdipCreateHatchBrush(style, fg, bg, &hatch)
	b.nativeBrush = (*GpBrush)(hatch)
	return b
}

func (b *HatchBrush) AsBrush() *Brush {
	if b == nil {
		return nil
	}
	return &b.Brush
}

// ==================== TextureBrush ====================

func NewTextureBrush(image *Image, wrapMode WrapMode) *TextureBrush {
	b := &TextureBrush{}
	if image == nil || image.nativeImage == nil {
		return b
	}
	var tex *GpTexture
	GdipCreateTexture(image.nativeImage, wrapMode, &tex)
	b.nativeBrush = (*GpBrush)(tex)
	return b
}

func (b *TextureBrush) AsBrush() *Brush {
	if b == nil {
		return nil
	}
	return &b.Brush
}

// ==================== PathGradientBrush ====================

func NewPathGradientBrush(path *GraphicsPath) *PathGradientBrush {
	b := &PathGradientBrush{}
	if path == nil || path.nativePath == nil {
		return b
	}
	var pg *GpPathGradient
	GdipCreatePathGradientFromPath(path.nativePath, &pg)
	b.nativeBrush = (*GpBrush)(pg)
	return b
}

func (b *PathGradientBrush) AsBrush() *Brush {
	if b == nil {
		return nil
	}
	return &b.Brush
}

func (b *PathGradientBrush) SetCenterColor(color *Color) {
	if b == nil || b.nativeBrush == nil || color == nil {
		return
	}
	GdipSetPathGradientCenterColor((*GpPathGradient)(b.nativeBrush), color.GetValue())
}

func (b *PathGradientBrush) SetSurroundColors(colors []*Color) {
	if b == nil || b.nativeBrush == nil || len(colors) == 0 {
		return
	}
	argbs := make([]ARGB, len(colors))
	for i, c := range colors {
		if c != nil {
			argbs[i] = c.GetValue()
		}
	}
	count := int32(len(argbs))
	GdipSetPathGradientSurroundColorsWithCount((*GpPathGradient)(b.nativeBrush), &argbs[0], &count)
}

func (b *PathGradientBrush) SetCenterPoint(pt *PointF) {
	if b == nil || b.nativeBrush == nil || pt == nil {
		return
	}
	GdipSetPathGradientCenterPoint((*GpPathGradient)(b.nativeBrush), pt)
}
