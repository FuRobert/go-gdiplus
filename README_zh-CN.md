# go-gdiplus

[English](./README.md) | **简体中文**

`go-gdiplus` 是面向 Go 语言的 Windows GDI+（`gdiplus.dll`）图形绘制封装库。它同时提供了底层 Flat C API 绑定与面向对象的高级封装接口，适用于 Windows 平台下的 2D 矢量绘图、文本排版与图像处理。

本项目 Fork 自 `shahfarhadreza/go-gdiplus`，扩展了渐变画刷、圆角图形、裁剪辅助方法、私有字体加载、多字体与 Emoji 字形自动回退、纯内存流图像编码以及与 Go 标准库 `image.Image` 的互转支持。

---

## 功能特性

- **多字体与 Emoji 自动回退（Font Fallback）**：通过 `gdi32.dll` 直接解析字体的 OpenType `cmap` 表（支持 Format 4 BMP 与 Format 12 全 Unicode 超平面 `U+0000..U+10FFFF`）并缓存字形覆盖范围。当首选字体缺失某些字符、生僻字、特殊符号或 Emoji 时，自动按字符分段回退至系统可用字体（如 `Segoe UI Emoji`、`Segoe UI Symbol`、`Microsoft YaHei`、`SimSun` 等）进行渲染。
- **渐变与图案画刷**：完整支持 `LinearGradientBrush`（两点渐变、矩形四向模式、任意角度渐变、多色插值渐变 `SetInterpolationColors`）、`PathGradientBrush`、`HatchBrush` 和 `TextureBrush`，并在 `Graphics` 上提供一行代码绘制渐变矩形与圆角渐变矩形的快捷方法。
- **矢量路径、圆角图形与裁剪**：内置 `GraphicsPath` 常用几何方法（`AddRoundedRectangle`、`AddString`、`AddPolygon`、`AddBezier`、`AddEllipse`），以及 `DrawRoundedRectangle`、`FillRoundedRectangle`、`DrawOutlinedString`（描边艺术字）、`DrawImageCircle`（圆形图片裁剪）和 `DrawImageRounded`（圆角图片裁剪）。
- **免安装私有字体集合（`PrivateFontCollection`）**：支持直接从本地文件路径或内存字节切片加载 `.ttf` / `.otf` 字体，无需将字体安装到操作系统。
- **纯内存流图像编码与 `image.Image` 互转**：`ToPngBytes()` 与 `ToJpegBytes(quality)` 基于 Windows COM `IStream`（`CreateStreamOnHGlobal`）直接在内存中完成编码，无需创建磁盘临时文件；支持 `*gdiplus.Bitmap` 与 Go 标准库 `image.Image` 直接互转。
- **幂等安全的资源释放**：所有对象的 `Dispose()` 方法均包含判空与内部指针置零保护，重复调用不会引发崩溃。

---

## 安装

```bash
go get -u github.com/furobert/go-gdiplus
```

在代码中导入：

```go
import gdi "github.com/furobert/go-gdiplus"
```

---

## 快速开始

```go
package main

import (
	"os"
	gdi "github.com/furobert/go-gdiplus"
)

func main() {
	// 1. 初始化 GDI+（进程启动时调用一次即可）
	input := &gdi.GdiplusStartupInput{GdiplusVersion: 1}
	output := &gdi.GdiplusStartupOutput{}
	if gdi.GdiplusStartup(input, output) != gdi.Ok {
		panic("GDI+ 初始化失败")
	}
	defer gdi.GdiplusShutdown()

	// 2. 创建 800x450 画布与 Graphics 上下文
	bmp := gdi.NewBitmap(800, 450, gdi.PixelFormat32bppARGB)
	defer bmp.Dispose()

	g := gdi.NewGraphicsFromBitmap(bmp)
	defer g.Dispose()
	g.EnableHighQuality() // 一键开启抗锯齿、双三次插值与高质量文本渲染

	// 3. 绘制对角线性渐变背景
	cStart, _ := gdi.NewColorFromHex("#4FACFE")
	cEnd, _ := gdi.NewColorFromHex("#00F2FE")
	g.FillLinearGradientRect(0, 0, 800, 450, cStart, cEnd, gdi.LinearGradientModeForwardDiagonal)

	// 4. 绘制半透明圆角矩形与边框
	cardBrush := gdi.NewSolidBrush(gdi.NewColor(255, 255, 255, 230))
	defer cardBrush.Dispose()
	g.FillRoundedRectangle(cardBrush.AsBrush(), 40, 40, 720, 370, 20)

	borderPen := gdi.NewPen(gdi.NewColor(255, 255, 255, 255), 2)
	defer borderPen.Dispose()
	g.DrawRoundedRectangle(borderPen, 40, 40, 720, 370, 20)

	// 5. 绘制包含 Emoji 与多语言字符的文本（自动多字体回退）
	font := gdi.NewFont("Segoe UI, Microsoft YaHei", 32, gdi.FontStyleBold, gdi.UnitPixel, nil)
	defer font.Dispose()

	textBrush := gdi.NewSolidBrush(gdi.NewColorRGB(30, 41, 59))
	defer textBrush.Dispose()

	g.DrawString("Hello GDI+ 🚀 ★ Unicode 与 Emoji 混合渲染", font, gdi.NewPointF(80, 90), textBrush.AsBrush())

	// 6. 在内存中编码为 PNG 字节流并保存
	pngBytes, err := bmp.ToPngBytes()
	if err == nil {
		_ = os.WriteFile("output.png", pngBytes, 0644)
	}
}
```

---

## 使用指南

### 1. 线性渐变与多色渐变（`LinearGradientBrush`）

```go
c1 := gdi.NewColorRGB(67, 233, 123)
c2 := gdi.NewColorRGB(56, 249, 215)

// 双色矩形渐变填充
// 支持模式：
// - gdi.LinearGradientModeHorizontal       (水平：左 -> 右)
// - gdi.LinearGradientModeVertical         (垂直：上 -> 下)
// - gdi.LinearGradientModeForwardDiagonal  (对角：左上 -> 右下)
// - gdi.LinearGradientModeBackwardDiagonal (反对角：右上 -> 左下)
g.FillLinearGradientRect(0, 0, 800, 400, c1, c2, gdi.LinearGradientModeVertical)

// 自定义角度（度数）矩形渐变填充
g.FillLinearGradientRectAngle(0, 0, 800, 400, c1, c2, 135.0)

// 多色（3 色及以上）插值渐变填充
c3 := gdi.NewColorRGB(79, 172, 254)
colors := []*gdi.Color{c1, c2, c3}
positions := []float32{0.0, 0.5, 1.0} // 首项必须为 0.0，末项必须为 1.0
g.FillMultiGradientRect(0, 0, 800, 400, colors, positions, gdi.LinearGradientModeHorizontal)

// 圆角矩形线性渐变填充
g.FillRoundedLinearGradientRect(50, 50, 300, 80, 16, c1, c2, gdi.LinearGradientModeHorizontal)

// 创建独立的 LinearGradientBrush 用于任意几何图形或路径填充
gradBrush := gdi.NewLinearGradientBrush(gdi.NewPointF(0, 0), gdi.NewPointF(0, 400), c1, c2)
defer gradBrush.Dispose()
g.FillEllipse(gradBrush.AsBrush(), 100, 100, 200, 200)
```

### 2. 多字体自动回退与 Emoji 渲染

自动字体回退默认处于开启状态（`AutoFontFallback = true`）。在调用 `g.DrawString`、`g.DrawStringEx`、`g.MeasureString` 或 `g.MeasureStringEx` 时：
- 如果文本中的所有字符在主 `Font` 中均存在字形，则直接调用原生 `GdipDrawString`。
- 如果文本中包含主字体缺失的字符（如 Emoji `😀`、`🔥`、特殊符号 `★` 或扩展汉字），库会自动按字符分段，并从回退字体链中选取包含该字形的字体（保持相同的字号 `emSize`、字重 `style` 和单位 `unit`）完成渲染，同时保留 `StringFormat` 的水平/垂直对齐与自动换行属性。

```go
// 在 NewFont 中直接传入逗号分隔的字体回退栈
font := gdi.NewFont("Arial, Microsoft YaHei, Segoe UI Emoji", 24, gdi.FontStyleRegular, gdi.UnitPixel, nil)
defer font.Dispose()

// 或为单个 Font 实例绑定自定义回退字体
font.WithFallbacks("Segoe UI Emoji", "SimSun")

// 自定义全局默认回退字体链
gdi.AddFallbackFont("Segoe UI Symbol", true) // prepend = true 表示插入到最前
gdi.SetFallbackFonts([]string{"Segoe UI Emoji", "Microsoft YaHei", "Arial"})

// 全局开启或关闭自动字体回退
gdi.SetAutoFontFallback(true)

// 查询某个字体家族是否包含指定 Unicode 字符
supported := gdi.FontSupportsRune("Arial", '😀')
```

### 3. 免安装加载私有字体（`PrivateFontCollection`）

直接从文件或内存字节流加载 `.ttf` / `.otf` 字体文件：

```go
pfc := gdi.NewPrivateFontCollection()
defer pfc.Dispose()

// 从文件路径加载（也可使用 pfc.AddMemoryFont(fontBytes) 从内存 []byte 加载）
_ = pfc.AddFontFile("./fonts/OpenSans-Regular.ttf")

families := pfc.Families()
if len(families) > 0 {
	font := gdi.NewFont(families[0], 28, gdi.FontStyleRegular, gdi.UnitPixel, pfc.AsFontCollection())
	defer font.Dispose()
}
```

### 4. 文本对齐、换行截断与描边文字（`StringFormat` & `DrawOutlinedString`）

```go
// 快速创建水平+垂直双居中的 StringFormat
centerFmt := gdi.NewCenteredStringFormat()
defer centerFmt.Dispose()

// 自定义对齐、禁换行与省略号截断
fmtCustom := gdi.NewStringFormat()
defer fmtCustom.Dispose()
fmtCustom.SetAlignment(gdi.StringAlignmentNear)            // 水平对齐：Near (左) / Center (中) / Far (右)
fmtCustom.SetLineAlignment(gdi.StringAlignmentCenter)      // 垂直对齐：Near (上) / Center (中) / Far (下)
fmtCustom.SetFormatFlags(gdi.StringFormatFlagsNoWrap)      // 禁用自动换行
fmtCustom.SetTrimming(gdi.StringTrimmingEllipsisCharacter) // 超出区域时末尾显示 "..."

rect := gdi.NewRectF(50, 50, 400, 60)
g.DrawStringEx("超出宽度限制时将自动使用省略号截断的长文本内容...", font, rect, fmtCustom, textBrush.AsBrush())

// 使用路径绘制带外框描边的文本
fillBrush := gdi.NewSolidBrush(gdi.NewColorRGB(255, 255, 255))
defer fillBrush.Dispose()
outlinePen := gdi.NewPen(gdi.NewColorRGB(30, 41, 59), 4)
outlinePen.SetLineJoin(gdi.LineJoinRound)
defer outlinePen.Dispose()

g.DrawOutlinedString("Outlined Title", "Segoe UI", gdi.FontStyleBold, 48, gdi.NewRectF(50, 150, 500, 80), nil, fillBrush.AsBrush(), outlinePen)
```

### 5. 圆角图形、图片裁剪与坐标变换

```go
img := gdi.NewImageFromFile("photo.jpg")
defer img.Dispose()

// 在 (60, 60) 处将图片裁剪为直径 120 的圆形并绘制
g.DrawImageCircle(img, 60, 60, 120)

// 将图片裁剪为圆角半径 16 的圆角矩形并绘制
g.DrawImageRounded(img, 220, 60, 200, 120, 16)

// 保存画布状态，围绕 (320, 120) 旋转 15 度后绘制并恢复状态
state := g.Save()
g.RotateAt(15.0, 320, 120)
g.DrawImageRect(img, 220, 60, 200, 120)
g.Restore(state)
```

### 6. 内存字节流读写与 `image.Image` 互转

```go
// 直接从内存字节切片（PNG、JPEG、BMP、GIF）解码创建 Bitmap
bmpFromMem, err := gdi.NewBitmapFromBytes(encodedBytes)
if err == nil {
	defer bmpFromMem.Dispose()
}

// 与 Go 标准库 image.Image 互转
stdImg, err := bmp.ToGoImage()            // *gdi.Bitmap -> *image.NRGBA
bmp2 := gdi.NewBitmapFromGoImage(stdImg)  // image.Image -> *gdi.Bitmap
defer bmp2.Dispose()

// 在内存中直接编码为 PNG 或 JPEG 字节流
pngData, err := bmp.ToPngBytes()
jpgData, err := bmp.ToJpegBytes(85)

// 直接保存到文件
_ = bmp.SaveAsPng("output.png")
_ = bmp.SaveAsJpeg("output.jpg", 90)
```

---

## 资源释放说明

GDI+ 对象底层分配在非托管系统内存中，不受 Go 垃圾回收器自动管理。请在每次创建对象（`NewBitmap`、`NewImageFromFile`、`NewGraphicsFromBitmap`、`NewFont`、`NewPen`、`NewSolidBrush`、`NewLinearGradientBrush`、`NewPath`、`NewStringFormat`、`NewPrivateFontCollection`）后搭配 `defer obj.Dispose()`，或在循环内部使用完毕后显式调用 `.Dispose()`。
