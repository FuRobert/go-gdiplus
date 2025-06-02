package gdiplus

import (
	"github.com/go-ole/go-ole"
)

const LANG_NEUTRAL = 0x00

// Status
const (
	Ok                        GpStatus = 0
	GenericError              GpStatus = 1
	InvalidParameter          GpStatus = 2
	OutOfMemory               GpStatus = 3
	ObjectBusy                GpStatus = 4
	InsufficientBuffer        GpStatus = 5
	NotImplemented            GpStatus = 6
	Win32Error                GpStatus = 7
	WrongState                GpStatus = 8
	Aborted                   GpStatus = 9
	FileNotFound              GpStatus = 10
	ValueOverflow             GpStatus = 11
	AccessDenied              GpStatus = 12
	UnknownImageFormat        GpStatus = 13
	FontFamilyNotFound        GpStatus = 14
	FontStyleNotFound         GpStatus = 15
	NotTrueTypeFont           GpStatus = 16
	UnsupportedGdiplusVersion GpStatus = 17
	GdiplusNotInitialized     GpStatus = 18
	PropertyNotFound          GpStatus = 19
	PropertyNotSupported      GpStatus = 20
	ProfileNotFound           GpStatus = 21
)

// Unit
const (
	UnitWorld      = 0 // 0 -- World coordinate (non-physical unit)
	UnitDisplay    = 1 // 1 -- Variable -- for PageTransform only
	UnitPixel      = 2 // 2 -- Each unit is one device pixel.
	UnitPoint      = 3 // 3 -- Each unit is a printer's point, or 1/72 inch.
	UnitInch       = 4 // 4 -- Each unit is 1 inch.
	UnitDocument   = 5 // 5 -- Each unit is 1/300 inch.
	UnitMillimeter = 6 // 6 -- Each unit is 1 millimeter.
)

const (
	AlphaShift = 24
	RedShift   = 16
	GreenShift = 8
	BlueShift  = 0
)

const (
	AlphaMask = 0xff000000
	RedMask   = 0x00ff0000
	GreenMask = 0x0000ff00
	BlueMask  = 0x000000ff
)

// FontStyle
const (
	FontStyleRegular    = 0
	FontStyleBold       = 1
	FontStyleItalic     = 2
	FontStyleBoldItalic = 3
	FontStyleUnderline  = 4
	FontStyleStrikeout  = 8
)

// QualityMode
const (
	QualityModeInvalid = iota - 1
	QualityModeDefault
	QualityModeLow  // Best performance
	QualityModeHigh // Best rendering quality
)

// Alpha Compositing mode
const (
	CompositingModeSourceOver = iota // 0
	CompositingModeSourceCopy        // 1
)

// Alpha Compositing quality
const (
	CompositingQualityInvalid = iota + QualityModeInvalid
	CompositingQualityDefault
	CompositingQualityHighSpeed
	CompositingQualityHighQuality
	CompositingQualityGammaCorrected
	CompositingQualityAssumeLinear
)

// InterpolationMode
const (
	InterpolationModeInvalid = iota + QualityModeInvalid
	InterpolationModeDefault
	InterpolationModeLowQuality
	InterpolationModeHighQuality
	InterpolationModeBilinear
	InterpolationModeBicubic
	InterpolationModeNearestNeighbor
	InterpolationModeHighQualityBilinear
	InterpolationModeHighQualityBicubic
)

// SmoothingMode
const (
	SmoothingModeInvalid = iota + QualityModeInvalid
	SmoothingModeDefault
	SmoothingModeHighSpeed
	SmoothingModeHighQuality
	SmoothingModeNone
	SmoothingModeAntiAlias

/*
#if (GDIPVER >= 0x0110)

	SmoothingModeAntiAlias8x4 = SmoothingModeAntiAlias,
	SmoothingModeAntiAlias8x8

#endif //(GDIPVER >= 0x0110)
*/
)

// Pixel Format Mode
const (
	PixelOffsetModeInvalid = iota + QualityModeInvalid
	PixelOffsetModeDefault
	PixelOffsetModeHighSpeed
	PixelOffsetModeHighQuality
	PixelOffsetModeNone // No pixel offset
	PixelOffsetModeHalf // Offset by -0.5, -0.5 for fast anti-alias perf
)

// Text Rendering Hint
const (
	TextRenderingHintSystemDefault            = iota // Glyph with system default rendering hint
	TextRenderingHintSingleBitPerPixelGridFit        // Glyph bitmap with hinting
	TextRenderingHintSingleBitPerPixel               // Glyph bitmap without hinting
	TextRenderingHintAntiAliasGridFit                // Glyph anti-alias bitmap with hinting
	TextRenderingHintAntiAlias                       // Glyph anti-alias bitmap without hinting
	TextRenderingHintClearTypeGridFit                // Glyph CT bitmap with hinting
)

// Fill mode constants
const (
	FillModeAlternate = iota // 0
	FillModeWinding          // 1
)

// BrushType
const (
	BrushTypeSolidColor GpBrushType = iota
	BrushTypeHatchFill
	BrushTypeTextureFill
	BrushTypePathGradient
	BrushTypeLinearGradient
)

// LineCap
const (
	LineCapFlat GpLineCap = iota
	LineCapSquare
	LineCapRound
	LineCapTriangle
	LineCapNoAnchor
	LineCapSquareAnchor
	LineCapRoundAnchor
	LineCapDiamondAnchor
	LineCapArrowAnchor
	LineCapCustom
	LineCapAnchorMask
)

// LineJoin
const (
	LineJoinMiter GpLineJoin = iota
	LineJoinBevel
	LineJoinRound
	LineJoinMiterClipped
)

// DashCap
const (
	DashCapFlat GpDashCap = iota
	DashCapRound
	DashCapTriangle
)

// DashStyle
const (
	DashStyleSolid GpDashStyle = iota
	DashStyleDash
	DashStyleDot
	DashStyleDashDot
	DashStyleDashDotDot
	DashStyleCustom
)

// PenAlignment
const (
	PenAlignmentCenter GpPenAlignment = iota
	PenAlignmentInset
)

// MatrixOrder
const (
	MatrixOrderPrepend GpMatrixOrder = iota
	MatrixOrderAppend
)

// PenType
const (
	PenTypeSolidColor GpPenType = iota
	PenTypeHatchFill
	PenTypeTextureFill
	PenTypePathGradient
	PenTypeLinearGradient
	PenTypeUnknown
)

var (
	ClsidPNGEncoder = ole.GUID{Data1: 0x557CF406, Data2: 0x1A04, Data3: 0x11D3, Data4: [8]byte{0x9A, 0x73, 0x00, 0x00, 0xF8, 0x1E, 0xF3, 0x2E}}
	ClsidJPEGEncoder = ole.GUID{Data1: 0x557CF401, Data2: 0x1A04, Data3: 0x11D3, Data4: [8]byte{0x9A, 0x73, 0x00, 0x00, 0xF8, 0x1E, 0xF3, 0x2E}}
)

var (
	EncoderQuality = ole.GUID{Data1: 0x1D5BE4B5, Data2: 0xFA4A, Data3: 0x452D, Data4: [8]byte{0x9C, 0xDD, 0x5D, 0xB3, 0x51, 0x05, 0xE7, 0xEB}}
)

const (
	EncoderParameterValueTypeLong uint32 = 4 
)