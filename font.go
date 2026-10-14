package gdiplus

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

type FontCollection struct {
	nativeFontCollection *GpFontCollection
	isPrivate            bool
}

type PrivateFontCollection struct {
	FontCollection
}

type FontFamily struct {
	nativeFontFamily *GpFontFamily
	name             string
}

type Font struct {
	nativeFont *GpFont
	familyName string
	emSize     float32
	style      int32
	unit       GpUnit
	collection *FontCollection
	fallbacks  []string
}

// NewInstalledFontCollection returns the system's installed font collection.
func NewInstalledFontCollection() *FontCollection {
	fc := &FontCollection{}
	GdipNewInstalledFontCollection(&fc.nativeFontCollection)
	return fc
}

// NewPrivateFontCollection creates a private font collection that can load .ttf/.otf files or memory buffers.
func NewPrivateFontCollection() *PrivateFontCollection {
	pfc := &PrivateFontCollection{}
	pfc.isPrivate = true
	GdipNewPrivateFontCollection(&pfc.nativeFontCollection)
	return pfc
}

// AddFontFile loads a font file (.ttf / .otf / .ttc) into the private font collection.
func (pfc *PrivateFontCollection) AddFontFile(filePath string) error {
	if pfc == nil || pfc.nativeFontCollection == nil {
		return fmt.Errorf("PrivateFontCollection is nil")
	}
	pathUTF16, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return err
	}
	status := GdipPrivateAddFontFile(pfc.nativeFontCollection, pathUTF16)
	if status != Ok {
		return fmt.Errorf("GdipPrivateAddFontFile failed for %q: %s", filePath, status.String())
	}
	return nil
}

// AddMemoryFont loads a font directly from a byte slice in memory.
func (pfc *PrivateFontCollection) AddMemoryFont(fontBytes []byte) error {
	if pfc == nil || pfc.nativeFontCollection == nil {
		return fmt.Errorf("PrivateFontCollection is nil")
	}
	if len(fontBytes) == 0 {
		return fmt.Errorf("empty font byte slice")
	}
	status := GdipPrivateAddMemoryFont(pfc.nativeFontCollection, unsafe.Pointer(&fontBytes[0]), int32(len(fontBytes)))
	if status != Ok {
		return fmt.Errorf("GdipPrivateAddMemoryFont failed: %s", status.String())
	}
	return nil
}

// AsFontCollection returns the base *FontCollection to pass into NewFont or NewFontFamily.
func (pfc *PrivateFontCollection) AsFontCollection() *FontCollection {
	if pfc == nil {
		return nil
	}
	return &pfc.FontCollection
}

// Families returns the list of FontFamily names available in this FontCollection.
func (fc *FontCollection) Families() []string {
	if fc == nil || fc.nativeFontCollection == nil {
		return nil
	}
	var count int32
	if GdipGetFontCollectionFamilyCount(fc.nativeFontCollection, &count) != Ok || count <= 0 {
		return nil
	}
	families := make([]*GpFontFamily, count)
	var found int32
	if GdipGetFontCollectionFamilyList(fc.nativeFontCollection, count, &families[0], &found) != Ok {
		return nil
	}
	result := make([]string, 0, found)
	nameBuf := make([]uint16, 64)
	for i := int32(0); i < found; i++ {
		if families[i] == nil {
			continue
		}
		if GdipGetFamilyName(families[i], &nameBuf[0], LANG_NEUTRAL) == Ok {
			result = append(result, syscall.UTF16ToString(nameBuf))
		}
		GdipDeleteFontFamily(families[i])
	}
	return result
}

func (fc *FontCollection) Dispose() {
	if fc != nil && fc.nativeFontCollection != nil {
		if fc.isPrivate {
			GdipDeletePrivateFontCollection(&fc.nativeFontCollection)
		}
		fc.nativeFontCollection = nil
	}
}

// NewFontFamily creates a FontFamily by name.
// Supports CSS-style comma-separated fallback lists (e.g. "Segoe UI, Microsoft YaHei, Arial").
// If none of the requested families exist on the system, it automatically falls back to system default fonts
// so nativeFontFamily is never nil.
func NewFontFamily(familyName string, fontCollection *FontCollection) *FontFamily {
	fm := &FontFamily{name: familyName}
	var nativeColl *GpFontCollection
	if fontCollection != nil {
		nativeColl = fontCollection.nativeFontCollection
	}

	candidates := parseFontFamilyList(familyName)
	for _, candidate := range candidates {
		nameUTF16, err := syscall.UTF16PtrFromString(candidate)
		if err != nil {
			continue
		}
		if GdipCreateFontFamilyFromName(nameUTF16, nativeColl, &fm.nativeFontFamily) == Ok && fm.nativeFontFamily != nil {
			fm.name = candidate
			return fm
		}
		if nativeColl != nil {
			if GdipCreateFontFamilyFromName(nameUTF16, nil, &fm.nativeFontFamily) == Ok && fm.nativeFontFamily != nil {
				fm.name = candidate
				return fm
			}
		}
	}

	// Automatic fallback to system fonts if requested font is not installed
	for _, sysFont := range GetFallbackFonts() {
		nameUTF16, err := syscall.UTF16PtrFromString(sysFont)
		if err != nil {
			continue
		}
		if GdipCreateFontFamilyFromName(nameUTF16, nil, &fm.nativeFontFamily) == Ok && fm.nativeFontFamily != nil {
			fm.name = sysFont
			return fm
		}
	}

	// Ultimate fallback: Generic Sans Serif
	GdipGetGenericFontFamilySansSerif(&fm.nativeFontFamily)
	return fm
}

func (fm *FontFamily) Name() string {
	if fm == nil {
		return ""
	}
	return fm.name
}

func (fm *FontFamily) GetEmHeight(style int32) uint16 {
	if fm == nil || fm.nativeFontFamily == nil {
		return 0
	}
	var val uint16
	GdipGetEmHeight(fm.nativeFontFamily, style, &val)
	return val
}

func (fm *FontFamily) GetCellAscent(style int32) uint16 {
	if fm == nil || fm.nativeFontFamily == nil {
		return 0
	}
	var val uint16
	GdipGetCellAscent(fm.nativeFontFamily, style, &val)
	return val
}

func (fm *FontFamily) GetCellDescent(style int32) uint16 {
	if fm == nil || fm.nativeFontFamily == nil {
		return 0
	}
	var val uint16
	GdipGetCellDescent(fm.nativeFontFamily, style, &val)
	return val
}

func (fm *FontFamily) GetLineSpacing(style int32) uint16 {
	if fm == nil || fm.nativeFontFamily == nil {
		return 0
	}
	var val uint16
	GdipGetLineSpacing(fm.nativeFontFamily, style, &val)
	return val
}

func (fm *FontFamily) Dispose() {
	if fm != nil && fm.nativeFontFamily != nil {
		GdipDeleteFontFamily(fm.nativeFontFamily)
		fm.nativeFontFamily = nil
	}
}

// NewFont creates a Font with automatic family fallback and stores metadata for multi-font glyph fallback.
// You may pass a single font name (e.g. "Segoe UI") or a comma-separated font stack
// (e.g. "Segoe UI, Microsoft YaHei, Segoe UI Emoji").
func NewFont(familyName string, emSize float32, style int32, unit GpUnit, fontCollection *FontCollection) *Font {
	candidates := parseFontFamilyList(familyName)
	primaryName := familyName
	var extraFallbacks []string
	if len(candidates) > 0 {
		primaryName = candidates[0]
		if len(candidates) > 1 {
			extraFallbacks = candidates[1:]
		}
	}

	fontFamily := NewFontFamily(familyName, fontCollection)
	defer fontFamily.Dispose()

	resolvedName := fontFamily.Name()
	if resolvedName == "" {
		resolvedName = primaryName
	}

	f := &Font{
		familyName: resolvedName,
		emSize:     emSize,
		style:      style,
		unit:       unit,
		collection: fontCollection,
		fallbacks:  extraFallbacks,
	}
	GdipCreateFont(fontFamily.nativeFontFamily, emSize, style, unit, &f.nativeFont)
	return f
}

// WithFallbacks attaches custom fallback font family names to this Font instance and returns the Font for chaining.
func (f *Font) WithFallbacks(fallbackFamilies ...string) *Font {
	if f == nil {
		return nil
	}
	f.fallbacks = append([]string{}, fallbackFamilies...)
	return f
}

func (f *Font) FamilyName() string {
	if f == nil {
		return ""
	}
	return f.familyName
}

func (f *Font) Size() float32 {
	if f == nil {
		return 0
	}
	return f.emSize
}

func (f *Font) Style() int32 {
	if f == nil {
		return 0
	}
	return f.style
}

func (f *Font) Unit() GpUnit {
	if f == nil {
		return UnitPixel
	}
	return f.unit
}

func (f *Font) Dispose() {
	if f != nil && f.nativeFont != nil {
		GdipDeleteFont(f.nativeFont)
		f.nativeFont = nil
	}
}

func parseFontFamilyList(familyName string) []string {
	parts := strings.Split(familyName, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	if len(out) == 0 && strings.TrimSpace(familyName) != "" {
		out = append(out, strings.TrimSpace(familyName))
	}
	return out
}
