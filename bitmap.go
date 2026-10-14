package gdiplus

import (
	"fmt"
	"image"
	stdcolor "image/color"
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

var (
	libole32                   = windows.NewLazySystemDLL("ole32.dll")
	libkernel32                = windows.NewLazySystemDLL("kernel32.dll")
	procCreateStreamOnHGlobal  = libole32.NewProc("CreateStreamOnHGlobal")
	procGetHGlobalFromStream   = libole32.NewProc("GetHGlobalFromStream")
	procGlobalAlloc            = libkernel32.NewProc("GlobalAlloc")
	procGlobalFree             = libkernel32.NewProc("GlobalFree")
	procGlobalLock             = libkernel32.NewProc("GlobalLock")
	procGlobalUnlock           = libkernel32.NewProc("GlobalUnlock")
	procGlobalSize             = libkernel32.NewProc("GlobalSize")
)

const gmemMoveable = 0x0002

type Bitmap struct {
	Image
}

func NewBitmap(width, height int32, format PixelFormat) *Bitmap {
	bitmap := &Bitmap{}
	var nativeBitmap *GpBitmap
	status := GdipCreateBitmapFromScan0(width, height, 0, format, nil, &nativeBitmap)
	if status != Ok {
		log.Panicln(status.String())
	}
	bitmap.nativeImage = (*GpImage)(nativeBitmap)
	return bitmap
}

func NewBitmapEx(width, height, stride int32, format PixelFormat, scan0 *byte) *Bitmap {
	bitmap := &Bitmap{}
	var nativeBitmap *GpBitmap
	GdipCreateBitmapFromScan0(width, height, stride, format, scan0, &nativeBitmap)
	bitmap.nativeImage = (*GpImage)(nativeBitmap)
	return bitmap
}

func NewBitmapFromHBITMAP(hbitmap HBITMAP) *Bitmap {
	bitmap := &Bitmap{}
	var nativeBitmap *GpBitmap
	GdipCreateBitmapFromHBITMAP(hbitmap, 0, &nativeBitmap)
	bitmap.nativeImage = (*GpImage)(nativeBitmap)
	return bitmap
}

func NewBitmapFromFile(fileName string) *Bitmap {
	bitmap := &Bitmap{}
	fileNameUTF16, err := syscall.UTF16PtrFromString(fileName)
	if err != nil {
		return bitmap
	}
	var nativeBitmap *GpBitmap
	GdipCreateBitmapFromFile(fileNameUTF16, &nativeBitmap)
	bitmap.nativeImage = (*GpImage)(nativeBitmap)
	return bitmap
}

// NewBitmapFromBytes creates a Bitmap directly from encoded PNG/JPEG/BMP/GIF bytes in memory.
func NewBitmapFromBytes(data []byte) (*Bitmap, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty image byte slice")
	}
	hGlobal, _, _ := procGlobalAlloc.Call(gmemMoveable, uintptr(len(data)))
	if hGlobal != 0 {
		ptr, _, _ := procGlobalLock.Call(hGlobal)
		if ptr != 0 {
			dst := unsafe.Slice((*byte)(unsafe.Pointer(ptr)), len(data))
			copy(dst, data)
			procGlobalUnlock.Call(hGlobal)

			var stream *ole.IUnknown
			hr, _, _ := procCreateStreamOnHGlobal.Call(hGlobal, 1, uintptr(unsafe.Pointer(&stream)))
			if hr == 0 && stream != nil {
				defer stream.Release()
				var nativeBitmap *GpBitmap
				if GdipCreateBitmapFromStream(stream, &nativeBitmap) == Ok && nativeBitmap != nil {
					return &Bitmap{Image: Image{nativeImage: (*GpImage)(nativeBitmap)}}, nil
				}
			} else {
				procGlobalFree.Call(hGlobal)
			}
		} else {
			procGlobalFree.Call(hGlobal)
		}
	}

	// Fallback via temporary file
	tmpFile, err := os.CreateTemp("", "gdiplus_in_*")
	if err != nil {
		return nil, err
	}
	tmpName := tmpFile.Name()
	_, writeErr := tmpFile.Write(data)
	tmpFile.Close()
	defer os.Remove(tmpName)
	if writeErr != nil {
		return nil, writeErr
	}
	bmp := NewBitmapFromFile(tmpName)
	if bmp.nativeImage == nil {
		return nil, fmt.Errorf("GdipCreateBitmapFromFile failed")
	}
	// Clone into an independent memory bitmap so tmpName can be deleted immediately
	w := int32(bmp.GetWidth())
	h := int32(bmp.GetHeight())
	cloned := bmp.CloneArea(0, 0, w, h, PixelFormat32bppARGB)
	bmp.Dispose()
	return cloned, nil
}

// NewBitmapFromGoImage converts any Go standard library image.Image into a *Bitmap.
func NewBitmapFromGoImage(img image.Image) *Bitmap {
	if img == nil {
		return nil
	}
	bounds := img.Bounds()
	w := int32(bounds.Dx())
	h := int32(bounds.Dy())
	if w <= 0 || h <= 0 {
		return nil
	}
	bmp := NewBitmap(w, h, PixelFormat32bppARGB)
	rect := &Rect{X: 0, Y: 0, Width: w, Height: h}
	var bd BitmapData
	if GdipBitmapLockBits((*GpBitmap)(bmp.nativeImage), rect, ImageLockModeWrite, PixelFormat32bppARGB, &bd) == Ok && bd.Scan0 != nil {
		stride := int(bd.Stride)
		buf := unsafe.Slice((*byte)(bd.Scan0), int(h)*stride)
		for y := 0; y < int(h); y++ {
			row := buf[y*stride : y*stride+int(w)*4]
			for x := 0; x < int(w); x++ {
				r, g, b, a := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
				// Convert from pre-multiplied 16-bit RGBA to non-premultiplied 8-bit BGRA
				var r8, g8, b8, a8 byte
				a8 = byte(a >> 8)
				if a > 0 {
					r8 = byte((r * 0xFFFF / a) >> 8)
					g8 = byte((g * 0xFFFF / a) >> 8)
					b8 = byte((b * 0xFFFF / a) >> 8)
				}
				idx := x * 4
				row[idx+0] = b8
				row[idx+1] = g8
				row[idx+2] = r8
				row[idx+3] = a8
			}
		}
		GdipBitmapUnlockBits((*GpBitmap)(bmp.nativeImage), &bd)
	}
	return bmp
}

// ToGoImage converts the GDI+ Bitmap into a standard library *image.NRGBA.
func (bitmap *Bitmap) ToGoImage() (*image.NRGBA, error) {
	if bitmap == nil || bitmap.nativeImage == nil {
		return nil, fmt.Errorf("bitmap is nil")
	}
	w := int32(bitmap.GetWidth())
	h := int32(bitmap.GetHeight())
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("invalid bitmap dimensions %dx%d", w, h)
	}
	out := image.NewNRGBA(image.Rect(0, 0, int(w), int(h)))
	rect := &Rect{X: 0, Y: 0, Width: w, Height: h}
	var bd BitmapData
	if status := GdipBitmapLockBits((*GpBitmap)(bitmap.nativeImage), rect, ImageLockModeRead, PixelFormat32bppARGB, &bd); status == Ok && bd.Scan0 != nil {
		defer GdipBitmapUnlockBits((*GpBitmap)(bitmap.nativeImage), &bd)
		stride := int(bd.Stride)
		src := unsafe.Slice((*byte)(bd.Scan0), int(h)*stride)
		for y := 0; y < int(h); y++ {
			srcRow := src[y*stride : y*stride+int(w)*4]
			dstRow := out.Pix[y*out.Stride : y*out.Stride+int(w)*4]
			for x := 0; x < int(w); x++ {
				si := x * 4
				// BGRA -> RGBA
				dstRow[si+0] = srcRow[si+2]
				dstRow[si+1] = srcRow[si+1]
				dstRow[si+2] = srcRow[si+0]
				dstRow[si+3] = srcRow[si+3]
			}
		}
		return out, nil
	}

	// Fallback via GetPixel
	for y := int32(0); y < h; y++ {
		for x := int32(0); x < w; x++ {
			c := bitmap.GetPixel(x, y)
			out.SetNRGBA(int(x), int(y), stdcolor.NRGBA{
				R: c.GetR(),
				G: c.GetG(),
				B: c.GetB(),
				A: c.GetA(),
			})
		}
	}
	return out, nil
}

// GetPixel returns the Color of the pixel at (x, y).
func (bitmap *Bitmap) GetPixel(x, y int32) *Color {
	c := &Color{}
	if bitmap != nil && bitmap.nativeImage != nil {
		GdipBitmapGetPixel((*GpBitmap)(bitmap.nativeImage), x, y, &c.Argb)
	}
	return c
}

// SetPixel sets the Color of the pixel at (x, y).
func (bitmap *Bitmap) SetPixel(x, y int32, color *Color) {
	if bitmap != nil && bitmap.nativeImage != nil && color != nil {
		GdipBitmapSetPixel((*GpBitmap)(bitmap.nativeImage), x, y, color.GetValue())
	}
}

// CloneArea clones a rectangular sub-region of the Bitmap.
func (bitmap *Bitmap) CloneArea(x, y, width, height int32, format PixelFormat) *Bitmap {
	out := &Bitmap{}
	if bitmap == nil || bitmap.nativeImage == nil {
		return out
	}
	var dst *GpBitmap
	if GdipCloneBitmapAreaI(x, y, width, height, format, (*GpBitmap)(bitmap.nativeImage), &dst) == Ok {
		out.nativeImage = (*GpImage)(dst)
	}
	return out
}

// LockBits locks a rectangular portion of the bitmap for direct memory access.
func (bitmap *Bitmap) LockBits(rect *Rect, flags uint32, format PixelFormat) (*BitmapData, error) {
	if bitmap == nil || bitmap.nativeImage == nil {
		return nil, fmt.Errorf("bitmap is nil")
	}
	bd := &BitmapData{}
	status := GdipBitmapLockBits((*GpBitmap)(bitmap.nativeImage), rect, flags, format, bd)
	if status != Ok {
		return nil, fmt.Errorf("GdipBitmapLockBits failed: %s", status.String())
	}
	return bd, nil
}

// UnlockBits unlocks a bitmap previously locked with LockBits.
func (bitmap *Bitmap) UnlockBits(bd *BitmapData) error {
	if bitmap == nil || bitmap.nativeImage == nil || bd == nil {
		return fmt.Errorf("bitmap or BitmapData is nil")
	}
	status := GdipBitmapUnlockBits((*GpBitmap)(bitmap.nativeImage), bd)
	if status != Ok {
		return fmt.Errorf("GdipBitmapUnlockBits failed: %s", status.String())
	}
	return nil
}

func (bitmap *Bitmap) Dispose() {
	if bitmap != nil && bitmap.nativeImage != nil {
		GdipDisposeImage(bitmap.nativeImage)
		bitmap.nativeImage = nil
	}
}

// saveToMemoryStream encodes the bitmap directly into a memory byte slice via COM IStream without disk I/O.
func (bitmap *Bitmap) saveToMemoryStream(clsid *ole.GUID, encoderParams *EncoderParameters) ([]byte, error) {
	if bitmap == nil || bitmap.nativeImage == nil {
		return nil, fmt.Errorf("bitmap is nil")
	}
	var stream *ole.IUnknown
	hr, _, _ := procCreateStreamOnHGlobal.Call(0, 1, uintptr(unsafe.Pointer(&stream)))
	if hr != 0 || stream == nil {
		return nil, fmt.Errorf("CreateStreamOnHGlobal failed: 0x%x", hr)
	}
	defer stream.Release()

	status := GdipSaveImageToStream((*GpBitmap)(bitmap.nativeImage), stream, clsid, encoderParams)
	if status != Ok {
		return nil, fmt.Errorf("GdipSaveImageToStream failed: %s", status.String())
	}

	var hGlobal uintptr
	hr, _, _ = procGetHGlobalFromStream.Call(uintptr(unsafe.Pointer(stream)), uintptr(unsafe.Pointer(&hGlobal)))
	if hr != 0 || hGlobal == 0 {
		return nil, fmt.Errorf("GetHGlobalFromStream failed: 0x%x", hr)
	}

	size, _, _ := procGlobalSize.Call(hGlobal)
	if size == 0 {
		return nil, fmt.Errorf("GlobalSize returned 0")
	}

	ptr, _, _ := procGlobalLock.Call(hGlobal)
	if ptr == 0 {
		return nil, fmt.Errorf("GlobalLock failed")
	}
	defer procGlobalUnlock.Call(hGlobal)

	out := make([]byte, size)
	copy(out, unsafe.Slice((*byte)(unsafe.Pointer(ptr)), size))
	return out, nil
}

func (bitmap *Bitmap) ToPngBytes() ([]byte, error) {
	if bitmap == nil || bitmap.nativeImage == nil {
		return nil, fmt.Errorf("bitmap is nil")
	}
	if data, err := bitmap.saveToMemoryStream(&ClsidPNGEncoder, nil); err == nil && len(data) > 0 {
		return data, nil
	}

	// Fallback to temp file if memory stream fails
	tempFile, err := os.CreateTemp("", "gdiplus_temp_*.png")
	if err != nil {
		return nil, err
	}
	tempFileName := tempFile.Name()
	if err := tempFile.Close(); err != nil {
		os.Remove(tempFileName)
		return nil, err
	}
	defer os.Remove(tempFileName)

	fileNameUTF16, err := syscall.UTF16PtrFromString(tempFileName)
	if err != nil {
		return nil, err
	}

	status := GdipSaveImageToFile(
		(*GpBitmap)(bitmap.nativeImage),
		fileNameUTF16,
		&ClsidPNGEncoder,
		nil,
	)
	if status != Ok {
		return nil, fmt.Errorf("%s", status.String())
	}

	return os.ReadFile(tempFileName)
}

// ToJpegBytes encodes the bitmap to JPEG bytes with quality in [0, 100].
func (bitmap *Bitmap) ToJpegBytes(quality int) ([]byte, error) {
	if bitmap == nil || bitmap.nativeImage == nil {
		return nil, fmt.Errorf("bitmap is nil")
	}
	if quality < 0 || quality > 100 {
		return nil, fmt.Errorf("JPEG quality must be between 0 and 100, got %d", quality)
	}

	qualityValue := uint32(quality)
	encoderParams := &EncoderParameters{
		Count: 1,
		Parameter: [1]EncoderParameter{
			{
				Guid:           EncoderQuality,
				NumberOfValues: 1,
				TypeAPI:        EncoderParameterValueTypeLong,
				Value:          uintptr(unsafe.Pointer(&qualityValue)),
			},
		},
	}

	if data, err := bitmap.saveToMemoryStream(&ClsidJPEGEncoder, encoderParams); err == nil && len(data) > 0 {
		return data, nil
	}

	tempFile, err := os.CreateTemp("", "gdiplus_temp_*.jpg")
	if err != nil {
		return nil, err
	}
	tempFileName := tempFile.Name()
	if err := tempFile.Close(); err != nil {
		os.Remove(tempFileName)
		return nil, err
	}
	defer os.Remove(tempFileName)

	fileNameUTF16, err := syscall.UTF16PtrFromString(tempFileName)
	if err != nil {
		return nil, err
	}

	status := GdipSaveImageToFile(
		(*GpBitmap)(bitmap.nativeImage),
		fileNameUTF16,
		&ClsidJPEGEncoder,
		encoderParams,
	)
	if status != Ok {
		return nil, fmt.Errorf("GdipSaveImageToFile save JPEG failed: %s (quality: %d)", status.String(), quality)
	}

	return os.ReadFile(tempFileName)
}

// SaveAsPng saves the bitmap to a PNG file (e.g. "C:/images/output.png").
func (bitmap *Bitmap) SaveAsPng(filePath string) error {
	if bitmap == nil || bitmap.nativeImage == nil {
		return fmt.Errorf("bitmap is nil")
	}
	filePathUTF16, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return err
	}

	status := GdipSaveImageToFile(
		(*GpBitmap)(bitmap.nativeImage),
		filePathUTF16,
		&ClsidPNGEncoder,
		nil,
	)
	if status != Ok {
		return fmt.Errorf("GdipSaveImageToFile failed to save PNG to '%s': %s", filePath, status.String())
	}
	return nil
}

// SaveAsJpeg saves the bitmap to a JPEG file with quality in [0, 100].
func (bitmap *Bitmap) SaveAsJpeg(filePath string, quality int) error {
	if bitmap == nil || bitmap.nativeImage == nil {
		return fmt.Errorf("bitmap is nil")
	}
	if quality < 0 || quality > 100 {
		return fmt.Errorf("JPEG quality must be between 0 and 100, got %d", quality)
	}

	filePathUTF16, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return fmt.Errorf("transform file path to UTF16 failed: %w", err)
	}

	qualityValue := uint32(quality)
	encoderParams := &EncoderParameters{
		Count: 1,
		Parameter: [1]EncoderParameter{
			{
				Guid:           EncoderQuality,
				NumberOfValues: 1,
				TypeAPI:        EncoderParameterValueTypeLong,
				Value:          uintptr(unsafe.Pointer(&qualityValue)),
			},
		},
	}

	status := GdipSaveImageToFile(
		(*GpBitmap)(bitmap.nativeImage),
		filePathUTF16,
		&ClsidJPEGEncoder,
		encoderParams,
	)
	if status != Ok {
		return fmt.Errorf("GdipSaveImageToFile save JPEG failed: %s (quality: %d)", status.String(), quality)
	}
	return nil
}
