# go-gdiplus

[English Documentation](./README.md) | **中文文档**

`go-gdiplus` 是面向 Go 语言的轻量、高性能 Windows GDI+ 绘图封装库。它既完整保留了底层 `gdiplus.dll` 的 Flat C API，又提供了符合 Go 直觉的面向对象高级封装，专为长图生成、数据卡片与看板渲染、报表输出、服务端自动化批量制图等高并发绘图场景打造。

Fork 自 `shahfarhadreza/go-gdiplus`，由 `furobert` 大幅重构扩展：新增**系统多字体与 Emoji/特殊符号自动回退引擎**、**一键线性/多色/角度渐变画刷**、**圆角矩形与圆形头像裁剪**、**私有字体文件（TTF/OTF）免安装加载**、**纯内存流（COM IStream）零临时文件 PNG/JPEG 编码**以及**防重复释放安全保护**。

---

## 核心特性

1. **多字体与 Emoji / 特殊符号自动回退（Font Fallback）**
   * 彻底解决原版 GDI+ 单一字体在绘制自定义艺术字或装饰字体时遇到 Emoji（`😀`, `🔥`, `📊`）、特殊符号（`★`, `♪`, `♡`）、多国语言扩展字符或生僻汉字显示为空白方框的问题。
   * 底层直接通过 `gdi32.dll` 解析字体 OpenType `cmap` 表（支持 Format 4 BMP 与 Format 12 全 Unicode 超平面）并建立内存缓存（单字符判定耗时 `< 20ns`），按字符自动切换至 `Segoe UI Emoji`、`Segoe UI Symbol`、`Microsoft YaHei` 等系统可用字体分段渲染。
   * 支持类 CSS 逗号分隔字体栈（如 `"CustomFont, Microsoft YaHei, Segoe UI Emoji"`）。即使首选字体在目标机器上未安装，也会自动降级为系统可用字体，绝不产生空指针。
2. **原生渐变画刷与一行代码画渐变背景**
   * 告别逐行循环 `NewPen` + `DrawLineI` 的繁琐写法！内置 `LinearGradientBrush`、`PathGradientBrush`、`HatchBrush`、`TextureBrush`。
   * 提供 `g.FillLinearGradientRect`、`g.FillLinearGradientRectAngle`、`g.FillMultiGradientRect`、`g.FillRoundedLinearGradientRect` 等快捷方法，一行代码完成水平、垂直、对角、任意角度及多色插值渐变。
3. **圆角卡片、圆形头像裁剪与描边艺术字**
   * 内置 `FillRoundedRectangle`、`DrawRoundedRectangle`、`FillCircle`、`DrawCircle`。
   * 内置 `DrawImageCircle`（一行代码绘制圆形头像/缩略图）与 `DrawImageRounded`（圆角矩形图片裁剪）。
   * 内置 `DrawOutlinedString`，基于 `GraphicsPath.AddString` 实现高质量描边+填充艺术字。
4. **纯内存流图片读写与 Go `image.Image` 互转**
   * `ToPngBytes()` 与 `ToJpegBytes(quality)` 优先通过 Windows COM `IStream`（`CreateStreamOnHGlobal`）直接在内存中完成编码，无需在磁盘上创建临时文件。
   * 支持 `NewBitmapFromBytes(data)` 直接从内存字节切片加载图片，以及 `NewBitmapFromGoImage` / `ToGoImage` 与 Go 标准库无缝互转。
5. **私有字体免安装加载（`PrivateFontCollection`）**
   * 支持直接从资产目录加载 `.ttf` / `.otf` 字体文件或内存字节流，无需将字体安装到系统 `C:\Windows\Fonts`。
6. **幂等安全的资源销毁（`Dispose`）**
   * 所有对象的 `Dispose()` 方法均做判空与指针置零保护，重复调用 `Dispose()` 绝不会引发底层 C 崩溃。

---

## 安装

```bash
go get -u github.com/furobert/go-gdiplus
```

导入包：

```go
import gdi "github.com/furobert/go-gdiplus"
```

---

## 快速上手示例

```go
package main

import (
	"os"
	gdi "github.com/furobert/go-gdiplus"
)

func main() {
	// 1. 初始化 GDI+（整个进程生命周期只需初始化一次）
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
	g.EnableHighQuality() // 一键开启抗锯齿、双三次插值与 ClearType 文本渲染

	// 3. 一行代码绘制渐变背景（左上 -> 右下对角渐变）
	cStart, _ := gdi.NewColorFromHex("#FFDEE9")
	cEnd, _ := gdi.NewColorFromHex("#B5FFFC")
	g.FillLinearGradientRect(0, 0, 800, 450, cStart, cEnd, gdi.LinearGradientModeForwardDiagonal)

	// 4. 绘制半透明白色圆角卡片与边框
	cardBrush := gdi.NewSolidBrush(gdi.NewColor(255, 255, 255, 220))
	defer cardBrush.Dispose()
	g.FillRoundedRectangle(cardBrush.AsBrush(), 40, 40, 720, 370, 24)

	borderPen := gdi.NewPen(gdi.NewColor(255, 182, 193, 255), 3)
	defer borderPen.Dispose()
	g.DrawRoundedRectangle(borderPen, 40, 40, 720, 370, 24)

	// 5. 绘制包含 Emoji、特殊符号与多语言扩展字符的文本（自动多字体回退！）
	// 即使主字体为字符集不全的自定义美术字体，缺失的 Emoji 与符号也会自动用系统字体补齐
	font := gdi.NewFont("CustomDisplayFont, Microsoft YaHei", 32, gdi.FontStyleBold, gdi.UnitPixel, nil)
	defer font.Dispose()

	textBrush := gdi.NewSolidBrush(gdi.NewColorRGB(50, 50, 70))
	defer textBrush.Dispose()

	g.DrawString("📊 综合数据指标分析 🔥 ★ 100% 达成！", font, gdi.NewPointF(80, 90), textBrush.AsBrush())

	// 6. 导出为 JPEG 字节流（纯内存编码）或直接保存为文件
	jpegBytes, err := bmp.ToJpegBytes(90)
	if err == nil {
		_ = os.WriteFile("card.jpg", jpegBytes, 0644)
	}
}
```

---

## 核心功能详解

### 1. 渐变背景与渐变画刷（`LinearGradientBrush`）

不再需要手动写 `for y := 0; y < height; y++` 循环画线！你可以使用 `Graphics` 上的快捷方法或手动创建 `LinearGradientBrush`：

```go
c1 := gdi.NewColorRGB(255, 210, 230)
c2 := gdi.NewColorRGB(190, 225, 255)

// 方式 A：一行代码填充矩形渐变
// 支持模式：
// - gdi.LinearGradientModeHorizontal       (水平：左 -> 右)
// - gdi.LinearGradientModeVertical         (垂直：上 -> 下)
// - gdi.LinearGradientModeForwardDiagonal  (对角：左上 -> 右下)
// - gdi.LinearGradientModeBackwardDiagonal (反对角：右上 -> 左下)
g.FillLinearGradientRect(0, 0, 1000, 600, c1, c2, gdi.LinearGradientModeVertical)

// 方式 B：一行代码按任意角度（如 45° 或 135°）填充矩形渐变
g.FillLinearGradientRectAngle(0, 0, 1000, 600, c1, c2, 135.0)

// 方式 C：一行代码填充多色（3 色及以上）平滑渐变
c3 := gdi.NewColorRGB(240, 255, 220)
colors := []*gdi.Color{c1, c2, c3}
positions := []float32{0.0, 0.5, 1.0} // 首项必须为 0.0，末项必须为 1.0
g.FillMultiGradientRect(0, 0, 1000, 600, colors, positions, gdi.LinearGradientModeVertical)

// 方式 D：一行代码填充圆角矩形渐变（适合按钮、进度条、徽章）
g.FillRoundedLinearGradientRect(50, 50, 300, 80, 20, c1, c2, gdi.LinearGradientModeHorizontal)

// 方式 E：直接创建 LinearGradientBrush 对象用于任意路径或文字填充
gradBrush := gdi.NewLinearGradientBrush(gdi.NewPointF(0, 0), gdi.NewPointF(0, 600), c1, c2)
defer gradBrush.Dispose()
g.FillEllipse(gradBrush.AsBrush(), 100, 100, 200, 200)
```

### 2. 多字体回退与 Emoji / 特殊符号支持

在默认情况下，`go-gdiplus` 已全局开启自动字体回退（`AutoFontFallback = true`）。当调用 `g.DrawString`、`g.DrawStringEx`、`g.MeasureString`、`g.MeasureStringEx` 时：
* 如果整段文本的所有字符都在当前 `Font` 中存在字形，则直接走原生 `GdipDrawString` 零开销快速通道；
* 如果文本中包含当前字体缺失的字符（如 `😀`、`⏳`、`★`、生僻字等），库会自动按字符检索回退链并用匹配的字体（保持相同的字号 `emSize`、字重 `style` 和单位 `unit`）分段绘制，同时完整保留 `StringFormat` 的水平居中/右对齐、垂直居中与自动换行逻辑。

```go
// 1. 方式一：在 NewFont 中直接传入逗号分隔的优先回退栈
font := gdi.NewFont("CustomFont, Microsoft YaHei, Segoe UI Emoji", 28, gdi.FontStyleRegular, gdi.UnitPixel, nil)
defer font.Dispose()

// 2. 方式二：为单个 Font 实例指定专属回退字体列表
font.WithFallbacks("Segoe UI Emoji", "SimSun")

// 3. 方式三：全局添加或修改默认回退字体链
gdi.AddFallbackFont("HarmonyOS Sans SC", true) // true 表示插入到最高优先级

// 4. 检测某个字体是否包含指定字符
supported := gdi.FontSupportsRune("CustomFont", '😀') // 返回 false
```

### 3. 免安装加载本地 `.ttf` / `.otf` 字体文件（`PrivateFontCollection`）

```go
pfc := gdi.NewPrivateFontCollection()
defer pfc.Dispose()

// 从文件加载自定义字体（也可使用 pfc.AddMemoryFont(ttfBytes) 从内存加载）
_ = pfc.AddFontFile("./assets/fonts/CustomFont.ttf")

// 查看加载成功的字体家族名称
families := pfc.Families() // 例如 ["CustomFont"]

// 使用该私有字体集合创建 Font
font := gdi.NewFont(families[0], 36, gdi.FontStyleRegular, gdi.UnitPixel, pfc.AsFontCollection())
defer font.Dispose()
```

### 4. 文本对齐、换行、省略号与描边字（`StringFormat` & `DrawOutlinedString`）

```go
// 快速创建水平+垂直双居中的 StringFormat
centerFmt := gdi.NewCenteredStringFormat()
defer centerFmt.Dispose()

// 或者自定义对齐与超出省略号
fmtCustom := gdi.NewStringFormat()
defer fmtCustom.Dispose()
fmtCustom.SetAlignment(gdi.StringAlignmentNear)                // 水平左对齐 (Near / Center / Far)
fmtCustom.SetLineAlignment(gdi.StringAlignmentCenter)          // 垂直居中
fmtCustom.SetFormatFlags(gdi.StringFormatFlagsNoWrap)          // 禁止自动换行
fmtCustom.SetTrimming(gdi.StringTrimmingEllipsisCharacter)     // 超出宽度时自动显示 "..."

rect := gdi.NewRectF(50, 50, 400, 60)
g.DrawStringEx("这是一段很长很长的文本内容...", font, rect, fmtCustom, textBrush.AsBrush())

// 绘制带描边的艺术字（外框 Pen + 内部 Brush）
fillBrush := gdi.NewSolidBrush(gdi.NewColorRGB(255, 255, 255))
defer fillBrush.Dispose()
outlinePen := gdi.NewPen(gdi.NewColorRGB(255, 105, 180), 4)
outlinePen.SetLineJoin(gdi.LineJoinRound) // 圆角描边防尖刺
defer outlinePen.Dispose()

g.DrawOutlinedString("GRAPHICS ENGINE", "Microsoft YaHei", gdi.FontStyleBold, 48, gdi.NewRectF(80, 200, 600, 80), nil, fillBrush.AsBrush(), outlinePen)
```

### 5. 圆形头像、圆角图片裁剪与坐标变换

```go
// 从本地文件或内存字节加载图像
avatar := gdi.NewImageFromFile("avatar.jpg")
defer avatar.Dispose()

// 一行代码绘制圆形图像（自动抗锯齿裁剪并还原画布裁剪区）
g.DrawImageCircle(avatar, 60, 60, 120) // 在 (60, 60) 处绘制直径 120 的圆形图像

// 一行代码绘制圆角矩形图片
g.DrawImageRounded(avatar, 220, 60, 200, 120, 16) // 圆角半径 16

// 围绕指定中心点旋转绘制
state := g.Save()
g.RotateAt(15.0, 320, 120) // 围绕 (320, 120) 顺时针旋转 15 度
g.DrawImageRect(avatar, 220, 60, 200, 120)
g.Restore(state)
```

### 6. 内存字节流读写与 Go 标准库 `image.Image` 互转

```go
// 1. 直接从网络下载或内存中的 []byte 创建 Bitmap（零临时文件）
bmpFromNet, err := gdi.NewBitmapFromBytes(downloadedBytes)
if err == nil {
	defer bmpFromNet.Dispose()
}

// 2. 与 Go 标准库 image.Image 互转
goImg, err := bmp.ToGoImage()              // *gdi.Bitmap -> *image.NRGBA
bmp2 := gdi.NewBitmapFromGoImage(goImg)    // image.Image -> *gdi.Bitmap
defer bmp2.Dispose()

// 3. 导出为 PNG / JPEG 内存字节流（基于 COM IStream 纯内存编码）
pngBytes, err := bmp.ToPngBytes()
jpgBytes, err := bmp.ToJpegBytes(85)
```

---

## 资源释放最佳实践

由于 GDI+ 对象分配在 Windows 非托管（Unmanaged）内存中，不受 Go GC 自动回收管理，请务必遵循以下习惯：
1. 每创建一个 `*Bitmap`、`*Image`、`*Graphics`、`*Font`、`*FontFamily`、`*Pen`、`*SolidBrush`、`*LinearGradientBrush`、`*GraphicsPath`、`*StringFormat`、`*PrivateFontCollection`，立即在下一行写上 `defer xxx.Dispose()`。
2. 避免在高频 `for` 循环内部创建未释放的临时对象；若必须在循环内创建，请在单次迭代末尾显式调用 `.Dispose()`（本库的 `.Dispose()` 是幂等安全的，重复调用不会报错）。