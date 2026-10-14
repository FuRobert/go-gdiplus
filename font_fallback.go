package gdiplus

import (
	"encoding/binary"
	"sort"
	"strings"
	"sync"
	"syscall"
	"unicode"
	"unsafe"

	"golang.org/x/sys/windows"
)

// GDI32 procedures for accurate OpenType cmap reading and glyph existence checking
var (
	libgdi32               = windows.NewLazySystemDLL("gdi32.dll")
	procCreateCompatibleDC = libgdi32.NewProc("CreateCompatibleDC")
	procDeleteDC           = libgdi32.NewProc("DeleteDC")
	procCreateFontW        = libgdi32.NewProc("CreateFontW")
	procSelectObject       = libgdi32.NewProc("SelectObject")
	procDeleteObject       = libgdi32.NewProc("DeleteObject")
	procGetFontData        = libgdi32.NewProc("GetFontData")
	procGetGlyphIndicesW   = libgdi32.NewProc("GetGlyphIndicesW")
)

const (
	ggiMarkNonexistingGlyphs = 0x0001
	cmapTableTag             = 0x70616D63 // 'cmap' in little-endian ('c','m','a','p')
	gdiError                 = 0xFFFFFFFF
)

type unicodeRange struct {
	start rune
	end   rune
}

type fontCoverage struct {
	ranges    []unicodeRange // sorted, non-overlapping inclusive ranges from OpenType cmap
	hasCmap   bool
	familyUTF [32]uint16
}

var (
	fallbackMu       sync.RWMutex
	autoFontFallback = true
	defaultFallbacks = []string{
		"Segoe UI Emoji",
		"Segoe UI Symbol",
		"Microsoft YaHei UI",
		"Microsoft YaHei",
		"SimSun",
		"Yu Gothic UI",
		"Meiryo",
		"Malgun Gothic",
		"Arial Unicode MS",
		"Segoe UI",
		"Tahoma",
	}

	coverageCache sync.Map // map[string]*fontCoverage (key: lowercase familyName)

	fallbackFontPoolMu sync.RWMutex
	fallbackFontPool   = make(map[fallbackFontKey]*Font)
)

type fallbackFontKey struct {
	family string
	size   uint32 // math.Float32bits(emSize)
	style  int32
	unit   GpUnit
}

// SetAutoFontFallback enables or disables automatic multi-font fallback when drawing/measuring text.
// Enabled by default.
func SetAutoFontFallback(enabled bool) {
	fallbackMu.Lock()
	autoFontFallback = enabled
	fallbackMu.Unlock()
}

// IsAutoFontFallbackEnabled returns whether automatic multi-font fallback is enabled.
func IsAutoFontFallbackEnabled() bool {
	fallbackMu.RLock()
	defer fallbackMu.RUnlock()
	return autoFontFallback
}

// SetFallbackFonts replaces the global fallback font family list.
func SetFallbackFonts(families []string) {
	fallbackMu.Lock()
	cleaned := make([]string, 0, len(families))
	for _, f := range families {
		f = strings.TrimSpace(f)
		if f != "" {
			cleaned = append(cleaned, f)
		}
	}
	if len(cleaned) > 0 {
		defaultFallbacks = cleaned
	}
	fallbackMu.Unlock()
}

// AddFallbackFont prepends or appends a fallback font family to the global fallback list.
func AddFallbackFont(familyName string, prepend bool) {
	familyName = strings.TrimSpace(familyName)
	if familyName == "" {
		return
	}
	fallbackMu.Lock()
	defer fallbackMu.Unlock()
	for _, existing := range defaultFallbacks {
		if strings.EqualFold(existing, familyName) {
			return
		}
	}
	if prepend {
		defaultFallbacks = append([]string{familyName}, defaultFallbacks...)
	} else {
		defaultFallbacks = append(defaultFallbacks, familyName)
	}
}

// GetFallbackFonts returns a copy of the current global fallback font list.
func GetFallbackFonts() []string {
	fallbackMu.RLock()
	defer fallbackMu.RUnlock()
	out := make([]string, len(defaultFallbacks))
	copy(out, defaultFallbacks)
	return out
}

func clearFallbackFontCache() {
	fallbackFontPoolMu.Lock()
	for k, f := range fallbackFontPool {
		if f != nil && f.nativeFont != nil {
			GdipDeleteFont(f.nativeFont)
			f.nativeFont = nil
		}
		delete(fallbackFontPool, k)
	}
	fallbackFontPoolMu.Unlock()
}

func getPooledFallbackFont(familyName string, emSize float32, style int32, unit GpUnit) *Font {
	key := fallbackFontKey{
		family: strings.ToLower(familyName),
		size:   *(*uint32)(unsafe.Pointer(&emSize)),
		style:  style,
		unit:   unit,
	}
	fallbackFontPoolMu.RLock()
	if f, ok := fallbackFontPool[key]; ok && f != nil && f.nativeFont != nil {
		fallbackFontPoolMu.RUnlock()
		return f
	}
	fallbackFontPoolMu.RUnlock()

	fallbackFontPoolMu.Lock()
	defer fallbackFontPoolMu.Unlock()
	if f, ok := fallbackFontPool[key]; ok && f != nil && f.nativeFont != nil {
		return f
	}

	fm := &FontFamily{name: familyName}
	nameUTF16, err := syscall.UTF16PtrFromString(familyName)
	if err != nil {
		return nil
	}
	if GdipCreateFontFamilyFromName(nameUTF16, nil, &fm.nativeFontFamily) != Ok || fm.nativeFontFamily == nil {
		return nil
	}
	defer fm.Dispose()

	f := &Font{
		familyName: familyName,
		emSize:     emSize,
		style:      style,
		unit:       unit,
	}
	if GdipCreateFont(fm.nativeFontFamily, emSize, style, unit, &f.nativeFont) != Ok || f.nativeFont == nil {
		return nil
	}
	fallbackFontPool[key] = f
	return f
}

// FontSupportsRune checks if the specified font family contains a glyph for the given Unicode rune.
func FontSupportsRune(familyName string, r rune) bool {
	if isIgnorableControlRune(r) {
		return true
	}
	cov := getOrLoadFontCoverage(familyName)
	if cov == nil {
		return false
	}
	if cov.hasCmap {
		idx := sort.Search(len(cov.ranges), func(i int) bool {
			return cov.ranges[i].end >= r
		})
		return idx < len(cov.ranges) && cov.ranges[idx].start <= r && r <= cov.ranges[idx].end
	}
	if r > 0xFFFF {
		return false
	}
	return checkGlyphViaGDI(familyName, uint16(r))
}

func isIgnorableControlRune(r rune) bool {
	switch r {
	case '\n', '\r', '\t', ' ', '\u00A0', '\u200B', '\u200C', '\u200D', '\uFE0E', '\uFE0F':
		return true
	}
	return unicode.IsControl(r)
}

func isEmojiModifierOrJoiner(r rune) bool {
	if r == '\u200D' || r == '\uFE0F' || r == '\uFE0E' || r == '\u20E3' {
		return true
	}
	if (r >= 0x1F3FB && r <= 0x1F3FF) || (r >= 0xE0020 && r <= 0xE007F) {
		return true
	}
	return false
}

func isLikelyEmoji(r rune) bool {
	return (r >= 0x1F000 && r <= 0x1FAFF) ||
		(r >= 0x2600 && r <= 0x27BF) ||
		(r >= 0x2300 && r <= 0x23FF) ||
		(r >= 0x2B00 && r <= 0x2BFF) ||
		(r >= 0x1F1E6 && r <= 0x1F1FF)
}

func getOrLoadFontCoverage(familyName string) *fontCoverage {
	key := strings.ToLower(strings.TrimSpace(familyName))
	if key == "" {
		return nil
	}
	if val, ok := coverageCache.Load(key); ok {
		return val.(*fontCoverage)
	}

	cov := loadFontCmapFromGDI(familyName)
	coverageCache.Store(key, cov)
	return cov
}

func loadFontCmapFromGDI(familyName string) *fontCoverage {
	cov := &fontCoverage{}
	utf16Slice, err := syscall.UTF16FromString(familyName)
	if err != nil {
		return cov
	}
	for i := 0; i < len(utf16Slice) && i < 31; i++ {
		cov.familyUTF[i] = utf16Slice[i]
	}

	hdc, _, _ := procCreateCompatibleDC.Call(0)
	if hdc == 0 {
		return cov
	}
	defer procDeleteDC.Call(hdc)

	// CreateFontW with cell height 24 and the face name
	hfont, _, _ := procCreateFontW.Call(
		24, 0, 0, 0, 400, 0, 0, 0,
		1, // DEFAULT_CHARSET
		0, 0, 0, 0,
		uintptr(unsafe.Pointer(&cov.familyUTF[0])),
	)
	if hfont == 0 {
		return cov
	}
	defer procDeleteObject.Call(hfont)

	oldObj, _, _ := procSelectObject.Call(hdc, hfont)
	defer procSelectObject.Call(hdc, oldObj)

	size, _, _ := procGetFontData.Call(hdc, cmapTableTag, 0, 0, 0)
	if size == 0 || size == gdiError || size > 16*1024*1024 {
		return cov
	}

	buf := make([]byte, size)
	readBytes, _, _ := procGetFontData.Call(hdc, cmapTableTag, 0, uintptr(unsafe.Pointer(&buf[0])), size)
	if readBytes != size {
		return cov
	}

	ranges := parseOpenTypeCmap(buf)
	if len(ranges) > 0 {
		cov.ranges = ranges
		cov.hasCmap = true
	}
	return cov
}

func checkGlyphViaGDI(familyName string, ch uint16) bool {
	utf16Name, err := syscall.UTF16PtrFromString(familyName)
	if err != nil {
		return false
	}
	hdc, _, _ := procCreateCompatibleDC.Call(0)
	if hdc == 0 {
		return false
	}
	defer procDeleteDC.Call(hdc)

	hfont, _, _ := procCreateFontW.Call(
		24, 0, 0, 0, 400, 0, 0, 0,
		1, 0, 0, 0, 0,
		uintptr(unsafe.Pointer(utf16Name)),
	)
	if hfont == 0 {
		return false
	}
	defer procDeleteObject.Call(hfont)

	oldObj, _, _ := procSelectObject.Call(hdc, hfont)
	defer procSelectObject.Call(hdc, oldObj)

	var glyphIdx uint16 = 0xFFFF
	ret, _, _ := procGetGlyphIndicesW.Call(
		hdc,
		uintptr(unsafe.Pointer(&ch)),
		1,
		uintptr(unsafe.Pointer(&glyphIdx)),
		ggiMarkNonexistingGlyphs,
	)
	if ret == gdiError {
		return false
	}
	return glyphIdx != 0xFFFF && glyphIdx != 0
}

// parseOpenTypeCmap extracts Unicode character coverage ranges from an OpenType/TrueType 'cmap' table.
// Supports Format 12 (full Unicode including Emoji U+10000..U+10FFFF) and Format 4 (BMP U+0000..U+FFFF).
func parseOpenTypeCmap(data []byte) []unicodeRange {
	if len(data) < 4 {
		return nil
	}
	numTables := int(binary.BigEndian.Uint16(data[2:4]))
	if len(data) < 4+numTables*8 {
		return nil
	}

	var format12Offset uint32
	var format4Offset uint32

	for i := 0; i < numTables; i++ {
		rec := data[4+i*8 : 4+(i+1)*8]
		platformID := binary.BigEndian.Uint16(rec[0:2])
		encodingID := binary.BigEndian.Uint16(rec[2:4])
		offset := binary.BigEndian.Uint32(rec[4:8])
		if int(offset)+4 > len(data) {
			continue
		}
		if platformID == 0 || (platformID == 3 && (encodingID == 1 || encodingID == 10)) {
			format := binary.BigEndian.Uint16(data[offset : offset+2])
			if format == 12 && format12Offset == 0 {
				format12Offset = offset
			} else if format == 4 && format4Offset == 0 {
				format4Offset = offset
			}
		}
	}

	if format12Offset != 0 {
		if ranges := parseCmapFormat12(data[format12Offset:]); len(ranges) > 0 {
			return mergeUnicodeRanges(ranges)
		}
	}
	if format4Offset != 0 {
		if ranges := parseCmapFormat4(data[format4Offset:]); len(ranges) > 0 {
			return mergeUnicodeRanges(ranges)
		}
	}
	return nil
}

func parseCmapFormat12(sub []byte) []unicodeRange {
	if len(sub) < 16 {
		return nil
	}
	numGroups := int(binary.BigEndian.Uint32(sub[12:16]))
	if numGroups <= 0 || numGroups > 200000 || len(sub) < 16+numGroups*12 {
		return nil
	}
	ranges := make([]unicodeRange, 0, numGroups)
	for i := 0; i < numGroups; i++ {
		base := 16 + i*12
		startCode := binary.BigEndian.Uint32(sub[base : base+4])
		endCode := binary.BigEndian.Uint32(sub[base+4 : base+8])
		startGlyph := binary.BigEndian.Uint32(sub[base+8 : base+12])
		if startGlyph == 0 && startCode == 0xFFFF {
			continue
		}
		if startCode <= endCode && endCode <= 0x10FFFF {
			ranges = append(ranges, unicodeRange{start: rune(startCode), end: rune(endCode)})
		}
	}
	return ranges
}

func parseCmapFormat4(sub []byte) []unicodeRange {
	if len(sub) < 16 {
		return nil
	}
	segCountX2 := int(binary.BigEndian.Uint16(sub[6:8]))
	segCount := segCountX2 / 2
	if segCount <= 0 || len(sub) < 16+segCount*8 {
		return nil
	}
	endCodesOffset := 14
	startCodesOffset := 16 + segCount*2
	idDeltaOffset := 16 + segCount*4
	idRangeOffsetPos := 16 + segCount*6

	ranges := make([]unicodeRange, 0, segCount)
	for i := 0; i < segCount; i++ {
		endCode := binary.BigEndian.Uint16(sub[endCodesOffset+i*2 : endCodesOffset+i*2+2])
		startCode := binary.BigEndian.Uint16(sub[startCodesOffset+i*2 : startCodesOffset+i*2+2])
		idDelta := binary.BigEndian.Uint16(sub[idDeltaOffset+i*2 : idDeltaOffset+i*2+2])
		idRangeOffset := binary.BigEndian.Uint16(sub[idRangeOffsetPos+i*2 : idRangeOffsetPos+i*2+2])
		if startCode == 0xFFFF && endCode == 0xFFFF {
			break
		}
		if startCode > endCode {
			continue
		}
		if idRangeOffset == 0 {
			if startCode == endCode && ((startCode+idDelta)&0xFFFF) == 0 {
				continue
			}
			ranges = append(ranges, unicodeRange{start: rune(startCode), end: rune(endCode)})
		} else {
			ranges = append(ranges, unicodeRange{start: rune(startCode), end: rune(endCode)})
		}
	}
	return ranges
}

func mergeUnicodeRanges(ranges []unicodeRange) []unicodeRange {
	if len(ranges) <= 1 {
		return ranges
	}
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].start < ranges[j].start
	})
	merged := make([]unicodeRange, 0, len(ranges))
	cur := ranges[0]
	for i := 1; i < len(ranges); i++ {
		if ranges[i].start <= cur.end+1 {
			if ranges[i].end > cur.end {
				cur.end = ranges[i].end
			}
		} else {
			merged = append(merged, cur)
			cur = ranges[i]
		}
	}
	merged = append(merged, cur)
	return merged
}

// fontRun represents a contiguous segment of text rendered with a single resolved Font.
type fontRun struct {
	text string
	font *Font
}

// textNeedsFallback quickly checks if any rune in text is unsupported by the primary font.
func textNeedsFallback(text string, primaryFont *Font) bool {
	if primaryFont == nil || primaryFont.familyName == "" || !IsAutoFontFallbackEnabled() {
		return false
	}
	for _, r := range text {
		if !isIgnorableControlRune(r) && !FontSupportsRune(primaryFont.familyName, r) {
			return true
		}
	}
	return false
}

// resolveFontForRune finds the best Font for a single rune, starting with primaryFont,
// then primaryFont.fallbacks, then the global defaultFallbackFonts.
func resolveFontForRune(r rune, primaryFont *Font) *Font {
	if primaryFont == nil {
		return nil
	}
	if isIgnorableControlRune(r) || FontSupportsRune(primaryFont.familyName, r) {
		return primaryFont
	}

	var candidates []string
	if isLikelyEmoji(r) {
		candidates = append(candidates, "Segoe UI Emoji", "Segoe UI Symbol")
	}
	candidates = append(candidates, primaryFont.fallbacks...)
	candidates = append(candidates, GetFallbackFonts()...)

	for _, fam := range candidates {
		if strings.EqualFold(fam, primaryFont.familyName) {
			continue
		}
		if FontSupportsRune(fam, r) {
			if fbFont := getPooledFallbackFont(fam, primaryFont.emSize, primaryFont.style, primaryFont.unit); fbFont != nil {
				return fbFont
			}
		}
	}
	return primaryFont
}

// splitLineIntoRuns splits a single line (no '\n') into contiguous fontRun segments.
func splitLineIntoRuns(line string, primaryFont *Font) []fontRun {
	if line == "" {
		return nil
	}
	runes := []rune(line)
	var runs []fontRun
	var sb strings.Builder
	var currentFont *Font

	for i, r := range runes {
		var targetFont *Font
		if i > 0 && currentFont != nil && (r == ' ' || isEmojiModifierOrJoiner(r)) {
			targetFont = currentFont
		} else {
			targetFont = resolveFontForRune(r, primaryFont)
		}

		if currentFont == nil {
			currentFont = targetFont
		} else if targetFont != currentFont {
			runs = append(runs, fontRun{
				text: sb.String(),
				font: currentFont,
			})
			sb.Reset()
			currentFont = targetFont
		}
		sb.WriteRune(r)
	}
	if sb.Len() > 0 {
		runs = append(runs, fontRun{
			text: sb.String(),
			font: currentFont,
		})
	}
	return runs
}

func (g *Graphics) measureRunAdvance(text string, font *Font) (width float32, height float32) {
	if text == "" || font == nil || font.nativeFont == nil {
		return 0, 0
	}
	typo := NewGenericTypographicStringFormat()
	defer typo.Dispose()
	typo.SetFormatFlags(StringFormatFlagsMeasureTrailingSpaces | StringFormatFlagsNoWrap | StringFormatFlagsNoClip)

	layout := &RectF{X: 0, Y: 0, Width: 0, Height: 0}
	bound := &RectF{}
	textUTF16, _ := syscall.UTF16FromString(text)
	GdipMeasureString(g.nativeGraphics, &textUTF16[0], int32(len(textUTF16)-1), font.nativeFont, layout, typo.nativeFormat, bound, nil, nil)
	return bound.Width, bound.Height
}

func (g *Graphics) getPrimaryLineHeight(font *Font) float32 {
	_, h := g.measureRunAdvance("Ag中", font)
	if h <= 0 && font != nil {
		return font.emSize * 1.35
	}
	return h
}

// drawStringFallbackPoint renders multi-line or single-line text at (origin.X, origin.Y) with automatic font fallback.
func (g *Graphics) drawStringFallbackPoint(text string, font *Font, origin *PointF, brush *Brush) {
	lineHeight := g.getPrimaryLineHeight(font)
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")

	typo := NewGenericTypographicStringFormat()
	defer typo.Dispose()
	typo.SetFormatFlags(StringFormatFlagsMeasureTrailingSpaces | StringFormatFlagsNoWrap | StringFormatFlagsNoClip)

	curY := origin.Y
	for _, line := range lines {
		if line == "" {
			curY += lineHeight
			continue
		}
		runs := splitLineIntoRuns(line, font)
		if len(runs) == 1 {
			rect := &RectF{X: origin.X, Y: curY, Width: 0, Height: 0}
			textUTF16, _ := syscall.UTF16FromString(runs[0].text)
			GdipDrawString(g.nativeGraphics, &textUTF16[0], int32(len(textUTF16)-1), runs[0].font.nativeFont, rect, nil, brush.nativeBrush)
			curY += lineHeight
			continue
		}

		curX := origin.X
		for _, run := range runs {
			w, h := g.measureRunAdvance(run.text, run.font)
			yOffset := float32(0)
			if h > 0 && lineHeight > 0 && run.font != font {
				yOffset = (lineHeight - h) * 0.5
			}
			rect := &RectF{X: curX, Y: curY + yOffset, Width: 0, Height: 0}
			textUTF16, _ := syscall.UTF16FromString(run.text)
			GdipDrawString(g.nativeGraphics, &textUTF16[0], int32(len(textUTF16)-1), run.font.nativeFont, rect, typo.nativeFormat, brush.nativeBrush)
			curX += w
		}
		curY += lineHeight
	}
}

type wrappedLine struct {
	runs   []fontRun
	width  float32
	height float32
}

func (g *Graphics) layoutFallbackLines(text string, font *Font, maxWidth float32, noWrap bool) []wrappedLine {
	lineHeight := g.getPrimaryLineHeight(font)
	rawLines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	var result []wrappedLine

	for _, rawLine := range rawLines {
		if rawLine == "" {
			result = append(result, wrappedLine{height: lineHeight})
			continue
		}
		if noWrap || maxWidth <= 0 {
			runs := splitLineIntoRuns(rawLine, font)
			var totalW float32
			for _, r := range runs {
				w, _ := g.measureRunAdvance(r.text, r.font)
				totalW += w
			}
			result = append(result, wrappedLine{runs: runs, width: totalW, height: lineHeight})
			continue
		}

		var curLineBuilder strings.Builder
		var curW float32
		for _, r := range rawLine {
			chStr := string(r)
			rFont := resolveFontForRune(r, font)
			chW, _ := g.measureRunAdvance(chStr, rFont)
			if curLineBuilder.Len() > 0 && curW+chW > maxWidth && !isEmojiModifierOrJoiner(r) {
				lineStr := curLineBuilder.String()
				runs := splitLineIntoRuns(lineStr, font)
				result = append(result, wrappedLine{runs: runs, width: curW, height: lineHeight})
				curLineBuilder.Reset()
				curW = 0
			}
			curLineBuilder.WriteRune(r)
			curW += chW
		}
		if curLineBuilder.Len() > 0 {
			lineStr := curLineBuilder.String()
			runs := splitLineIntoRuns(lineStr, font)
			result = append(result, wrappedLine{runs: runs, width: curW, height: lineHeight})
		}
	}
	return result
}

func (g *Graphics) drawStringFallbackRect(text string, font *Font, rect *RectF, format *StringFormat, brush *Brush) {
	var hAlign, vAlign StringAlignment
	var flags int32
	if format != nil && format.nativeFormat != nil {
		hAlign = format.GetAlignment()
		vAlign = format.GetLineAlignment()
		flags = format.GetFormatFlags()
	}
	noWrap := (flags & StringFormatFlagsNoWrap) != 0
	lines := g.layoutFallbackLines(text, font, rect.Width, noWrap)
	if len(lines) == 0 {
		return
	}

	var totalHeight float32
	for _, l := range lines {
		totalHeight += l.height
	}

	startY := rect.Y
	if rect.Height > 0 {
		switch vAlign {
		case StringAlignmentCenter:
			startY = rect.Y + (rect.Height-totalHeight)*0.5
		case StringAlignmentFar:
			startY = rect.Y + rect.Height - totalHeight
		}
	}

	typo := NewGenericTypographicStringFormat()
	defer typo.Dispose()
	typo.SetFormatFlags(StringFormatFlagsMeasureTrailingSpaces | StringFormatFlagsNoWrap | StringFormatFlagsNoClip)

	curY := startY
	for _, l := range lines {
		if rect.Height > 0 && (flags&StringFormatFlagsNoClip) == 0 && curY > rect.Y+rect.Height {
			break
		}
		startX := rect.X
		if rect.Width > 0 {
			switch hAlign {
			case StringAlignmentCenter:
				startX = rect.X + (rect.Width-l.width)*0.5
			case StringAlignmentFar:
				startX = rect.X + rect.Width - l.width
			}
		}

		curX := startX
		for _, run := range l.runs {
			w, h := g.measureRunAdvance(run.text, run.font)
			yOffset := float32(0)
			if h > 0 && l.height > 0 && run.font != font {
				yOffset = (l.height - h) * 0.5
			}
			rRect := &RectF{X: curX, Y: curY + yOffset, Width: 0, Height: 0}
			textUTF16, _ := syscall.UTF16FromString(run.text)
			GdipDrawString(g.nativeGraphics, &textUTF16[0], int32(len(textUTF16)-1), run.font.nativeFont, rRect, typo.nativeFormat, brush.nativeBrush)
			curX += w
		}
		curY += l.height
	}
}

func (g *Graphics) measureStringFallbackRect(text string, font *Font, layoutRect *RectF, format *StringFormat, boundingBox *RectF, codepointsFitted *int32, linesFilled *int32) {
	if boundingBox == nil {
		return
	}
	var flags int32
	if format != nil && format.nativeFormat != nil {
		flags = format.GetFormatFlags()
	}
	maxWidth := float32(0)
	if layoutRect != nil {
		maxWidth = layoutRect.Width
	}
	noWrap := (flags & StringFormatFlagsNoWrap) != 0
	lines := g.layoutFallbackLines(text, font, maxWidth, noWrap)

	var maxW, totalH float32
	for _, l := range lines {
		if l.width > maxW {
			maxW = l.width
		}
		totalH += l.height
	}
	if layoutRect != nil {
		boundingBox.X = layoutRect.X
		boundingBox.Y = layoutRect.Y
	}
	boundingBox.Width = maxW
	boundingBox.Height = totalH
	if codepointsFitted != nil {
		*codepointsFitted = int32(len([]rune(text)))
	}
	if linesFilled != nil {
		*linesFilled = int32(len(lines))
	}
}
