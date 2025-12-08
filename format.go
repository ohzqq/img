package img

import (
	"encoding"
	"fmt"
	"image"
	"slices"
	"strings"

	"github.com/hhrutter/tiff"
)

var (
	_ encoding.TextUnmarshaler = new(Format)
	_ encoding.TextMarshaler   = Format(0)
)

// Format is an image file format.
type Format int

// Image file formats.
const (
	JPEG Format = iota
	PNG
	GIF
	TIFF
	BMP
	PDF
	WEBP
	HTML
	BASE64
)

var formatExts = [][]string{
	{".jpg", ".jpeg"},
	{".png"},
	{".gif"},
	{".tif", ".tiff"},
	{".bmp"},
	{".pdf"},
	{".webp"},
	{".html"},
	{".b64", ".uue"},
}

func (f Format) String() (format string) {
	defer func() {
		if err := recover(); err != nil {
			format = "unknown"
		}
	}()
	return formatExts[f][0]
}

// FormatFromExtension parses image format from filename extension:
// ".jpg" (or ".jpeg"), ".png", ".gif", ".tif" (or ".tiff"), ".bmp", ".pdf",
// ".b64 (or ".uue") and ".webp" are supported.
func FormatFromExtension(ext string) (Format, error) {
	ext = strings.ToLower(ext)
	for index, exts := range formatExts {
		if slices.Contains(exts, ext) {
			return Format(index), nil
		}
	}
	return -1, image.ErrFormat
}

func (f *Format) UnmarshalText(text []byte) error {
	format, err := FormatFromExtension(string(text))
	if err != nil {
		return err
	}
	*f = format
	return nil
}

func (f Format) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

var (
	_ encoding.TextUnmarshaler = new(TIFFCompression)
	_ encoding.TextMarshaler   = TIFFCompression(0)
)

// TIFFCompression describes the type of compression used in Options.
type TIFFCompression int

// Constants for supported TIFF compression types.
const (
	TIFFUncompressed TIFFCompression = iota
	TIFFDeflate
)

var tiffCompression = []string{
	"none",
	"deflate",
}

func (c TIFFCompression) value() tiff.CompressionType {
	switch c {
	case TIFFDeflate:
		return tiff.Deflate
	}
	return tiff.Uncompressed
}

func (c *TIFFCompression) UnmarshalText(text []byte) error {
	t := strings.ToLower(string(text))
	for index, tt := range tiffCompression {
		if t == tt {
			*c = TIFFCompression(index)
			return nil
		}
	}
	return fmt.Errorf("tiff: unsupported compression: %s", t)
}

func (c TIFFCompression) MarshalText() (b []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
			b = []byte("unknown")
		}
	}()
	ct := tiffCompression[c]
	return []byte(ct), nil
}
