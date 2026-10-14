package gdiplus

import (
	"math"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

var (
	gdiplusShutdown *windows.LazyProc
	gdiplusStartup  *windows.LazyProc
	// Graphics
	gdipCreateFromHDC          *windows.LazyProc
	gdipCreateFromHDC2         *windows.LazyProc
	gdipCreateFromHWND         *windows.LazyProc
	gdipCreateFromHWNDICM      *windows.LazyProc
	gdipDeleteGraphics         *windows.LazyProc
	gdipGetDC                  *windows.LazyProc
	gdipReleaseDC              *windows.LazyProc
	gdipSetInterpolationMode   *windows.LazyProc
	gdipSetSmoothingMode       *windows.LazyProc
	gdipSetPixelOffsetMode     *windows.LazyProc
	gdipSetCompositingQuality  *windows.LazyProc
	gdipSetCompositingMode     *windows.LazyProc
	gdipSetRenderingOrigin     *windows.LazyProc
	gdipSetTextRenderingHint   *windows.LazyProc
	gdipGraphicsClear          *windows.LazyProc
	gdipDrawLine               *windows.LazyProc
	gdipDrawLineI              *windows.LazyProc
	gdipDrawArc                *windows.LazyProc
	gdipDrawArcI               *windows.LazyProc
	gdipDrawBezier             *windows.LazyProc
	gdipDrawBezierI            *windows.LazyProc
	gdipDrawRectangle          *windows.LazyProc
	gdipDrawRectangleI         *windows.LazyProc
	gdipDrawEllipse            *windows.LazyProc
	gdipDrawEllipseI           *windows.LazyProc
	gdipDrawPie                *windows.LazyProc
	gdipDrawPieI               *windows.LazyProc
	gdipDrawPolygon            *windows.LazyProc
	gdipDrawPolygonI           *windows.LazyProc
	gdipDrawPath               *windows.LazyProc
	gdipDrawCurve              *windows.LazyProc
	gdipDrawCurve2             *windows.LazyProc
	gdipDrawClosedCurve        *windows.LazyProc
	gdipDrawString             *windows.LazyProc
	gdipDrawImage              *windows.LazyProc
	gdipDrawImageI             *windows.LazyProc
	gdipDrawImageRect          *windows.LazyProc
	gdipDrawImageRectI         *windows.LazyProc
	gdipFillRectangle          *windows.LazyProc
	gdipFillRectangleI         *windows.LazyProc
	gdipFillPolygon            *windows.LazyProc
	gdipFillPolygonI           *windows.LazyProc
	gdipFillPath               *windows.LazyProc
	gdipFillEllipse            *windows.LazyProc
	gdipFillEllipseI           *windows.LazyProc
	gdipFillPie                *windows.LazyProc
	gdipFillPieI               *windows.LazyProc
	gdipFillClosedCurve        *windows.LazyProc
	gdipMeasureString          *windows.LazyProc
	gdipMeasureCharacterRanges *windows.LazyProc
	gdipSaveGraphics           *windows.LazyProc
	gdipRestoreGraphics        *windows.LazyProc
	gdipFlush                  *windows.LazyProc
	gdipTranslateWorldTransform *windows.LazyProc
	gdipScaleWorldTransform     *windows.LazyProc
	gdipRotateWorldTransform    *windows.LazyProc
	gdipResetWorldTransform     *windows.LazyProc
	gdipSetClipRect            *windows.LazyProc
	gdipSetClipRectI           *windows.LazyProc
	gdipSetClipPath            *windows.LazyProc
	gdipResetClip              *windows.LazyProc
	// Pen
	gdipCreatePen1            *windows.LazyProc
	gdipCreatePen2            *windows.LazyProc
	gdipClonePen              *windows.LazyProc
	gdipDeletePen             *windows.LazyProc
	gdipSetPenWidth           *windows.LazyProc
	gdipGetPenWidth           *windows.LazyProc
	gdipSetPenLineCap197819   *windows.LazyProc
	gdipSetPenStartCap        *windows.LazyProc
	gdipSetPenEndCap          *windows.LazyProc
	gdipSetPenDashCap197819   *windows.LazyProc
	gdipGetPenStartCap        *windows.LazyProc
	gdipGetPenEndCap          *windows.LazyProc
	gdipGetPenDashCap197819   *windows.LazyProc
	gdipSetPenLineJoin        *windows.LazyProc
	gdipGetPenLineJoin        *windows.LazyProc
	gdipSetPenCustomStartCap  *windows.LazyProc
	gdipGetPenCustomStartCap  *windows.LazyProc
	gdipSetPenCustomEndCap    *windows.LazyProc
	gdipGetPenCustomEndCap    *windows.LazyProc
	gdipSetPenMiterLimit      *windows.LazyProc
	gdipGetPenMiterLimit      *windows.LazyProc
	gdipSetPenMode            *windows.LazyProc
	gdipGetPenMode            *windows.LazyProc
	gdipSetPenTransform       *windows.LazyProc
	gdipGetPenTransform       *windows.LazyProc
	gdipResetPenTransform     *windows.LazyProc
	gdipMultiplyPenTransform  *windows.LazyProc
	gdipTranslatePenTransform *windows.LazyProc
	gdipScalePenTransform     *windows.LazyProc
	gdipRotatePenTransform    *windows.LazyProc
	gdipSetPenColor           *windows.LazyProc
	gdipGetPenColor           *windows.LazyProc
	gdipSetPenBrushFill       *windows.LazyProc
	gdipGetPenBrushFill       *windows.LazyProc
	gdipGetPenFillType        *windows.LazyProc
	gdipGetPenDashStyle       *windows.LazyProc
	gdipSetPenDashStyle       *windows.LazyProc
	gdipGetPenDashOffset      *windows.LazyProc
	gdipSetPenDashOffset      *windows.LazyProc
	gdipGetPenDashCount       *windows.LazyProc
	gdipSetPenDashArray       *windows.LazyProc
	gdipGetPenDashArray       *windows.LazyProc
	gdipGetPenCompoundCount   *windows.LazyProc
	gdipSetPenCompoundArray   *windows.LazyProc
	gdipGetPenCompoundArray   *windows.LazyProc
	// Brush
	gdipCloneBrush   *windows.LazyProc
	gdipDeleteBrush  *windows.LazyProc
	gdipGetBrushType *windows.LazyProc
	// Solid Brush
	gdipCreateSolidFill   *windows.LazyProc
	gdipSetSolidFillColor *windows.LazyProc
	gdipGetSolidFillColor *windows.LazyProc
	// LinearGradientBrush
	gdipCreateLineBrush                  *windows.LazyProc
	gdipCreateLineBrushI                 *windows.LazyProc
	gdipCreateLineBrushFromRect          *windows.LazyProc
	gdipCreateLineBrushFromRectI         *windows.LazyProc
	gdipCreateLineBrushFromRectWithAngle *windows.LazyProc
	gdipCreateLineBrushFromRectWithAngleI *windows.LazyProc
	gdipSetLineColors                    *windows.LazyProc
	gdipGetLineColors                    *windows.LazyProc
	gdipSetLinePresetBlend               *windows.LazyProc
	gdipSetLineLinearBlend               *windows.LazyProc
	gdipSetLineSigmaBlend                *windows.LazyProc
	gdipSetLineWrapMode                  *windows.LazyProc
	gdipGetLineWrapMode                  *windows.LazyProc
	gdipSetLineGammaCorrection           *windows.LazyProc
	// HatchBrush / TextureBrush / PathGradientBrush
	gdipCreateHatchBrush                      *windows.LazyProc
	gdipCreateTexture                         *windows.LazyProc
	gdipCreatePathGradientFromPath            *windows.LazyProc
	gdipSetPathGradientCenterColor            *windows.LazyProc
	gdipSetPathGradientSurroundColorsWithCount *windows.LazyProc
	gdipSetPathGradientCenterPoint            *windows.LazyProc
	// Image
	gdipLoadImageFromFile       *windows.LazyProc
	gdipLoadImageFromStream     *windows.LazyProc
	gdipSaveImageToFile         *windows.LazyProc
	gdipSaveImageToStream       *windows.LazyProc
	gdipGetImageWidth           *windows.LazyProc
	gdipGetImageHeight          *windows.LazyProc
	gdipGetImageGraphicsContext *windows.LazyProc
	gdipDisposeImage            *windows.LazyProc
	// Bitmap
	gdipCreateBitmapFromScan0   *windows.LazyProc
	gdipCreateBitmapFromFile    *windows.LazyProc
	gdipCreateBitmapFromStream  *windows.LazyProc
	gdipCreateBitmapFromHBITMAP *windows.LazyProc
	gdipCreateHBITMAPFromBitmap *windows.LazyProc
	gdipBitmapGetPixel          *windows.LazyProc
	gdipBitmapSetPixel          *windows.LazyProc
	gdipBitmapLockBits          *windows.LazyProc
	gdipBitmapUnlockBits        *windows.LazyProc
	gdipCloneBitmapArea         *windows.LazyProc
	gdipCloneBitmapAreaI        *windows.LazyProc
	// Font & FontCollection
	gdipCreateFontFromDC              *windows.LazyProc
	gdipCreateFont                    *windows.LazyProc
	gdipDeleteFont                    *windows.LazyProc
	gdipNewInstalledFontCollection    *windows.LazyProc
	gdipNewPrivateFontCollection      *windows.LazyProc
	gdipDeletePrivateFontCollection   *windows.LazyProc
	gdipPrivateAddFontFile            *windows.LazyProc
	gdipPrivateAddMemoryFont          *windows.LazyProc
	gdipGetFontCollectionFamilyCount  *windows.LazyProc
	gdipGetFontCollectionFamilyList   *windows.LazyProc
	gdipCreateFontFamilyFromName      *windows.LazyProc
	gdipDeleteFontFamily              *windows.LazyProc
	gdipGetFamilyName                 *windows.LazyProc
	gdipGetGenericFontFamilySansSerif *windows.LazyProc
	gdipGetGenericFontFamilySerif     *windows.LazyProc
	gdipGetGenericFontFamilyMonospace *windows.LazyProc
	gdipGetEmHeight                   *windows.LazyProc
	gdipGetCellAscent                 *windows.LazyProc
	gdipGetCellDescent                *windows.LazyProc
	gdipGetLineSpacing                *windows.LazyProc
	// StringFormat
	gdipCreateStringFormat                *windows.LazyProc
	gdipDeleteStringFormat                *windows.LazyProc
	gdipStringFormatGetGenericTypographic *windows.LazyProc
	gdipStringFormatGetGenericDefault     *windows.LazyProc
	gdipCloneStringFormat                 *windows.LazyProc
	gdipSetStringFormatAlign              *windows.LazyProc
	gdipGetStringFormatAlign              *windows.LazyProc
	gdipSetStringFormatLineAlign          *windows.LazyProc
	gdipGetStringFormatLineAlign          *windows.LazyProc
	gdipSetStringFormatFlags              *windows.LazyProc
	gdipGetStringFormatFlags              *windows.LazyProc
	gdipSetStringFormatTrimming           *windows.LazyProc
	gdipGetStringFormatTrimming           *windows.LazyProc
	// Path
	gdipCreatePath        *windows.LazyProc
	gdipDeletePath        *windows.LazyProc
	gdipResetPath         *windows.LazyProc
	gdipStartPathFigure   *windows.LazyProc
	gdipAddPathArc        *windows.LazyProc
	gdipAddPathArcI       *windows.LazyProc
	gdipAddPathLine       *windows.LazyProc
	gdipAddPathLineI      *windows.LazyProc
	gdipAddPathRectangle  *windows.LazyProc
	gdipAddPathRectangleI *windows.LazyProc
	gdipAddPathEllipse    *windows.LazyProc
	gdipAddPathEllipseI   *windows.LazyProc
	gdipAddPathPie        *windows.LazyProc
	gdipAddPathPolygon    *windows.LazyProc
	gdipAddPathBezier     *windows.LazyProc
	gdipAddPathString     *windows.LazyProc
	gdipClosePathFigure   *windows.LazyProc
	gdipClosePathFigures  *windows.LazyProc
)

func init() {
	libgdiplus := windows.NewLazySystemDLL("gdiplus.dll")
	gdiplusShutdown = libgdiplus.NewProc("GdiplusShutdown")
	gdiplusStartup = libgdiplus.NewProc("GdiplusStartup")
	// Graphics
	gdipCreateFromHDC = libgdiplus.NewProc("GdipCreateFromHDC")
	gdipCreateFromHDC2 = libgdiplus.NewProc("GdipCreateFromHDC2")
	gdipCreateFromHWND = libgdiplus.NewProc("GdipCreateFromHWND")
	gdipCreateFromHWNDICM = libgdiplus.NewProc("GdipCreateFromHWNDICM")
	gdipDeleteGraphics = libgdiplus.NewProc("GdipDeleteGraphics")
	gdipGetDC = libgdiplus.NewProc("GdipGetDC")
	gdipReleaseDC = libgdiplus.NewProc("GdipReleaseDC")
	gdipSetCompositingMode = libgdiplus.NewProc("GdipSetCompositingMode")
	gdipSetRenderingOrigin = libgdiplus.NewProc("GdipSetRenderingOrigin")
	gdipSetCompositingQuality = libgdiplus.NewProc("GdipSetCompositingQuality")
	gdipSetSmoothingMode = libgdiplus.NewProc("GdipSetSmoothingMode")
	gdipSetPixelOffsetMode = libgdiplus.NewProc("GdipSetPixelOffsetMode")
	gdipSetInterpolationMode = libgdiplus.NewProc("GdipSetInterpolationMode")
	gdipSetTextRenderingHint = libgdiplus.NewProc("GdipSetTextRenderingHint")
	gdipGraphicsClear = libgdiplus.NewProc("GdipGraphicsClear")
	gdipDrawLine = libgdiplus.NewProc("GdipDrawLine")
	gdipDrawLineI = libgdiplus.NewProc("GdipDrawLineI")
	gdipDrawArc = libgdiplus.NewProc("GdipDrawArc")
	gdipDrawArcI = libgdiplus.NewProc("GdipDrawArcI")
	gdipDrawBezier = libgdiplus.NewProc("GdipDrawBezier")
	gdipDrawBezierI = libgdiplus.NewProc("GdipDrawBezierI")
	gdipDrawRectangle = libgdiplus.NewProc("GdipDrawRectangle")
	gdipDrawRectangleI = libgdiplus.NewProc("GdipDrawRectangleI")
	gdipDrawEllipse = libgdiplus.NewProc("GdipDrawEllipse")
	gdipDrawEllipseI = libgdiplus.NewProc("GdipDrawEllipseI")
	gdipDrawPie = libgdiplus.NewProc("GdipDrawPie")
	gdipDrawPieI = libgdiplus.NewProc("GdipDrawPieI")
	gdipDrawPolygonI = libgdiplus.NewProc("GdipDrawPolygonI")
	gdipDrawPolygon = libgdiplus.NewProc("GdipDrawPolygon")
	gdipDrawPath = libgdiplus.NewProc("GdipDrawPath")
	gdipDrawCurve = libgdiplus.NewProc("GdipDrawCurve")
	gdipDrawCurve2 = libgdiplus.NewProc("GdipDrawCurve2")
	gdipDrawClosedCurve = libgdiplus.NewProc("GdipDrawClosedCurve")
	gdipDrawString = libgdiplus.NewProc("GdipDrawString")
	gdipDrawImage = libgdiplus.NewProc("GdipDrawImage")
	gdipDrawImageI = libgdiplus.NewProc("GdipDrawImageI")
	gdipDrawImageRect = libgdiplus.NewProc("GdipDrawImageRect")
	gdipDrawImageRectI = libgdiplus.NewProc("GdipDrawImageRectI")
	gdipFillRectangle = libgdiplus.NewProc("GdipFillRectangle")
	gdipFillRectangleI = libgdiplus.NewProc("GdipFillRectangleI")
	gdipFillPolygon = libgdiplus.NewProc("GdipFillPolygon")
	gdipFillPolygonI = libgdiplus.NewProc("GdipFillPolygonI")
	gdipFillPath = libgdiplus.NewProc("GdipFillPath")
	gdipFillEllipse = libgdiplus.NewProc("GdipFillEllipse")
	gdipFillEllipseI = libgdiplus.NewProc("GdipFillEllipseI")
	gdipFillPie = libgdiplus.NewProc("GdipFillPie")
	gdipFillPieI = libgdiplus.NewProc("GdipFillPieI")
	gdipFillClosedCurve = libgdiplus.NewProc("GdipFillClosedCurve")
	gdipMeasureString = libgdiplus.NewProc("GdipMeasureString")
	gdipMeasureCharacterRanges = libgdiplus.NewProc("GdipMeasureCharacterRanges")
	gdipSaveGraphics = libgdiplus.NewProc("GdipSaveGraphics")
	gdipRestoreGraphics = libgdiplus.NewProc("GdipRestoreGraphics")
	gdipFlush = libgdiplus.NewProc("GdipFlush")
	gdipTranslateWorldTransform = libgdiplus.NewProc("GdipTranslateWorldTransform")
	gdipScaleWorldTransform = libgdiplus.NewProc("GdipScaleWorldTransform")
	gdipRotateWorldTransform = libgdiplus.NewProc("GdipRotateWorldTransform")
	gdipResetWorldTransform = libgdiplus.NewProc("GdipResetWorldTransform")
	gdipSetClipRect = libgdiplus.NewProc("GdipSetClipRect")
	gdipSetClipRectI = libgdiplus.NewProc("GdipSetClipRectI")
	gdipSetClipPath = libgdiplus.NewProc("GdipSetClipPath")
	gdipResetClip = libgdiplus.NewProc("GdipResetClip")
	// Pen
	gdipCreatePen1 = libgdiplus.NewProc("GdipCreatePen1")
	gdipCreatePen2 = libgdiplus.NewProc("GdipCreatePen2")
	gdipClonePen = libgdiplus.NewProc("GdipClonePen")
	gdipDeletePen = libgdiplus.NewProc("GdipDeletePen")
	gdipSetPenWidth = libgdiplus.NewProc("GdipSetPenWidth")
	gdipGetPenWidth = libgdiplus.NewProc("GdipGetPenWidth")
	gdipSetPenLineCap197819 = libgdiplus.NewProc("GdipSetPenLineCap197819")
	gdipSetPenStartCap = libgdiplus.NewProc("GdipSetPenStartCap")
	gdipSetPenEndCap = libgdiplus.NewProc("GdipSetPenEndCap")
	gdipSetPenDashCap197819 = libgdiplus.NewProc("GdipSetPenDashCap197819")
	gdipGetPenStartCap = libgdiplus.NewProc("GdipGetPenStartCap")
	gdipGetPenEndCap = libgdiplus.NewProc("GdipGetPenEndCap")
	gdipGetPenDashCap197819 = libgdiplus.NewProc("GdipGetPenDashCap197819")
	gdipSetPenLineJoin = libgdiplus.NewProc("GdipSetPenLineJoin")
	gdipGetPenLineJoin = libgdiplus.NewProc("GdipGetPenLineJoin")
	gdipSetPenCustomStartCap = libgdiplus.NewProc("GdipSetPenCustomStartCap")
	gdipGetPenCustomStartCap = libgdiplus.NewProc("GdipGetPenCustomStartCap")
	gdipSetPenCustomEndCap = libgdiplus.NewProc("GdipSetPenCustomEndCap")
	gdipGetPenCustomEndCap = libgdiplus.NewProc("GdipGetPenCustomEndCap")
	gdipSetPenMiterLimit = libgdiplus.NewProc("GdipSetPenMiterLimit")
	gdipGetPenMiterLimit = libgdiplus.NewProc("GdipGetPenMiterLimit")
	gdipSetPenMode = libgdiplus.NewProc("GdipSetPenMode")
	gdipGetPenMode = libgdiplus.NewProc("GdipGetPenMode")
	gdipSetPenTransform = libgdiplus.NewProc("GdipSetPenTransform")
	gdipGetPenTransform = libgdiplus.NewProc("GdipGetPenTransform")
	gdipResetPenTransform = libgdiplus.NewProc("GdipResetPenTransform")
	gdipMultiplyPenTransform = libgdiplus.NewProc("GdipMultiplyPenTransform")
	gdipTranslatePenTransform = libgdiplus.NewProc("GdipTranslatePenTransform")
	gdipScalePenTransform = libgdiplus.NewProc("GdipScalePenTransform")
	gdipRotatePenTransform = libgdiplus.NewProc("GdipRotatePenTransform")
	gdipSetPenColor = libgdiplus.NewProc("GdipSetPenColor")
	gdipGetPenColor = libgdiplus.NewProc("GdipGetPenColor")
	gdipSetPenBrushFill = libgdiplus.NewProc("GdipSetPenBrushFill")
	gdipGetPenBrushFill = libgdiplus.NewProc("GdipGetPenBrushFill")
	gdipGetPenFillType = libgdiplus.NewProc("GdipGetPenFillType")
	gdipGetPenDashStyle = libgdiplus.NewProc("GdipGetPenDashStyle")
	gdipSetPenDashStyle = libgdiplus.NewProc("GdipSetPenDashStyle")
	gdipGetPenDashOffset = libgdiplus.NewProc("GdipGetPenDashOffset")
	gdipSetPenDashOffset = libgdiplus.NewProc("GdipSetPenDashOffset")
	gdipGetPenDashCount = libgdiplus.NewProc("GdipGetPenDashCount")
	gdipSetPenDashArray = libgdiplus.NewProc("GdipSetPenDashArray")
	gdipGetPenDashArray = libgdiplus.NewProc("GdipGetPenDashArray")
	gdipGetPenCompoundCount = libgdiplus.NewProc("GdipGetPenCompoundCount")
	gdipSetPenCompoundArray = libgdiplus.NewProc("GdipSetPenCompoundArray")
	gdipGetPenCompoundArray = libgdiplus.NewProc("GdipGetPenCompoundArray")
	// Brush
	gdipCloneBrush = libgdiplus.NewProc("GdipCloneBrush")
	gdipDeleteBrush = libgdiplus.NewProc("GdipDeleteBrush")
	gdipGetBrushType = libgdiplus.NewProc("GdipGetBrushType")
	// Solid Brush
	gdipCreateSolidFill = libgdiplus.NewProc("GdipCreateSolidFill")
	gdipSetSolidFillColor = libgdiplus.NewProc("GdipSetSolidFillColor")
	gdipGetSolidFillColor = libgdiplus.NewProc("GdipGetSolidFillColor")
	// LinearGradientBrush
	gdipCreateLineBrush = libgdiplus.NewProc("GdipCreateLineBrush")
	gdipCreateLineBrushI = libgdiplus.NewProc("GdipCreateLineBrushI")
	gdipCreateLineBrushFromRect = libgdiplus.NewProc("GdipCreateLineBrushFromRect")
	gdipCreateLineBrushFromRectI = libgdiplus.NewProc("GdipCreateLineBrushFromRectI")
	gdipCreateLineBrushFromRectWithAngle = libgdiplus.NewProc("GdipCreateLineBrushFromRectWithAngle")
	gdipCreateLineBrushFromRectWithAngleI = libgdiplus.NewProc("GdipCreateLineBrushFromRectWithAngleI")
	gdipSetLineColors = libgdiplus.NewProc("GdipSetLineColors")
	gdipGetLineColors = libgdiplus.NewProc("GdipGetLineColors")
	gdipSetLinePresetBlend = libgdiplus.NewProc("GdipSetLinePresetBlend")
	gdipSetLineLinearBlend = libgdiplus.NewProc("GdipSetLineLinearBlend")
	gdipSetLineSigmaBlend = libgdiplus.NewProc("GdipSetLineSigmaBlend")
	gdipSetLineWrapMode = libgdiplus.NewProc("GdipSetLineWrapMode")
	gdipGetLineWrapMode = libgdiplus.NewProc("GdipGetLineWrapMode")
	gdipSetLineGammaCorrection = libgdiplus.NewProc("GdipSetLineGammaCorrection")
	// Hatch / Texture / PathGradient
	gdipCreateHatchBrush = libgdiplus.NewProc("GdipCreateHatchBrush")
	gdipCreateTexture = libgdiplus.NewProc("GdipCreateTexture")
	gdipCreatePathGradientFromPath = libgdiplus.NewProc("GdipCreatePathGradientFromPath")
	gdipSetPathGradientCenterColor = libgdiplus.NewProc("GdipSetPathGradientCenterColor")
	gdipSetPathGradientSurroundColorsWithCount = libgdiplus.NewProc("GdipSetPathGradientSurroundColorsWithCount")
	gdipSetPathGradientCenterPoint = libgdiplus.NewProc("GdipSetPathGradientCenterPoint")
	// Image
	gdipLoadImageFromFile = libgdiplus.NewProc("GdipLoadImageFromFile")
	gdipLoadImageFromStream = libgdiplus.NewProc("GdipLoadImageFromStream")
	gdipSaveImageToFile = libgdiplus.NewProc("GdipSaveImageToFile")
	gdipSaveImageToStream = libgdiplus.NewProc("GdipSaveImageToStream")
	gdipGetImageWidth = libgdiplus.NewProc("GdipGetImageWidth")
	gdipGetImageHeight = libgdiplus.NewProc("GdipGetImageHeight")
	gdipGetImageGraphicsContext = libgdiplus.NewProc("GdipGetImageGraphicsContext")
	gdipDisposeImage = libgdiplus.NewProc("GdipDisposeImage")
	// Bitmap
	gdipCreateBitmapFromScan0 = libgdiplus.NewProc("GdipCreateBitmapFromScan0")
	gdipCreateBitmapFromFile = libgdiplus.NewProc("GdipCreateBitmapFromFile")
	gdipCreateBitmapFromStream = libgdiplus.NewProc("GdipCreateBitmapFromStream")
	gdipCreateBitmapFromHBITMAP = libgdiplus.NewProc("GdipCreateBitmapFromHBITMAP")
	gdipCreateHBITMAPFromBitmap = libgdiplus.NewProc("GdipCreateHBITMAPFromBitmap")
	gdipBitmapGetPixel = libgdiplus.NewProc("GdipBitmapGetPixel")
	gdipBitmapSetPixel = libgdiplus.NewProc("GdipBitmapSetPixel")
	gdipBitmapLockBits = libgdiplus.NewProc("GdipBitmapLockBits")
	gdipBitmapUnlockBits = libgdiplus.NewProc("GdipBitmapUnlockBits")
	gdipCloneBitmapArea = libgdiplus.NewProc("GdipCloneBitmapArea")
	gdipCloneBitmapAreaI = libgdiplus.NewProc("GdipCloneBitmapAreaI")
	// Font
	gdipCreateFontFromDC = libgdiplus.NewProc("GdipCreateFontFromDC")
	gdipCreateFont = libgdiplus.NewProc("GdipCreateFont")
	gdipDeleteFont = libgdiplus.NewProc("GdipDeleteFont")
	gdipNewInstalledFontCollection = libgdiplus.NewProc("GdipNewInstalledFontCollection")
	gdipNewPrivateFontCollection = libgdiplus.NewProc("GdipNewPrivateFontCollection")
	gdipDeletePrivateFontCollection = libgdiplus.NewProc("GdipDeletePrivateFontCollection")
	gdipPrivateAddFontFile = libgdiplus.NewProc("GdipPrivateAddFontFile")
	gdipPrivateAddMemoryFont = libgdiplus.NewProc("GdipPrivateAddMemoryFont")
	gdipGetFontCollectionFamilyCount = libgdiplus.NewProc("GdipGetFontCollectionFamilyCount")
	gdipGetFontCollectionFamilyList = libgdiplus.NewProc("GdipGetFontCollectionFamilyList")
	gdipCreateFontFamilyFromName = libgdiplus.NewProc("GdipCreateFontFamilyFromName")
	gdipDeleteFontFamily = libgdiplus.NewProc("GdipDeleteFontFamily")
	gdipGetFamilyName = libgdiplus.NewProc("GdipGetFamilyName")
	gdipGetGenericFontFamilySansSerif = libgdiplus.NewProc("GdipGetGenericFontFamilySansSerif")
	gdipGetGenericFontFamilySerif = libgdiplus.NewProc("GdipGetGenericFontFamilySerif")
	gdipGetGenericFontFamilyMonospace = libgdiplus.NewProc("GdipGetGenericFontFamilyMonospace")
	gdipGetEmHeight = libgdiplus.NewProc("GdipGetEmHeight")
	gdipGetCellAscent = libgdiplus.NewProc("GdipGetCellAscent")
	gdipGetCellDescent = libgdiplus.NewProc("GdipGetCellDescent")
	gdipGetLineSpacing = libgdiplus.NewProc("GdipGetLineSpacing")
	// StringFormat
	gdipCreateStringFormat = libgdiplus.NewProc("GdipCreateStringFormat")
	gdipDeleteStringFormat = libgdiplus.NewProc("GdipDeleteStringFormat")
	gdipStringFormatGetGenericTypographic = libgdiplus.NewProc("GdipStringFormatGetGenericTypographic")
	gdipStringFormatGetGenericDefault = libgdiplus.NewProc("GdipStringFormatGetGenericDefault")
	gdipCloneStringFormat = libgdiplus.NewProc("GdipCloneStringFormat")
	gdipSetStringFormatAlign = libgdiplus.NewProc("GdipSetStringFormatAlign")
	gdipGetStringFormatAlign = libgdiplus.NewProc("GdipGetStringFormatAlign")
	gdipSetStringFormatLineAlign = libgdiplus.NewProc("GdipSetStringFormatLineAlign")
	gdipGetStringFormatLineAlign = libgdiplus.NewProc("GdipGetStringFormatLineAlign")
	gdipSetStringFormatFlags = libgdiplus.NewProc("GdipSetStringFormatFlags")
	gdipGetStringFormatFlags = libgdiplus.NewProc("GdipGetStringFormatFlags")
	gdipSetStringFormatTrimming = libgdiplus.NewProc("GdipSetStringFormatTrimming")
	gdipGetStringFormatTrimming = libgdiplus.NewProc("GdipGetStringFormatTrimming")
	// Path
	gdipCreatePath = libgdiplus.NewProc("GdipCreatePath")
	gdipDeletePath = libgdiplus.NewProc("GdipDeletePath")
	gdipResetPath = libgdiplus.NewProc("GdipResetPath")
	gdipStartPathFigure = libgdiplus.NewProc("GdipStartPathFigure")
	gdipAddPathArc = libgdiplus.NewProc("GdipAddPathArc")
	gdipAddPathArcI = libgdiplus.NewProc("GdipAddPathArcI")
	gdipAddPathLine = libgdiplus.NewProc("GdipAddPathLine")
	gdipAddPathLineI = libgdiplus.NewProc("GdipAddPathLineI")
	gdipAddPathRectangle = libgdiplus.NewProc("GdipAddPathRectangle")
	gdipAddPathRectangleI = libgdiplus.NewProc("GdipAddPathRectangleI")
	gdipAddPathEllipse = libgdiplus.NewProc("GdipAddPathEllipse")
	gdipAddPathEllipseI = libgdiplus.NewProc("GdipAddPathEllipseI")
	gdipAddPathPie = libgdiplus.NewProc("GdipAddPathPie")
	gdipAddPathPolygon = libgdiplus.NewProc("GdipAddPathPolygon")
	gdipAddPathBezier = libgdiplus.NewProc("GdipAddPathBezier")
	gdipAddPathString = libgdiplus.NewProc("GdipAddPathString")
	gdipClosePathFigure = libgdiplus.NewProc("GdipClosePathFigure")
	gdipClosePathFigures = libgdiplus.NewProc("GdipClosePathFigures")
}

var (
	token uintptr
)

func GdiplusShutdown() {
	clearFallbackFontCache()
	syscall.Syscall(gdiplusShutdown.Addr(), 1,
		token,
		0,
		0)
}

func GdiplusStartup(input *GdiplusStartupInput, output *GdiplusStartupOutput) GpStatus {
	ret, _, _ := syscall.Syscall(gdiplusStartup.Addr(), 3,
		uintptr(unsafe.Pointer(&token)),
		uintptr(unsafe.Pointer(input)),
		uintptr(unsafe.Pointer(output)))

	return GpStatus(ret)
}

// Graphics
func GdipCreateFromHDC(hdc HDC, graphics **GpGraphics) GpStatus {
	ret, _, _ := gdipCreateFromHDC.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipCreateFromHDC2(hdc HDC, hDevice HANDLE, graphics **GpGraphics) GpStatus {
	ret, _, _ := gdipCreateFromHDC2.Call(
		uintptr(hdc),
		uintptr(hDevice),
		uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipCreateFromHWND(hwnd HWND, graphics **GpGraphics) GpStatus {
	ret, _, _ := gdipCreateFromHWND.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipCreateFromHWNDICM(hwnd HWND, graphics **GpGraphics) GpStatus {
	ret, _, _ := gdipCreateFromHWNDICM.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipDeleteGraphics(graphics *GpGraphics) GpStatus {
	ret, _, _ := gdipDeleteGraphics.Call(uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipGetDC(graphics *GpGraphics, hdc *HDC) GpStatus {
	ret, _, _ := gdipGetDC.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(hdc)))
	return GpStatus(ret)
}

func GdipReleaseDC(graphics *GpGraphics, hdc HDC) GpStatus {
	ret, _, _ := gdipReleaseDC.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(hdc))
	return GpStatus(ret)
}

func GdipSetCompositingMode(graphics *GpGraphics, mode int32) GpStatus {
	ret, _, _ := gdipSetCompositingMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return GpStatus(ret)
}

func GdipSetRenderingOrigin(graphics *GpGraphics, x, y int32) GpStatus {
	ret, _, _ := gdipSetRenderingOrigin.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(x),
		uintptr(y))
	return GpStatus(ret)
}

func GdipSetCompositingQuality(graphics *GpGraphics, quality int32) GpStatus {
	ret, _, _ := gdipSetCompositingQuality.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(quality))
	return GpStatus(ret)
}

func GdipSetInterpolationMode(graphics *GpGraphics, mode int32) GpStatus {
	ret, _, _ := gdipSetInterpolationMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return GpStatus(ret)
}

func GdipSetPixelOffsetMode(graphics *GpGraphics, mode int32) GpStatus {
	ret, _, _ := gdipSetPixelOffsetMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return GpStatus(ret)
}

func GdipSetSmoothingMode(graphics *GpGraphics, mode int32) GpStatus {
	ret, _, _ := gdipSetSmoothingMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return GpStatus(ret)
}

func GdipSetTextRenderingHint(graphics *GpGraphics, hint int32) GpStatus {
	ret, _, _ := gdipSetTextRenderingHint.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(hint))
	return GpStatus(ret)
}

func GdipGraphicsClear(graphics *GpGraphics, color ARGB) GpStatus {
	ret, _, _ := gdipGraphicsClear.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(color))
	return GpStatus(ret)
}

func GdipSaveGraphics(graphics *GpGraphics, state *uint32) GpStatus {
	ret, _, _ := gdipSaveGraphics.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(state)))
	return GpStatus(ret)
}

func GdipRestoreGraphics(graphics *GpGraphics, state uint32) GpStatus {
	ret, _, _ := gdipRestoreGraphics.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(state))
	return GpStatus(ret)
}

func GdipFlush(graphics *GpGraphics, intention int32) GpStatus {
	ret, _, _ := gdipFlush.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(intention))
	return GpStatus(ret)
}

func GdipTranslateWorldTransform(graphics *GpGraphics, dx, dy float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := gdipTranslateWorldTransform.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(math.Float32bits(dx)),
		uintptr(math.Float32bits(dy)),
		uintptr(order))
	return GpStatus(ret)
}

func GdipScaleWorldTransform(graphics *GpGraphics, sx, sy float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := gdipScaleWorldTransform.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(math.Float32bits(sx)),
		uintptr(math.Float32bits(sy)),
		uintptr(order))
	return GpStatus(ret)
}

func GdipRotateWorldTransform(graphics *GpGraphics, angle float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := gdipRotateWorldTransform.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(math.Float32bits(angle)),
		uintptr(order))
	return GpStatus(ret)
}

func GdipResetWorldTransform(graphics *GpGraphics) GpStatus {
	ret, _, _ := gdipResetWorldTransform.Call(uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipSetClipRect(graphics *GpGraphics, x, y, width, height float32, combineMode CombineMode) GpStatus {
	ret, _, _ := gdipSetClipRect.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(combineMode))
	return GpStatus(ret)
}

func GdipSetClipRectI(graphics *GpGraphics, x, y, width, height int32, combineMode CombineMode) GpStatus {
	ret, _, _ := gdipSetClipRectI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(combineMode))
	return GpStatus(ret)
}

func GdipSetClipPath(graphics *GpGraphics, path *GpPath, combineMode CombineMode) GpStatus {
	ret, _, _ := gdipSetClipPath.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(path)),
		uintptr(combineMode))
	return GpStatus(ret)
}

func GdipResetClip(graphics *GpGraphics) GpStatus {
	ret, _, _ := gdipResetClip.Call(uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipDrawLine(graphics *GpGraphics, pen *GpPen, x1, y1, x2, y2 float32) GpStatus {
	ret, _, _ := gdipDrawLine.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x1)),
		uintptr(math.Float32bits(y1)),
		uintptr(math.Float32bits(x2)),
		uintptr(math.Float32bits(y2)))
	return GpStatus(ret)
}

func GdipDrawLineI(graphics *GpGraphics, pen *GpPen, x1, y1, x2, y2 int32) GpStatus {
	ret, _, _ := gdipDrawLineI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x1),
		uintptr(y1),
		uintptr(x2),
		uintptr(y2))
	return GpStatus(ret)
}

func GdipDrawArc(graphics *GpGraphics, pen *GpPen, x, y, width, height, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := gdipDrawArc.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipDrawArcI(graphics *GpGraphics, pen *GpPen, x, y, width, height int32, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := gdipDrawArcI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipDrawBezier(graphics *GpGraphics, pen *GpPen, x1, y1, x2, y2, x3, y3, x4, y4 float32) GpStatus {
	ret, _, _ := gdipDrawBezier.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x1)),
		uintptr(math.Float32bits(y1)),
		uintptr(math.Float32bits(x2)),
		uintptr(math.Float32bits(y2)),
		uintptr(math.Float32bits(x3)),
		uintptr(math.Float32bits(y3)),
		uintptr(math.Float32bits(x4)),
		uintptr(math.Float32bits(y4)))
	return GpStatus(ret)
}

func GdipDrawBezierI(graphics *GpGraphics, pen *GpPen, x1, y1, x2, y2, x3, y3, x4, y4 int32) GpStatus {
	ret, _, _ := gdipDrawBezierI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x1),
		uintptr(y1),
		uintptr(x2),
		uintptr(y2),
		uintptr(x3),
		uintptr(y3),
		uintptr(x4),
		uintptr(y4))
	return GpStatus(ret)
}

func GdipDrawRectangle(graphics *GpGraphics, pen *GpPen, x, y, width, height float32) GpStatus {
	ret, _, _ := gdipDrawRectangle.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipDrawRectangleI(graphics *GpGraphics, pen *GpPen, x, y, width, height int32) GpStatus {
	ret, _, _ := gdipDrawRectangleI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipDrawEllipse(graphics *GpGraphics, pen *GpPen, x, y, width, height float32) GpStatus {
	ret, _, _ := gdipDrawEllipse.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipDrawEllipseI(graphics *GpGraphics, pen *GpPen, x, y, width, height int32) GpStatus {
	ret, _, _ := gdipDrawEllipseI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipDrawPie(graphics *GpGraphics, pen *GpPen, x, y, width, height, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := gdipDrawPie.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipDrawPieI(graphics *GpGraphics, pen *GpPen, x, y, width, height int32, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := gdipDrawPieI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipDrawPolygon(graphics *GpGraphics, pen *GpPen, points *PointF, count int32) GpStatus {
	ret, _, _ := gdipDrawPolygon.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipDrawPolygonI(graphics *GpGraphics, pen *GpPen, points *Point, count int32) GpStatus {
	ret, _, _ := gdipDrawPolygonI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipDrawPath(graphics *GpGraphics, pen *GpPen, path *GpPath) GpStatus {
	ret, _, _ := gdipDrawPath.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipDrawCurve(graphics *GpGraphics, pen *GpPen, points *PointF, count int32) GpStatus {
	ret, _, _ := gdipDrawCurve.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipDrawCurve2(graphics *GpGraphics, pen *GpPen, points *PointF, count int32, tension float32) GpStatus {
	ret, _, _ := gdipDrawCurve2.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count),
		uintptr(math.Float32bits(tension)))
	return GpStatus(ret)
}

func GdipDrawClosedCurve(graphics *GpGraphics, pen *GpPen, points *PointF, count int32) GpStatus {
	ret, _, _ := gdipDrawClosedCurve.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipDrawString(graphics *GpGraphics, text *uint16, length int32, font *GpFont, layoutRect *RectF, stringFormat *GpStringFormat, brush *GpBrush) GpStatus {
	ret, _, _ := gdipDrawString.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(text)),
		uintptr(length),
		uintptr(unsafe.Pointer(font)),
		uintptr(unsafe.Pointer(layoutRect)),
		uintptr(unsafe.Pointer(stringFormat)),
		uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}

func GdipDrawImage(graphics *GpGraphics, image *GpImage, x, y float32) GpStatus {
	ret, _, _ := gdipDrawImage.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)))
	return GpStatus(ret)
}

func GdipDrawImageI(graphics *GpGraphics, image *GpImage, x, y int32) GpStatus {
	ret, _, _ := gdipDrawImageI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(x),
		uintptr(y))
	return GpStatus(ret)
}

func GdipDrawImageRect(graphics *GpGraphics, image *GpImage, x, y, width, height float32) GpStatus {
	ret, _, _ := gdipDrawImageRect.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipDrawImageRectI(graphics *GpGraphics, image *GpImage, x, y, width, height int32) GpStatus {
	ret, _, _ := gdipDrawImageRectI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipFillRectangle(graphics *GpGraphics, brush *GpBrush, x, y, width, height float32) GpStatus {
	ret, _, _ := gdipFillRectangle.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipFillRectangleI(graphics *GpGraphics, brush *GpBrush, x, y, width, height int32) GpStatus {
	ret, _, _ := gdipFillRectangleI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipFillEllipse(graphics *GpGraphics, brush *GpBrush, x, y, width, height float32) GpStatus {
	ret, _, _ := gdipFillEllipse.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipFillEllipseI(graphics *GpGraphics, brush *GpBrush, x, y, width, height int32) GpStatus {
	ret, _, _ := gdipFillEllipseI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipFillPie(graphics *GpGraphics, brush *GpBrush, x, y, width, height, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := gdipFillPie.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipFillPieI(graphics *GpGraphics, brush *GpBrush, x, y, width, height int32, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := gdipFillPieI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipFillPolygon(graphics *GpGraphics, brush *GpBrush, points *PointF, count int32, fillMode int32) GpStatus {
	ret, _, _ := gdipFillPolygon.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count),
		uintptr(fillMode))
	return GpStatus(ret)
}

func GdipFillPolygonI(graphics *GpGraphics, brush *GpBrush, points *Point, count int32, fillMode int32) GpStatus {
	ret, _, _ := gdipFillPolygonI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count),
		uintptr(fillMode))
	return GpStatus(ret)
}

func GdipFillClosedCurve(graphics *GpGraphics, brush *GpBrush, points *PointF, count int32) GpStatus {
	ret, _, _ := gdipFillClosedCurve.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipFillPath(graphics *GpGraphics, brush *GpBrush, path *GpPath) GpStatus {
	ret, _, _ := gdipFillPath.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipMeasureString(
	graphics *GpGraphics, text *uint16,
	length int32, font *GpFont, layoutRect *RectF,
	stringFormat *GpStringFormat, boundingBox *RectF,
	codepointsFitted *int32, linesFilled *int32) GpStatus {

	ret, _, _ := gdipMeasureString.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(text)),
		uintptr(length),
		uintptr(unsafe.Pointer(font)),
		uintptr(unsafe.Pointer(layoutRect)),
		uintptr(unsafe.Pointer(stringFormat)),
		uintptr(unsafe.Pointer(boundingBox)),
		uintptr(unsafe.Pointer(codepointsFitted)),
		uintptr(unsafe.Pointer(linesFilled)))
	return GpStatus(ret)
}

func GdipMeasureCharacterRanges(
	graphics *GpGraphics, text *uint16,
	length int32, font *GpFont, layoutRect *RectF,
	stringFormat *GpStringFormat, regionCount int32,
	regions **GpRegion) GpStatus {

	ret, _, _ := gdipMeasureCharacterRanges.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(text)),
		uintptr(length),
		uintptr(unsafe.Pointer(font)),
		uintptr(unsafe.Pointer(layoutRect)),
		uintptr(unsafe.Pointer(stringFormat)),
		uintptr(regionCount),
		uintptr(unsafe.Pointer(regions)))
	return GpStatus(ret)
}

// Pen
func GdipCreatePen1(color ARGB, width float32, unit GpUnit, pen **GpPen) GpStatus {
	ret, _, _ := gdipCreatePen1.Call(
		uintptr(color),
		uintptr(math.Float32bits(width)),
		uintptr(unit),
		uintptr(unsafe.Pointer(pen)))
	return GpStatus(ret)
}

func GdipCreatePen2(brush *GpBrush, width float32, unit GpUnit, pen **GpPen) GpStatus {
	ret, _, _ := gdipCreatePen2.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(math.Float32bits(width)),
		uintptr(unit),
		uintptr(unsafe.Pointer(pen)))
	return GpStatus(ret)
}

func GdipClonePen(pen *GpPen, clonepen **GpPen) GpStatus {
	ret, _, _ := gdipClonePen.Call(uintptr(unsafe.Pointer(pen)), uintptr(unsafe.Pointer(clonepen)))
	return GpStatus(ret)
}

func GdipDeletePen(pen *GpPen) GpStatus {
	ret, _, _ := gdipDeletePen.Call(uintptr(unsafe.Pointer(pen)))
	return GpStatus(ret)
}

func GdipSetPenWidth(pen *GpPen, width float32) GpStatus {
	ret, _, _ := gdipSetPenWidth.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(width)))
	return GpStatus(ret)
}

func GdipGetPenWidth(pen *GpPen, width *float32) GpStatus {
	var penWidth uint32
	ret, _, _ := gdipGetPenWidth.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(&penWidth)))
	*width = math.Float32frombits(penWidth)
	return GpStatus(ret)
}

func GdipSetPenLineCap197819(pen *GpPen, startCap, endCap GpLineCap, dashCap GpDashCap) GpStatus {
	ret, _, _ := gdipSetPenLineCap197819.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(startCap),
		uintptr(endCap),
		uintptr(dashCap))
	return GpStatus(ret)
}
func GdipSetPenStartCap(pen *GpPen, startCap GpLineCap) GpStatus {
	ret, _, _ := gdipSetPenStartCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(startCap))
	return GpStatus(ret)
}
func GdipSetPenEndCap(pen *GpPen, endCap GpLineCap) GpStatus {
	ret, _, _ := gdipSetPenEndCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(endCap))
	return GpStatus(ret)
}
func GdipSetPenDashCap197819(pen *GpPen, dashCap GpDashCap) GpStatus {
	ret, _, _ := gdipSetPenDashCap197819.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(dashCap))
	return GpStatus(ret)
}
func GdipGetPenStartCap(pen *GpPen, startCap *GpLineCap) GpStatus {
	ret, _, _ := gdipGetPenStartCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(startCap)))
	return GpStatus(ret)
}
func GdipGetPenEndCap(pen *GpPen, endCap *GpLineCap) GpStatus {
	ret, _, _ := gdipGetPenEndCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(endCap)))
	return GpStatus(ret)
}
func GdipGetPenDashCap197819(pen *GpPen, dashCap *GpDashCap) GpStatus {
	ret, _, _ := gdipGetPenDashCap197819.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dashCap)))
	return GpStatus(ret)
}
func GdipSetPenLineJoin(pen *GpPen, lineJoin GpLineJoin) GpStatus {
	ret, _, _ := gdipSetPenLineJoin.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(lineJoin))
	return GpStatus(ret)
}
func GdipGetPenLineJoin(pen *GpPen, lineJoin *GpLineJoin) GpStatus {
	ret, _, _ := gdipGetPenLineJoin.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(lineJoin)))
	return GpStatus(ret)
}
func GdipSetPenCustomStartCap(pen *GpPen, customCap *GpCustomLineCap) GpStatus {
	ret, _, _ := gdipSetPenCustomStartCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(customCap)))
	return GpStatus(ret)
}
func GdipGetPenCustomStartCap(pen *GpPen, customCap **GpCustomLineCap) GpStatus {
	ret, _, _ := gdipGetPenCustomStartCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(customCap)))
	return GpStatus(ret)
}
func GdipSetPenCustomEndCap(pen *GpPen, customCap *GpCustomLineCap) GpStatus {
	ret, _, _ := gdipSetPenCustomEndCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(customCap)))
	return GpStatus(ret)
}
func GdipGetPenCustomEndCap(pen *GpPen, customCap **GpCustomLineCap) GpStatus {
	ret, _, _ := gdipGetPenCustomEndCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(customCap)))
	return GpStatus(ret)
}
func GdipSetPenMiterLimit(pen *GpPen, miterLimit float32) GpStatus {
	ret, _, _ := gdipSetPenMiterLimit.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(miterLimit)))
	return GpStatus(ret)
}
func GdipGetPenMiterLimit(pen *GpPen, miterLimit *float32) GpStatus {
	var iMiterLimit uint32
	ret, _, _ := gdipGetPenMiterLimit.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(&iMiterLimit)))
	*miterLimit = math.Float32frombits(iMiterLimit)
	return GpStatus(ret)
}
func GdipSetPenMode(pen *GpPen, penMode GpPenAlignment) GpStatus {
	ret, _, _ := gdipSetPenMode.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(penMode))
	return GpStatus(ret)
}
func GdipGetPenMode(pen *GpPen, penMode *GpPenAlignment) GpStatus {
	ret, _, _ := gdipGetPenMode.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(penMode)))
	return GpStatus(ret)
}
func GdipSetPenTransform(pen *GpPen, matrix *GpMatrix) GpStatus {
	ret, _, _ := gdipSetPenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(matrix)))
	return GpStatus(ret)
}
func GdipGetPenTransform(pen *GpPen, matrix *GpMatrix) GpStatus {
	ret, _, _ := gdipGetPenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(matrix)))
	return GpStatus(ret)
}
func GdipResetPenTransform(pen *GpPen) GpStatus {
	ret, _, _ := gdipResetPenTransform.Call(uintptr(unsafe.Pointer(pen)))
	return GpStatus(ret)
}
func GdipMultiplyPenTransform(pen *GpPen, matrix *GpMatrix, order GpMatrixOrder) GpStatus {
	ret, _, _ := gdipMultiplyPenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(matrix)),
		uintptr(order))
	return GpStatus(ret)
}
func GdipTranslatePenTransform(pen *GpPen, dx, dy float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := gdipTranslatePenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(dx)),
		uintptr(math.Float32bits(dy)),
		uintptr(order))
	return GpStatus(ret)
}
func GdipScalePenTransform(pen *GpPen, sx, sy float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := gdipScalePenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(sx)),
		uintptr(math.Float32bits(sy)),
		uintptr(order))
	return GpStatus(ret)
}
func GdipRotatePenTransform(pen *GpPen, angle float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := gdipRotatePenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(angle)),
		uintptr(order))
	return GpStatus(ret)
}
func GdipSetPenColor(pen *GpPen, argb ARGB) GpStatus {
	ret, _, _ := gdipSetPenColor.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(argb))
	return GpStatus(ret)
}
func GdipGetPenColor(pen *GpPen, argb *ARGB) GpStatus {
	ret, _, _ := gdipGetPenColor.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(argb)))
	return GpStatus(ret)
}
func GdipSetPenBrushFill(pen *GpPen, brush *GpBrush) GpStatus {
	ret, _, _ := gdipSetPenBrushFill.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}
func GdipGetPenBrushFill(pen *GpPen, brush **GpBrush) GpStatus {
	ret, _, _ := gdipGetPenBrushFill.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}
func GdipGetPenFillType(pen *GpPen, penType *GpPenType) GpStatus {
	ret, _, _ := gdipGetPenFillType.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(penType)))
	return GpStatus(ret)
}
func GdipGetPenDashStyle(pen *GpPen, dashStyle *GpDashStyle) GpStatus {
	ret, _, _ := gdipGetPenDashStyle.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dashStyle)))
	return GpStatus(ret)
}
func GdipSetPenDashStyle(pen *GpPen, dashStyle GpDashStyle) GpStatus {
	ret, _, _ := gdipSetPenDashStyle.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(dashStyle))
	return GpStatus(ret)
}
func GdipGetPenDashOffset(pen *GpPen, offset *float32) GpStatus {
	var iOffset uint32
	ret, _, _ := gdipGetPenDashOffset.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(&iOffset)))
	*offset = math.Float32frombits(iOffset)
	return GpStatus(ret)
}
func GdipSetPenDashOffset(pen *GpPen, offset float32) GpStatus {
	ret, _, _ := gdipSetPenDashOffset.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(offset)))
	return GpStatus(ret)
}
func GdipGetPenDashCount(pen *GpPen, count *int32) GpStatus {
	ret, _, _ := gdipGetPenDashCount.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(count)))
	return GpStatus(ret)
}
func GdipSetPenDashArray(pen *GpPen, dash *float32, count int32) GpStatus {
	ret, _, _ := gdipSetPenDashArray.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dash)),
		uintptr(count))
	return GpStatus(ret)
}
func GdipGetPenDashArray(pen *GpPen, dash *float32, count int32) GpStatus {
	ret, _, _ := gdipGetPenDashArray.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dash)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipGetPenCompoundCount(pen *GpPen, count *int32) GpStatus {
	ret, _, _ := gdipGetPenCompoundCount.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(count)))
	return GpStatus(ret)
}

func GdipSetPenCompoundArray(pen *GpPen, dash *float32, count int32) GpStatus {
	ret, _, _ := gdipSetPenCompoundArray.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dash)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipGetPenCompoundArray(pen *GpPen, dash *float32, count int32) GpStatus {
	ret, _, _ := gdipGetPenCompoundArray.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dash)),
		uintptr(count))
	return GpStatus(ret)
}

// Brush
func GdipCloneBrush(brush *GpBrush, clone **GpBrush) GpStatus {
	ret, _, _ := gdipCloneBrush.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(clone)))
	return GpStatus(ret)
}

func GdipDeleteBrush(brush *GpBrush) GpStatus {
	ret, _, _ := gdipDeleteBrush.Call(uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}

func GdipGetBrushType(brush *GpBrush, brushType *GpBrushType) GpStatus {
	ret, _, _ := gdipGetBrushType.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(brushType)))
	return GpStatus(ret)
}

// Solid Brush
func GdipCreateSolidFill(color ARGB, brush **GpSolidFill) GpStatus {
	ret, _, _ := gdipCreateSolidFill.Call(
		uintptr(color),
		uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}

func GdipSetSolidFillColor(brush *GpBrush, color ARGB) GpStatus {
	ret, _, _ := gdipSetSolidFillColor.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(color))
	return GpStatus(ret)
}

func GdipGetSolidFillColor(brush *GpBrush, color *ARGB) GpStatus {
	ret, _, _ := gdipGetSolidFillColor.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(color)))
	return GpStatus(ret)
}

// LinearGradientBrush
func GdipCreateLineBrush(point1, point2 *PointF, color1, color2 ARGB, wrapMode WrapMode, lineGradient **GpLineGradient) GpStatus {
	ret, _, _ := gdipCreateLineBrush.Call(
		uintptr(unsafe.Pointer(point1)),
		uintptr(unsafe.Pointer(point2)),
		uintptr(color1),
		uintptr(color2),
		uintptr(wrapMode),
		uintptr(unsafe.Pointer(lineGradient)))
	return GpStatus(ret)
}

func GdipCreateLineBrushI(point1, point2 *Point, color1, color2 ARGB, wrapMode WrapMode, lineGradient **GpLineGradient) GpStatus {
	ret, _, _ := gdipCreateLineBrushI.Call(
		uintptr(unsafe.Pointer(point1)),
		uintptr(unsafe.Pointer(point2)),
		uintptr(color1),
		uintptr(color2),
		uintptr(wrapMode),
		uintptr(unsafe.Pointer(lineGradient)))
	return GpStatus(ret)
}

func GdipCreateLineBrushFromRect(rect *RectF, color1, color2 ARGB, mode LinearGradientMode, wrapMode WrapMode, lineGradient **GpLineGradient) GpStatus {
	ret, _, _ := gdipCreateLineBrushFromRect.Call(
		uintptr(unsafe.Pointer(rect)),
		uintptr(color1),
		uintptr(color2),
		uintptr(mode),
		uintptr(wrapMode),
		uintptr(unsafe.Pointer(lineGradient)))
	return GpStatus(ret)
}

func GdipCreateLineBrushFromRectI(rect *Rect, color1, color2 ARGB, mode LinearGradientMode, wrapMode WrapMode, lineGradient **GpLineGradient) GpStatus {
	ret, _, _ := gdipCreateLineBrushFromRectI.Call(
		uintptr(unsafe.Pointer(rect)),
		uintptr(color1),
		uintptr(color2),
		uintptr(mode),
		uintptr(wrapMode),
		uintptr(unsafe.Pointer(lineGradient)))
	return GpStatus(ret)
}

func GdipCreateLineBrushFromRectWithAngle(rect *RectF, color1, color2 ARGB, angle float32, isAngleScalable int32, wrapMode WrapMode, lineGradient **GpLineGradient) GpStatus {
	ret, _, _ := gdipCreateLineBrushFromRectWithAngle.Call(
		uintptr(unsafe.Pointer(rect)),
		uintptr(color1),
		uintptr(color2),
		uintptr(math.Float32bits(angle)),
		uintptr(isAngleScalable),
		uintptr(wrapMode),
		uintptr(unsafe.Pointer(lineGradient)))
	return GpStatus(ret)
}

func GdipCreateLineBrushFromRectWithAngleI(rect *Rect, color1, color2 ARGB, angle float32, isAngleScalable int32, wrapMode WrapMode, lineGradient **GpLineGradient) GpStatus {
	ret, _, _ := gdipCreateLineBrushFromRectWithAngleI.Call(
		uintptr(unsafe.Pointer(rect)),
		uintptr(color1),
		uintptr(color2),
		uintptr(math.Float32bits(angle)),
		uintptr(isAngleScalable),
		uintptr(wrapMode),
		uintptr(unsafe.Pointer(lineGradient)))
	return GpStatus(ret)
}

func GdipSetLineColors(brush *GpLineGradient, color1, color2 ARGB) GpStatus {
	ret, _, _ := gdipSetLineColors.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(color1),
		uintptr(color2))
	return GpStatus(ret)
}

func GdipGetLineColors(brush *GpLineGradient, colors *ARGB) GpStatus {
	ret, _, _ := gdipGetLineColors.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(colors)))
	return GpStatus(ret)
}

func GdipSetLinePresetBlend(brush *GpLineGradient, blend *ARGB, positions *float32, count int32) GpStatus {
	ret, _, _ := gdipSetLinePresetBlend.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(blend)),
		uintptr(unsafe.Pointer(positions)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipSetLineLinearBlend(brush *GpLineGradient, focus, scale float32) GpStatus {
	ret, _, _ := gdipSetLineLinearBlend.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(math.Float32bits(focus)),
		uintptr(math.Float32bits(scale)))
	return GpStatus(ret)
}

func GdipSetLineSigmaBlend(brush *GpLineGradient, focus, scale float32) GpStatus {
	ret, _, _ := gdipSetLineSigmaBlend.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(math.Float32bits(focus)),
		uintptr(math.Float32bits(scale)))
	return GpStatus(ret)
}

func GdipSetLineWrapMode(brush *GpLineGradient, wrapMode WrapMode) GpStatus {
	ret, _, _ := gdipSetLineWrapMode.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(wrapMode))
	return GpStatus(ret)
}

func GdipGetLineWrapMode(brush *GpLineGradient, wrapMode *WrapMode) GpStatus {
	ret, _, _ := gdipGetLineWrapMode.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(wrapMode)))
	return GpStatus(ret)
}

func GdipSetLineGammaCorrection(brush *GpLineGradient, useGammaCorrection int32) GpStatus {
	ret, _, _ := gdipSetLineGammaCorrection.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(useGammaCorrection))
	return GpStatus(ret)
}

// Hatch / Texture / PathGradient
func GdipCreateHatchBrush(hatchStyle HatchStyle, forecol, backcol ARGB, brush **GpHatch) GpStatus {
	ret, _, _ := gdipCreateHatchBrush.Call(
		uintptr(hatchStyle),
		uintptr(forecol),
		uintptr(backcol),
		uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}

func GdipCreateTexture(image *GpImage, wrapMode WrapMode, texture **GpTexture) GpStatus {
	ret, _, _ := gdipCreateTexture.Call(
		uintptr(unsafe.Pointer(image)),
		uintptr(wrapMode),
		uintptr(unsafe.Pointer(texture)))
	return GpStatus(ret)
}

func GdipCreatePathGradientFromPath(path *GpPath, polyGradient **GpPathGradient) GpStatus {
	ret, _, _ := gdipCreatePathGradientFromPath.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(unsafe.Pointer(polyGradient)))
	return GpStatus(ret)
}

func GdipSetPathGradientCenterColor(brush *GpPathGradient, colors ARGB) GpStatus {
	ret, _, _ := gdipSetPathGradientCenterColor.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(colors))
	return GpStatus(ret)
}

func GdipSetPathGradientSurroundColorsWithCount(brush *GpPathGradient, color *ARGB, count *int32) GpStatus {
	ret, _, _ := gdipSetPathGradientSurroundColorsWithCount.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(color)),
		uintptr(unsafe.Pointer(count)))
	return GpStatus(ret)
}

func GdipSetPathGradientCenterPoint(brush *GpPathGradient, points *PointF) GpStatus {
	ret, _, _ := gdipSetPathGradientCenterPoint.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(points)))
	return GpStatus(ret)
}

// Font
func GdipCreateFontFromDC(hdc HDC, font **GpFont) GpStatus {
	ret, _, _ := gdipCreateFontFromDC.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(font)))
	return GpStatus(ret)
}

func GdipCreateFont(fontFamily *GpFontFamily, emSize float32, style int32, unit GpUnit, font **GpFont) GpStatus {
	ret, _, _ := gdipCreateFont.Call(
		uintptr(unsafe.Pointer(fontFamily)),
		uintptr(math.Float32bits(emSize)),
		uintptr(style),
		uintptr(unit),
		uintptr(unsafe.Pointer(font)))
	return GpStatus(ret)
}

func GdipDeleteFont(font *GpFont) GpStatus {
	ret, _, _ := gdipDeleteFont.Call(uintptr(unsafe.Pointer(font)))
	return GpStatus(ret)
}

func GdipNewInstalledFontCollection(fontCollection **GpFontCollection) GpStatus {
	ret, _, _ := gdipNewInstalledFontCollection.Call(uintptr(unsafe.Pointer(fontCollection)))
	return GpStatus(ret)
}

func GdipNewPrivateFontCollection(fontCollection **GpFontCollection) GpStatus {
	ret, _, _ := gdipNewPrivateFontCollection.Call(uintptr(unsafe.Pointer(fontCollection)))
	return GpStatus(ret)
}

func GdipDeletePrivateFontCollection(fontCollection **GpFontCollection) GpStatus {
	ret, _, _ := gdipDeletePrivateFontCollection.Call(uintptr(unsafe.Pointer(fontCollection)))
	return GpStatus(ret)
}

func GdipPrivateAddFontFile(fontCollection *GpFontCollection, filename *uint16) GpStatus {
	ret, _, _ := gdipPrivateAddFontFile.Call(
		uintptr(unsafe.Pointer(fontCollection)),
		uintptr(unsafe.Pointer(filename)))
	return GpStatus(ret)
}

func GdipPrivateAddMemoryFont(fontCollection *GpFontCollection, memory unsafe.Pointer, length int32) GpStatus {
	ret, _, _ := gdipPrivateAddMemoryFont.Call(
		uintptr(unsafe.Pointer(fontCollection)),
		uintptr(memory),
		uintptr(length))
	return GpStatus(ret)
}

func GdipGetFontCollectionFamilyCount(fontCollection *GpFontCollection, numFound *int32) GpStatus {
	ret, _, _ := gdipGetFontCollectionFamilyCount.Call(
		uintptr(unsafe.Pointer(fontCollection)),
		uintptr(unsafe.Pointer(numFound)))
	return GpStatus(ret)
}

func GdipGetFontCollectionFamilyList(fontCollection *GpFontCollection, numSought int32, gpfamilies **GpFontFamily, numFound *int32) GpStatus {
	ret, _, _ := gdipGetFontCollectionFamilyList.Call(
		uintptr(unsafe.Pointer(fontCollection)),
		uintptr(numSought),
		uintptr(unsafe.Pointer(gpfamilies)),
		uintptr(unsafe.Pointer(numFound)))
	return GpStatus(ret)
}

func GdipCreateFontFamilyFromName(name *uint16, fontCollection *GpFontCollection, fontFamily **GpFontFamily) GpStatus {
	ret, _, _ := gdipCreateFontFamilyFromName.Call(
		uintptr(unsafe.Pointer(name)),
		uintptr(unsafe.Pointer(fontCollection)),
		uintptr(unsafe.Pointer(fontFamily)))
	return GpStatus(ret)
}

func GdipDeleteFontFamily(fontFamily *GpFontFamily) GpStatus {
	ret, _, _ := gdipDeleteFontFamily.Call(uintptr(unsafe.Pointer(fontFamily)))
	return GpStatus(ret)
}

func GdipGetFamilyName(family *GpFontFamily, name *uint16, language uint16) GpStatus {
	ret, _, _ := gdipGetFamilyName.Call(
		uintptr(unsafe.Pointer(family)),
		uintptr(unsafe.Pointer(name)),
		uintptr(language))
	return GpStatus(ret)
}

func GdipGetGenericFontFamilySansSerif(nativeFamily **GpFontFamily) GpStatus {
	ret, _, _ := gdipGetGenericFontFamilySansSerif.Call(uintptr(unsafe.Pointer(nativeFamily)))
	return GpStatus(ret)
}

func GdipGetGenericFontFamilySerif(nativeFamily **GpFontFamily) GpStatus {
	ret, _, _ := gdipGetGenericFontFamilySerif.Call(uintptr(unsafe.Pointer(nativeFamily)))
	return GpStatus(ret)
}

func GdipGetGenericFontFamilyMonospace(nativeFamily **GpFontFamily) GpStatus {
	ret, _, _ := gdipGetGenericFontFamilyMonospace.Call(uintptr(unsafe.Pointer(nativeFamily)))
	return GpStatus(ret)
}

func GdipGetEmHeight(family *GpFontFamily, style int32, EmHeight *uint16) GpStatus {
	ret, _, _ := gdipGetEmHeight.Call(
		uintptr(unsafe.Pointer(family)),
		uintptr(style),
		uintptr(unsafe.Pointer(EmHeight)))
	return GpStatus(ret)
}

func GdipGetCellAscent(family *GpFontFamily, style int32, CellAscent *uint16) GpStatus {
	ret, _, _ := gdipGetCellAscent.Call(
		uintptr(unsafe.Pointer(family)),
		uintptr(style),
		uintptr(unsafe.Pointer(CellAscent)))
	return GpStatus(ret)
}

func GdipGetCellDescent(family *GpFontFamily, style int32, CellDescent *uint16) GpStatus {
	ret, _, _ := gdipGetCellDescent.Call(
		uintptr(unsafe.Pointer(family)),
		uintptr(style),
		uintptr(unsafe.Pointer(CellDescent)))
	return GpStatus(ret)
}

func GdipGetLineSpacing(family *GpFontFamily, style int32, LineSpacing *uint16) GpStatus {
	ret, _, _ := gdipGetLineSpacing.Call(
		uintptr(unsafe.Pointer(family)),
		uintptr(style),
		uintptr(unsafe.Pointer(LineSpacing)))
	return GpStatus(ret)
}

// StringFormat
func GdipCreateStringFormat(formatAttributes int32, language uint16, format **GpStringFormat) GpStatus {
	ret, _, _ := gdipCreateStringFormat.Call(
		uintptr(formatAttributes),
		uintptr(language),
		uintptr(unsafe.Pointer(format)))
	return GpStatus(ret)
}

func GdipStringFormatGetGenericTypographic(format **GpStringFormat) GpStatus {
	ret, _, _ := gdipStringFormatGetGenericTypographic.Call(uintptr(unsafe.Pointer(format)))
	return GpStatus(ret)
}

func GdipStringFormatGetGenericDefault(format **GpStringFormat) GpStatus {
	ret, _, _ := gdipStringFormatGetGenericDefault.Call(uintptr(unsafe.Pointer(format)))
	return GpStatus(ret)
}

func GdipCloneStringFormat(format *GpStringFormat, newFormat **GpStringFormat) GpStatus {
	ret, _, _ := gdipCloneStringFormat.Call(
		uintptr(unsafe.Pointer(format)),
		uintptr(unsafe.Pointer(newFormat)))
	return GpStatus(ret)
}

func GdipDeleteStringFormat(format *GpStringFormat) GpStatus {
	ret, _, _ := gdipDeleteStringFormat.Call(uintptr(unsafe.Pointer(format)))
	return GpStatus(ret)
}

func GdipSetStringFormatAlign(format *GpStringFormat, align StringAlignment) GpStatus {
	ret, _, _ := gdipSetStringFormatAlign.Call(
		uintptr(unsafe.Pointer(format)),
		uintptr(align))
	return GpStatus(ret)
}

func GdipGetStringFormatAlign(format *GpStringFormat, align *StringAlignment) GpStatus {
	ret, _, _ := gdipGetStringFormatAlign.Call(
		uintptr(unsafe.Pointer(format)),
		uintptr(unsafe.Pointer(align)))
	return GpStatus(ret)
}

func GdipSetStringFormatLineAlign(format *GpStringFormat, align StringAlignment) GpStatus {
	ret, _, _ := gdipSetStringFormatLineAlign.Call(
		uintptr(unsafe.Pointer(format)),
		uintptr(align))
	return GpStatus(ret)
}

func GdipGetStringFormatLineAlign(format *GpStringFormat, align *StringAlignment) GpStatus {
	ret, _, _ := gdipGetStringFormatLineAlign.Call(
		uintptr(unsafe.Pointer(format)),
		uintptr(unsafe.Pointer(align)))
	return GpStatus(ret)
}

func GdipSetStringFormatFlags(format *GpStringFormat, flags int32) GpStatus {
	ret, _, _ := gdipSetStringFormatFlags.Call(
		uintptr(unsafe.Pointer(format)),
		uintptr(flags))
	return GpStatus(ret)
}

func GdipGetStringFormatFlags(format *GpStringFormat, flags *int32) GpStatus {
	ret, _, _ := gdipGetStringFormatFlags.Call(
		uintptr(unsafe.Pointer(format)),
		uintptr(unsafe.Pointer(flags)))
	return GpStatus(ret)
}

func GdipSetStringFormatTrimming(format *GpStringFormat, trimming StringTrimming) GpStatus {
	ret, _, _ := gdipSetStringFormatTrimming.Call(
		uintptr(unsafe.Pointer(format)),
		uintptr(trimming))
	return GpStatus(ret)
}

func GdipGetStringFormatTrimming(format *GpStringFormat, trimming *StringTrimming) GpStatus {
	ret, _, _ := gdipGetStringFormatTrimming.Call(
		uintptr(unsafe.Pointer(format)),
		uintptr(unsafe.Pointer(trimming)))
	return GpStatus(ret)
}

// Path
func GdipCreatePath(brushMode int32, path **GpPath) GpStatus {
	ret, _, _ := gdipCreatePath.Call(uintptr(brushMode), uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipDeletePath(path *GpPath) GpStatus {
	ret, _, _ := gdipDeletePath.Call(uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipResetPath(path *GpPath) GpStatus {
	ret, _, _ := gdipResetPath.Call(uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipStartPathFigure(path *GpPath) GpStatus {
	ret, _, _ := gdipStartPathFigure.Call(uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipAddPathArc(path *GpPath, x, y, width, height, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := gdipAddPathArc.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipAddPathArcI(path *GpPath, x, y, width, height int32, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := gdipAddPathArcI.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipAddPathLine(path *GpPath, x1, y1, x2, y2 float32) GpStatus {
	ret, _, _ := gdipAddPathLine.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(math.Float32bits(x1)),
		uintptr(math.Float32bits(y1)),
		uintptr(math.Float32bits(x2)),
		uintptr(math.Float32bits(y2)))
	return GpStatus(ret)
}

func GdipAddPathLineI(path *GpPath, x1, y1, x2, y2 int32) GpStatus {
	ret, _, _ := gdipAddPathLineI.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(x1),
		uintptr(y1),
		uintptr(x2),
		uintptr(y2))
	return GpStatus(ret)
}

func GdipAddPathRectangle(path *GpPath, x, y, width, height float32) GpStatus {
	ret, _, _ := gdipAddPathRectangle.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipAddPathRectangleI(path *GpPath, x, y, width, height int32) GpStatus {
	ret, _, _ := gdipAddPathRectangleI.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipAddPathEllipse(path *GpPath, x, y, width, height float32) GpStatus {
	ret, _, _ := gdipAddPathEllipse.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipAddPathEllipseI(path *GpPath, x, y, width, height int32) GpStatus {
	ret, _, _ := gdipAddPathEllipseI.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipAddPathPie(path *GpPath, x, y, width, height, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := gdipAddPathPie.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipAddPathPolygon(path *GpPath, points *PointF, count int32) GpStatus {
	ret, _, _ := gdipAddPathPolygon.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipAddPathBezier(path *GpPath, x1, y1, x2, y2, x3, y3, x4, y4 float32) GpStatus {
	ret, _, _ := gdipAddPathBezier.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(math.Float32bits(x1)),
		uintptr(math.Float32bits(y1)),
		uintptr(math.Float32bits(x2)),
		uintptr(math.Float32bits(y2)),
		uintptr(math.Float32bits(x3)),
		uintptr(math.Float32bits(y3)),
		uintptr(math.Float32bits(x4)),
		uintptr(math.Float32bits(y4)))
	return GpStatus(ret)
}

func GdipAddPathString(path *GpPath, str *uint16, length int32, family *GpFontFamily, style int32, emSize float32, layoutRect *RectF, format *GpStringFormat) GpStatus {
	ret, _, _ := gdipAddPathString.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(unsafe.Pointer(str)),
		uintptr(length),
		uintptr(unsafe.Pointer(family)),
		uintptr(style),
		uintptr(math.Float32bits(emSize)),
		uintptr(unsafe.Pointer(layoutRect)),
		uintptr(unsafe.Pointer(format)))
	return GpStatus(ret)
}

func GdipClosePathFigure(path *GpPath) GpStatus {
	ret, _, _ := gdipClosePathFigure.Call(uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipClosePathFigures(path *GpPath) GpStatus {
	ret, _, _ := gdipClosePathFigures.Call(uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

// Image
func GdipGetImageGraphicsContext(image *GpImage, graphics **GpGraphics) GpStatus {
	ret, _, _ := gdipGetImageGraphicsContext.Call(
		uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipLoadImageFromFile(filename *uint16, image **GpImage) GpStatus {
	ret, _, _ := gdipLoadImageFromFile.Call(
		uintptr(unsafe.Pointer(filename)),
		uintptr(unsafe.Pointer(image)))
	return GpStatus(ret)
}

func GdipLoadImageFromStream(stream *ole.IUnknown, image **GpImage) GpStatus {
	ret, _, _ := gdipLoadImageFromStream.Call(
		uintptr(unsafe.Pointer(stream)),
		uintptr(unsafe.Pointer(image)))
	return GpStatus(ret)
}

func GdipSaveImageToFile(image *GpBitmap, filename *uint16, clsidEncoder *ole.GUID, encoderParams *EncoderParameters) GpStatus {
	ret, _, _ := gdipSaveImageToFile.Call(uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(filename)), uintptr(unsafe.Pointer(clsidEncoder)),
		uintptr(unsafe.Pointer(encoderParams)))
	return GpStatus(ret)
}

func GdipGetImageWidth(image *GpImage, width *uint32) GpStatus {
	ret, _, _ := gdipGetImageWidth.Call(uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(width)))
	return GpStatus(ret)
}

func GdipGetImageHeight(image *GpImage, height *uint32) GpStatus {
	ret, _, _ := gdipGetImageHeight.Call(uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(height)))
	return GpStatus(ret)
}

func GdipDisposeImage(image *GpImage) GpStatus {
	ret, _, _ := syscall.Syscall(gdipDisposeImage.Addr(), 1,
		uintptr(unsafe.Pointer(image)),
		0,
		0)
	return GpStatus(ret)
}

// Bitmap
func GdipCreateBitmapFromFile(filename *uint16, bitmap **GpBitmap) GpStatus {
	ret, _, _ := syscall.Syscall(gdipCreateBitmapFromFile.Addr(), 2,
		uintptr(unsafe.Pointer(filename)),
		uintptr(unsafe.Pointer(bitmap)),
		0)
	return GpStatus(ret)
}

func GdipCreateBitmapFromStream(stream *ole.IUnknown, bitmap **GpBitmap) GpStatus {
	ret, _, _ := gdipCreateBitmapFromStream.Call(
		uintptr(unsafe.Pointer(stream)),
		uintptr(unsafe.Pointer(bitmap)))
	return GpStatus(ret)
}

func GdipCreateBitmapFromHBITMAP(hbm HBITMAP, hpal HPALETTE, bitmap **GpBitmap) GpStatus {
	ret, _, _ := syscall.Syscall(gdipCreateBitmapFromHBITMAP.Addr(), 3,
		uintptr(hbm),
		uintptr(hpal),
		uintptr(unsafe.Pointer(bitmap)))
	return GpStatus(ret)
}

func GdipCreateHBITMAPFromBitmap(bitmap *GpBitmap, hbmReturn *HBITMAP, background ARGB) GpStatus {
	ret, _, _ := syscall.Syscall(gdipCreateHBITMAPFromBitmap.Addr(), 3,
		uintptr(unsafe.Pointer(bitmap)),
		uintptr(unsafe.Pointer(hbmReturn)),
		uintptr(background))
	return GpStatus(ret)
}

func GdipCreateBitmapFromScan0(width, height, stride int32, format PixelFormat, scan0 *byte, bitmap **GpBitmap) GpStatus {
	ret, _, _ := gdipCreateBitmapFromScan0.Call(
		uintptr(width),
		uintptr(height),
		uintptr(stride),
		uintptr(format),
		uintptr(unsafe.Pointer(scan0)),
		uintptr(unsafe.Pointer(bitmap)))
	return GpStatus(ret)
}

func GdipBitmapGetPixel(bitmap *GpBitmap, x, y int32, color *ARGB) GpStatus {
	ret, _, _ := gdipBitmapGetPixel.Call(
		uintptr(unsafe.Pointer(bitmap)),
		uintptr(x),
		uintptr(y),
		uintptr(unsafe.Pointer(color)))
	return GpStatus(ret)
}

func GdipBitmapSetPixel(bitmap *GpBitmap, x, y int32, color ARGB) GpStatus {
	ret, _, _ := gdipBitmapSetPixel.Call(
		uintptr(unsafe.Pointer(bitmap)),
		uintptr(x),
		uintptr(y),
		uintptr(color))
	return GpStatus(ret)
}

func GdipBitmapLockBits(bitmap *GpBitmap, rect *Rect, flags uint32, format PixelFormat, lockedBitmapData *BitmapData) GpStatus {
	ret, _, _ := gdipBitmapLockBits.Call(
		uintptr(unsafe.Pointer(bitmap)),
		uintptr(unsafe.Pointer(rect)),
		uintptr(flags),
		uintptr(format),
		uintptr(unsafe.Pointer(lockedBitmapData)))
	return GpStatus(ret)
}

func GdipBitmapUnlockBits(bitmap *GpBitmap, lockedBitmapData *BitmapData) GpStatus {
	ret, _, _ := gdipBitmapUnlockBits.Call(
		uintptr(unsafe.Pointer(bitmap)),
		uintptr(unsafe.Pointer(lockedBitmapData)))
	return GpStatus(ret)
}

func GdipCloneBitmapArea(x, y, width, height float32, format PixelFormat, srcBitmap *GpBitmap, dstBitmap **GpBitmap) GpStatus {
	ret, _, _ := gdipCloneBitmapArea.Call(
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(format),
		uintptr(unsafe.Pointer(srcBitmap)),
		uintptr(unsafe.Pointer(dstBitmap)))
	return GpStatus(ret)
}

func GdipCloneBitmapAreaI(x, y, width, height int32, format PixelFormat, srcBitmap *GpBitmap, dstBitmap **GpBitmap) GpStatus {
	ret, _, _ := gdipCloneBitmapAreaI.Call(
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(format),
		uintptr(unsafe.Pointer(srcBitmap)),
		uintptr(unsafe.Pointer(dstBitmap)))
	return GpStatus(ret)
}

func GdipSaveImageToStream(image *GpBitmap, stream *ole.IUnknown, clsidEncoder *ole.GUID, encoderParams *EncoderParameters) GpStatus {
	ret, _, _ := gdipSaveImageToStream.Call(
		uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(stream)),
		uintptr(unsafe.Pointer(clsidEncoder)),
		uintptr(unsafe.Pointer(encoderParams)))
	return GpStatus(ret)
}
