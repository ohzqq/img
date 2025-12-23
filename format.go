package img

import (
	"encoding"
	"fmt"
	"image"
	"io"
	"path/filepath"
	"slices"
	"strings"

	//"github.com/HugoSmits86/nativewebp"
	"github.com/bep/imagemeta"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/gen2brain/webp"
	"github.com/hhrutter/tiff"
	"github.com/samber/lo"
	"github.com/sunshineplan/imgconv"
	"github.com/sunshineplan/pdf"
	"golang.org/x/image/bmp"
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
	URL
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

func (f Format) Save(output string, img image.Image, opts ...EncodeOption) error {
	return NewEncoder(f, opts...).Save(output, img)
}

func (f Format) SaveAll(output string, imgs []image.Image, opts ...EncodeOption) error {
	return NewEncoder(f, opts...).SaveAll(output, imgs)
}

func (f Format) Encode(w io.Writer, img image.Image, opts ...EncodeOption) error {
	return NewEncoder(f, opts...).Encode(w, img)
}

func (f Format) DecodeAnimatedWebP(r io.Reader) (*webp.WEBP, error) {
	return webp.DecodeAll(r)
}

func (f Format) Decode(r io.Reader) (image.Image, error) {
	switch f {
	case PDF:
		img, err := pdf.Decode(r)
		if err != nil {
			return nil, fmt.Errorf("pdf.Decode %w", err)
		}
		return img, nil
	case WEBP:
		img, err := webp.Decode(r)
		if err != nil {
			return nil, fmt.Errorf("this webp is probably animated, use DecodeAnimatedWebP %w", err)
		}
		return img, nil
	case TIFF:
		img, err := tiff.Decode(r)
		if err != nil {
			return img, fmt.Errorf("tiff.Decode %w", err)
		}
		return img, nil
	case BMP:
		img, err := bmp.Decode(r)
		if err != nil {
			return nil, fmt.Errorf("bmp.Decode %w", err)
		}
		return img, nil
	case PNG, JPEG, GIF:
		img, _, err := image.Decode(r)
		if err != nil {
			return nil, fmt.Errorf("error from image.Decode %w", err)
		}
		return img, nil
	}
	return nil, fmt.Errorf("error from f.Decode %w", image.ErrFormat)
}

func (f Format) metaFmt() imagemeta.ImageFormat {
	switch f {
	case JPEG:
		return imagemeta.JPEG
	case PNG:
		return imagemeta.PNG
	case TIFF:
		return imagemeta.TIFF
	case WEBP:
		return imagemeta.WebP
	default:
		return imagemeta.ImageFormatAuto
	}
}

func (f Format) convFmt() imgconv.Format {
	return imgconv.Format(int(f))
}

func (f Format) ImageType() imagetype.ImageType {
	return imagetype.FromString(f.String())
}

// MimeType returns the mimetype of the image format.
func (f Format) MimeType() string {
	return f.ImageType().String()
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

// FormatFromExtension parses image format from filename extension:
// ".jpg" (or ".jpeg"), ".png", ".gif", ".tif" (or ".tiff"), ".bmp", ".pdf",
// ".b64 (or ".uue") and ".webp" are supported.
func FormatFromExtension(ext string) (Format, error) {
	ext = strings.ToLower(ext)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	for index, exts := range formatExts {
		if slices.Contains(exts, ext) {
			return Format(index), nil
		}
	}
	return -1, image.ErrFormat
}

func HasExt(name string) bool {
	return filepath.Ext(name) != ""
}

func Ext(name string) (string, bool) {
	ext := strings.ToLower(filepath.Ext(name))
	return ext, ext != ""
}

func IsValidFormat(name string) bool {
	if !HasExt(name) {
		return false
	}
	return slices.Contains(lo.Flatten(formatExts), filepath.Ext(strings.ToLower(name)))
}

func FormatFromFilename(name string) (Format, error) {
	return FormatFromExtension(filepath.Ext(name))
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
