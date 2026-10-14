package gdiplus

type StringFormat struct {
	nativeFormat *GpStringFormat
}

func NewStringFormat() *StringFormat {
	format := &StringFormat{}
	GdipCreateStringFormat(0, LANG_NEUTRAL, &format.nativeFormat)
	return format
}

func NewStringFormatWithFlags(flags int32) *StringFormat {
	format := &StringFormat{}
	GdipCreateStringFormat(flags, LANG_NEUTRAL, &format.nativeFormat)
	return format
}

func NewGenericTypographicStringFormat() *StringFormat {
	format := &StringFormat{}
	var base *GpStringFormat
	GdipStringFormatGetGenericTypographic(&base)
	if base != nil {
		GdipCloneStringFormat(base, &format.nativeFormat)
	} else {
		GdipCreateStringFormat(0, LANG_NEUTRAL, &format.nativeFormat)
	}
	return format
}

func NewGenericDefaultStringFormat() *StringFormat {
	format := &StringFormat{}
	var base *GpStringFormat
	GdipStringFormatGetGenericDefault(&base)
	if base != nil {
		GdipCloneStringFormat(base, &format.nativeFormat)
	} else {
		GdipCreateStringFormat(0, LANG_NEUTRAL, &format.nativeFormat)
	}
	return format
}

// NewCenteredStringFormat creates a StringFormat with horizontal and vertical centering.
func NewCenteredStringFormat() *StringFormat {
	sf := NewStringFormat()
	sf.SetAlignment(StringAlignmentCenter)
	sf.SetLineAlignment(StringAlignmentCenter)
	return sf
}

// Clone creates an independent copy of this StringFormat.
func (format *StringFormat) Clone() *StringFormat {
	if format == nil || format.nativeFormat == nil {
		return NewStringFormat()
	}
	clone := &StringFormat{}
	GdipCloneStringFormat(format.nativeFormat, &clone.nativeFormat)
	return clone
}

// SetAlignment sets horizontal text alignment (StringAlignmentNear / Center / Far).
func (format *StringFormat) SetAlignment(align StringAlignment) *StringFormat {
	if format != nil && format.nativeFormat != nil {
		GdipSetStringFormatAlign(format.nativeFormat, align)
	}
	return format
}

// GetAlignment returns horizontal text alignment.
func (format *StringFormat) GetAlignment() StringAlignment {
	if format == nil || format.nativeFormat == nil {
		return StringAlignmentNear
	}
	var align StringAlignment
	GdipGetStringFormatAlign(format.nativeFormat, &align)
	return align
}

// SetLineAlignment sets vertical line alignment (StringAlignmentNear / Center / Far).
func (format *StringFormat) SetLineAlignment(align StringAlignment) *StringFormat {
	if format != nil && format.nativeFormat != nil {
		GdipSetStringFormatLineAlign(format.nativeFormat, align)
	}
	return format
}

// GetLineAlignment returns vertical line alignment.
func (format *StringFormat) GetLineAlignment() StringAlignment {
	if format == nil || format.nativeFormat == nil {
		return StringAlignmentNear
	}
	var align StringAlignment
	GdipGetStringFormatLineAlign(format.nativeFormat, &align)
	return align
}

// SetFormatFlags sets StringFormatFlags (e.g., StringFormatFlagsNoWrap | StringFormatFlagsMeasureTrailingSpaces).
func (format *StringFormat) SetFormatFlags(flags int32) *StringFormat {
	if format != nil && format.nativeFormat != nil {
		GdipSetStringFormatFlags(format.nativeFormat, flags)
	}
	return format
}

// GetFormatFlags returns the current StringFormatFlags.
func (format *StringFormat) GetFormatFlags() int32 {
	if format == nil || format.nativeFormat == nil {
		return 0
	}
	var flags int32
	GdipGetStringFormatFlags(format.nativeFormat, &flags)
	return flags
}

// SetTrimming sets how overflowing text is trimmed (e.g., StringTrimmingEllipsisCharacter).
func (format *StringFormat) SetTrimming(trimming StringTrimming) *StringFormat {
	if format != nil && format.nativeFormat != nil {
		GdipSetStringFormatTrimming(format.nativeFormat, trimming)
	}
	return format
}

// GetTrimming returns the current StringTrimming mode.
func (format *StringFormat) GetTrimming() StringTrimming {
	if format == nil || format.nativeFormat == nil {
		return StringTrimmingNone
	}
	var trimming StringTrimming
	GdipGetStringFormatTrimming(format.nativeFormat, &trimming)
	return trimming
}

func (format *StringFormat) Dispose() {
	if format != nil && format.nativeFormat != nil {
		GdipDeleteStringFormat(format.nativeFormat)
		format.nativeFormat = nil
	}
}
