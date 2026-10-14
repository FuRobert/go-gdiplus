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
)

// Pixel Offset Mode
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

// LinearGradientMode specifies the direction of a linear gradient
const (
	LinearGradientModeHorizontal       LinearGradientMode = 0 // Left to right
	LinearGradientModeVertical         LinearGradientMode = 1 // Top to bottom
	LinearGradientModeForwardDiagonal  LinearGradientMode = 2 // Upper-left to lower-right
	LinearGradientModeBackwardDiagonal LinearGradientMode = 3 // Upper-right to lower-left
)

// WrapMode specifies how a texture or gradient is tiled
const (
	WrapModeTile       WrapMode = 0
	WrapModeTileFlipX  WrapMode = 1
	WrapModeTileFlipY  WrapMode = 2
	WrapModeTileFlipXY WrapMode = 3
	WrapModeClamp      WrapMode = 4
)

// HatchStyle specifies the pattern used by a HatchBrush
const (
	HatchStyleHorizontal             HatchStyle = 0
	HatchStyleVertical               HatchStyle = 1
	HatchStyleForwardDiagonal        HatchStyle = 2
	HatchStyleBackwardDiagonal       HatchStyle = 3
	HatchStyleCross                  HatchStyle = 4
	HatchStyleDiagonalCross          HatchStyle = 5
	HatchStyle05Percent              HatchStyle = 6
	HatchStyle10Percent              HatchStyle = 7
	HatchStyle20Percent              HatchStyle = 8
	HatchStyle25Percent              HatchStyle = 9
	HatchStyle30Percent              HatchStyle = 10
	HatchStyle40Percent              HatchStyle = 11
	HatchStyle50Percent              HatchStyle = 12
	HatchStyle60Percent              HatchStyle = 13
	HatchStyle70Percent              HatchStyle = 14
	HatchStyle75Percent              HatchStyle = 15
	HatchStyle80Percent              HatchStyle = 16
	HatchStyle90Percent              HatchStyle = 17
	HatchStyleLightDownwardDiagonal  HatchStyle = 18
	HatchStyleLightUpwardDiagonal    HatchStyle = 19
	HatchStyleDarkDownwardDiagonal   HatchStyle = 20
	HatchStyleDarkUpwardDiagonal     HatchStyle = 21
	HatchStyleWideDownwardDiagonal   HatchStyle = 22
	HatchStyleWideUpwardDiagonal     HatchStyle = 23
	HatchStyleLightVertical          HatchStyle = 24
	HatchStyleLightHorizontal        HatchStyle = 25
	HatchStyleNarrowVertical         HatchStyle = 26
	HatchStyleNarrowHorizontal       HatchStyle = 27
	HatchStyleDarkVertical           HatchStyle = 28
	HatchStyleDarkHorizontal         HatchStyle = 29
	HatchStyleDashedDownwardDiagonal HatchStyle = 30
	HatchStyleDashedUpwardDiagonal   HatchStyle = 31
	HatchStyleDashedHorizontal       HatchStyle = 32
	HatchStyleDashedVertical         HatchStyle = 33
	HatchStyleSmallConfetti          HatchStyle = 34
	HatchStyleLargeConfetti          HatchStyle = 35
	HatchStyleZigZag                 HatchStyle = 36
	HatchStyleWave                   HatchStyle = 37
	HatchStyleDiagonalBrick          HatchStyle = 38
	HatchStyleHorizontalBrick        HatchStyle = 39
	HatchStyleWeave                  HatchStyle = 40
	HatchStylePlaid                  HatchStyle = 41
	HatchStyleDivot                  HatchStyle = 42
	HatchStyleDottedGrid             HatchStyle = 43
	HatchStyleDottedDiamond          HatchStyle = 44
	HatchStyleShingle                HatchStyle = 45
	HatchStyleTrellis                HatchStyle = 46
	HatchStyleSphere                 HatchStyle = 47
	HatchStyleSmallGrid              HatchStyle = 48
	HatchStyleSmallCheckerBoard      HatchStyle = 49
	HatchStyleLargeCheckerBoard      HatchStyle = 50
	HatchStyleOutlinedDiamond        HatchStyle = 51
	HatchStyleSolidDiamond           HatchStyle = 52
)

// StringAlignment specifies text alignment (horizontal or vertical)
const (
	StringAlignmentNear   StringAlignment = 0 // Left (or Top for LineAlignment)
	StringAlignmentCenter StringAlignment = 1 // Center (or Middle for LineAlignment)
	StringAlignmentFar    StringAlignment = 2 // Right (or Bottom for LineAlignment)
)

// StringTrimming specifies how characters are trimmed from a string that does not fit
const (
	StringTrimmingNone              StringTrimming = 0
	StringTrimmingCharacter         StringTrimming = 1
	StringTrimmingWord              StringTrimming = 2
	StringTrimmingEllipsisCharacter StringTrimming = 3
	StringTrimmingEllipsisWord      StringTrimming = 4
	StringTrimmingEllipsisPath      StringTrimming = 5
)

// StringFormatFlags specifies text layout and display options
const (
	StringFormatFlagsDirectionRightToLeft  int32 = 0x00000001
	StringFormatFlagsDirectionVertical     int32 = 0x00000002
	StringFormatFlagsNoFitBlackBox         int32 = 0x00000004
	StringFormatFlagsDisplayFormatControl  int32 = 0x00000020
	StringFormatFlagsNoFontFallback        int32 = 0x00000400
	StringFormatFlagsMeasureTrailingSpaces int32 = 0x00000800
	StringFormatFlagsNoWrap                int32 = 0x00001000
	StringFormatFlagsLineLimit             int32 = 0x00002000
	StringFormatFlagsNoClip                int32 = 0x00004000
)

// CombineMode specifies how a new clipping region is combined with an existing one
const (
	CombineModeReplace    CombineMode = 0
	CombineModeIntersect  CombineMode = 1
	CombineModeUnion      CombineMode = 2
	CombineModeXor        CombineMode = 3
	CombineModeExclude    CombineMode = 4
	CombineModeComplement CombineMode = 5
)

// FlushIntention specifies when pending graphics operations are executed
const (
	FlushIntentionFlush int32 = 0
	FlushIntentionSync  int32 = 1
)

// ImageLockMode flags for Bitmap.LockBits
const (
	ImageLockModeRead         uint32 = 0x0001
	ImageLockModeWrite        uint32 = 0x0002
	ImageLockModeUserInputBuf uint32 = 0x0004
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
	ClsidBMPEncoder  = ole.GUID{Data1: 0x557CF400, Data2: 0x1A04, Data3: 0x11D3, Data4: [8]byte{0x9A, 0x73, 0x00, 0x00, 0xF8, 0x1E, 0xF3, 0x2E}}
	ClsidJPEGEncoder = ole.GUID{Data1: 0x557CF401, Data2: 0x1A04, Data3: 0x11D3, Data4: [8]byte{0x9A, 0x73, 0x00, 0x00, 0xF8, 0x1E, 0xF3, 0x2E}}
	ClsidGIFEncoder  = ole.GUID{Data1: 0x557CF402, Data2: 0x1A04, Data3: 0x11D3, Data4: [8]byte{0x9A, 0x73, 0x00, 0x00, 0xF8, 0x1E, 0xF3, 0x2E}}
	ClsidTIFFEncoder = ole.GUID{Data1: 0x557CF405, Data2: 0x1A04, Data3: 0x11D3, Data4: [8]byte{0x9A, 0x73, 0x00, 0x00, 0xF8, 0x1E, 0xF3, 0x2E}}
	ClsidPNGEncoder  = ole.GUID{Data1: 0x557CF406, Data2: 0x1A04, Data3: 0x11D3, Data4: [8]byte{0x9A, 0x73, 0x00, 0x00, 0xF8, 0x1E, 0xF3, 0x2E}}
)

var (
	EncoderQuality = ole.GUID{Data1: 0x1D5BE4B5, Data2: 0xFA4A, Data3: 0x452D, Data4: [8]byte{0x9C, 0xDD, 0x5D, 0xB3, 0x51, 0x05, 0xE7, 0xEB}}
)

const (
	EncoderParameterValueTypeLong uint32 = 4
)