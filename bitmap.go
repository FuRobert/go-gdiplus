package gdiplus

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

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
	fileNameUTF16, _ := syscall.UTF16PtrFromString(fileName)
	var nativeBitmap *GpBitmap
	GdipCreateBitmapFromFile(fileNameUTF16, &nativeBitmap)
	bitmap.nativeImage = (*GpImage)(nativeBitmap)
	return bitmap
}

func (bitmap *Bitmap) Dispose() {
	GdipDisposeImage(bitmap.nativeImage)
}

func (bitmap *Bitmap) ToPngBytes() ([]byte, error) {
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

	data, err := os.ReadFile(tempFileName)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// quality 0-100。
func (bitmap *Bitmap) ToJpegBytes(quality int) ([]byte, error) {
	if quality < 0 || quality > 100 {
		return nil, fmt.Errorf("JPEG quality must be between 0 and 100, got %d", quality)
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
		fileNameUTF16,
		&ClsidJPEGEncoder,
		encoderParams,
	)

	if status != Ok {
		return nil, fmt.Errorf("GdipSaveImageToFile save JPEG failed: %s (quality: %d)", status.String(), quality)
	}

	data, err := os.ReadFile(tempFileName)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// filePath for example :"C:/images/output.png"
func (bitmap *Bitmap) SaveAsPng(filePath string) error {
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

// filePath for example :"C:/images/output.jpg"
// quality 0-100。
func (bitmap *Bitmap) SaveAsJpeg(filePath string, quality int) error {
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
