package gdiplus

import (
	"fmt"
	"syscall"
)

type Image struct {
	nativeImage *GpImage
}

func NewImageFromFile(fileName string) *Image {
	image := &Image{}
	fileNameUTF16, err := syscall.UTF16PtrFromString(fileName)
	if err != nil {
		return image
	}
	GdipLoadImageFromFile(fileNameUTF16, &image.nativeImage)
	return image
}

// NewImageFromBytes loads a GDI+ Image directly from a byte slice in memory (PNG, JPEG, GIF, BMP, etc.).
func NewImageFromBytes(data []byte) (*Image, error) {
	bmp, err := NewBitmapFromBytes(data)
	if err != nil {
		return nil, err
	}
	if bmp == nil || bmp.nativeImage == nil {
		return nil, fmt.Errorf("failed to decode image from bytes")
	}
	return &bmp.Image, nil
}

func (image *Image) IsValid() bool {
	return image != nil && image.nativeImage != nil
}

func (image *Image) GetWidth() (width uint32) {
	if image == nil || image.nativeImage == nil {
		return 0
	}
	GdipGetImageWidth(image.nativeImage, &width)
	return
}

func (image *Image) GetHeight() (height uint32) {
	if image == nil || image.nativeImage == nil {
		return 0
	}
	GdipGetImageHeight(image.nativeImage, &height)
	return
}

func (image *Image) Dispose() {
	if image != nil && image.nativeImage != nil {
		GdipDisposeImage(image.nativeImage)
		image.nativeImage = nil
	}
}
