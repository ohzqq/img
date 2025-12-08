package img

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/HugoSmits86/nativewebp"
	"github.com/hhrutter/tiff"
	"github.com/sunshineplan/pdf"
	"golang.org/x/image/bmp"
)

// Encoder is format option
type Encoder struct {
	Format       Format
	EncodeOption []EncodeOption
}

// EncodeOption sets an optional parameter for the Encode and Save functions.
// https://github.com/disintegration/imaging
type EncodeOption func(*encodeConfig)

type encodeConfig struct {
	Quality               int
	gifNumColors          int
	gifQuantizer          draw.Quantizer
	gifDrawer             draw.Drawer
	pngCompressionLevel   png.CompressionLevel
	tiffCompressionType   TIFFCompression
	webpUseExtendedFormat bool
	background            color.Color
	toBase64              bool
	pages                 []image.Image
}

var defaultEncodeConfig = encodeConfig{
	Quality:             75,
	gifNumColors:        256,
	gifQuantizer:        nil,
	gifDrawer:           nil,
	pngCompressionLevel: png.DefaultCompression,
	tiffCompressionType: TIFFDeflate,
	background:          color.Transparent,
}

// NewEncoder initializes an encoder.
func NewEncoder(format Format, opts ...EncodeOption) *Encoder {
	return &Encoder{
		Format:       format,
		EncodeOption: opts,
	}
}

// Encode writes the image img to w in the specified format (JPEG, PNG, GIF,
// TIFF, BMP, PDF, WEBP, HTML, or BASE64).
func (f *Encoder) Encode(w io.Writer, img image.Image) error {
	cfg := defaultEncodeConfig
	for _, option := range f.EncodeOption {
		option(&cfg)
	}

	if cfg.background != nil {
		i := image.NewNRGBA(img.Bounds())
		draw.Draw(i, i.Bounds(), &image.Uniform{cfg.background}, img.Bounds().Min, draw.Src)
		draw.Draw(i, i.Bounds(), img, img.Bounds().Min, draw.Over)
		img = i
	}

	if f.toBase64 {
		var buf bytes.Buffer
		err := f.encode(&buf, img)
		if err != nil {
			return err
		}
		b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
		_, err = w.Write([]byte(b64))
		if err != nil {
			return err
		}
		return nil
	}

	return f.encode(w, img)
}

func (f *Encoder) encode(w io.Writer, img image.Image) error {
	switch f.Format {
	case JPEG:
		if nrgba, ok := img.(*image.NRGBA); ok && nrgba.Opaque() {
			rgba := &image.RGBA{
				Pix:    nrgba.Pix,
				Stride: nrgba.Stride,
				Rect:   nrgba.Rect,
			}
			return jpeg.Encode(w, rgba, &jpeg.Options{Quality: cfg.Quality})
		}
		return jpeg.Encode(w, img, &jpeg.Options{Quality: cfg.Quality})

	case PNG:
		encoder := png.Encoder{CompressionLevel: cfg.pngCompressionLevel}
		return encoder.Encode(w, img)

	case GIF:
		return gif.Encode(w, img, &gif.Options{
			NumColors: cfg.gifNumColors,
			Quantizer: cfg.gifQuantizer,
			Drawer:    cfg.gifDrawer,
		})

	case TIFF:
		return tiff.Encode(w, img, &tiff.Options{Compression: cfg.tiffCompressionType.value(), Predictor: true})

	case BMP:
		return bmp.Encode(w, img)

	case PDF:
		pages := []image.Image{img}
		pages = append(pages, cfg.pages...)
		return pdf.Encode(w, pages, &pdf.Options{Quality: cfg.Quality})

	case WEBP:
		return nativewebp.Encode(w, img, &nativewebp.Options{UseExtendedFormat: cfg.webpUseExtendedFormat})
	}

	return image.ErrFormat
}

// Quality returns an EncodeOption that sets the output JPEG or PDF quality.
// Quality ranges from 1 to 100 inclusive, higher is better.
func Quality(quality int) EncodeOption {
	return func(c *encodeConfig) {
		c.Quality = quality
	}
}

// GIFNumColors returns an EncodeOption that sets the maximum number of colors
// used in the GIF-encoded image. It ranges from 1 to 256.  Default is 256.
func GIFNumColors(numColors int) EncodeOption {
	return func(c *encodeConfig) {
		c.gifNumColors = numColors
	}
}

// GIFQuantizer returns an EncodeOption that sets the quantizer that is used to produce
// a palette of the GIF-encoded image.
func GIFQuantizer(quantizer draw.Quantizer) EncodeOption {
	return func(c *encodeConfig) {
		c.gifQuantizer = quantizer
	}
}

// GIFDrawer returns an EncodeOption that sets the drawer that is used to convert
// the source image to the desired palette of the GIF-encoded image.
func GIFDrawer(drawer draw.Drawer) EncodeOption {
	return func(c *encodeConfig) {
		c.gifDrawer = drawer
	}
}

// PNGCompressionLevel returns an EncodeOption that sets the compression level
// of the PNG-encoded image. Default is png.DefaultCompression.
func PNGCompressionLevel(level png.CompressionLevel) EncodeOption {
	return func(c *encodeConfig) {
		c.pngCompressionLevel = level
	}
}

// TIFFCompressionType returns an EncodeOption that sets the compression type
// of the TIFF-encoded image. Default is tiff.Deflate.
func TIFFCompressionType(compressionType TIFFCompression) EncodeOption {
	return func(c *encodeConfig) {
		c.tiffCompressionType = compressionType
	}
}

// WEBPUseExtendedFormat returns EncodeOption that determines whether to use extended format
// of the WEBP-encoded image. Default is false.
func WEBPUseExtendedFormat(b bool) EncodeOption {
	return func(c *encodeConfig) {
		c.webpUseExtendedFormat = b
	}
}

// BackgroundColor returns an EncodeOption that sets the background color.
func BackgroundColor(color color.Color) EncodeOption {
	return func(c *encodeConfig) {
		c.background = color
	}
}

// PDFPages returns an EncodeOption that sets multiple pages for pdf conversion.
func PDFPages(pages []image.Image) EncodeOption {
	return func(c *encodeConfig) {
		c.pages = pages
	}
}

// Base64 returns an EncodeOption that encodes the format to Base64.
func Base64(format Format) EncodeOption {
	return func(c *encodeConfig) {
		c.toBase64 = true
		c.Format = format
	}
}
