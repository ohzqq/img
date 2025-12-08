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
	"mime"
	"strings"

	"github.com/HugoSmits86/nativewebp"
	"github.com/hhrutter/tiff"
	"github.com/sunshineplan/pdf"
	"golang.org/x/image/bmp"
)

// Encoder is format option
type Encoder struct {
	Format                Format
	EncodeOption          []EncodeOption
	batch                 bool
	padding               string
	Quality               int
	gifNumColors          int
	gifQuantizer          draw.Quantizer
	gifDrawer             draw.Drawer
	pngCompressionLevel   png.CompressionLevel
	tiffCompressionType   TIFFCompression
	background            color.Color
	pages                 []image.Image
	toBase64              bool
	base64Fmt             Format
	webpUseExtendedFormat bool
	webpAnimation         *nativewebp.Animation
}

// NewEncoder initializes an encoder.
func NewEncoder(format Format, opts ...EncodeOption) *Encoder {
	enc := defaultEncodeConfig
	enc.Format = format
	for _, option := range opts {
		option(enc)
	}
	return enc
}

// Encode writes the image img to w in the specified format (JPEG, PNG, GIF,
// TIFF, BMP, PDF, WEBP, HTML, or BASE64).
func (f *Encoder) Encode(w io.Writer, img image.Image) error {
	if f.background != nil {
		i := image.NewNRGBA(img.Bounds())
		draw.Draw(i, i.Bounds(), &image.Uniform{f.background}, img.Bounds().Min, draw.Src)
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
		switch f.base64Fmt {
		case BASE64:
		case HTML:
			b64 = dataURL(f.Format, b64, true)
		case URL:
			b64 = dataURL(f.Format, b64, false)
		}
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
			return jpeg.Encode(w, rgba, &jpeg.Options{Quality: f.Quality})
		}
		return jpeg.Encode(w, img, &jpeg.Options{Quality: f.Quality})

	case PNG:
		encoder := png.Encoder{CompressionLevel: f.pngCompressionLevel}
		return encoder.Encode(w, img)

	case GIF:
		return gif.Encode(w, img, &gif.Options{
			NumColors: f.gifNumColors,
			Quantizer: f.gifQuantizer,
			Drawer:    f.gifDrawer,
		})

	case TIFF:
		return tiff.Encode(w, img, &tiff.Options{Compression: f.tiffCompressionType.value(), Predictor: true})

	case BMP:
		return bmp.Encode(w, img)

	case PDF:
		pages := []image.Image{img}
		pages = append(pages, f.pages...)
		return pdf.Encode(w, pages, &pdf.Options{Quality: f.Quality})

	case WEBP:
		webpOpts := &nativewebp.Options{UseExtendedFormat: f.webpUseExtendedFormat}
		if len(f.webpAnimation.Images) > 0 {
			return nativewebp.EncodeAll(w, f.webpAnimation, webpOpts)
		}
		return nativewebp.Encode(w, img, webpOpts)
	}

	return image.ErrFormat
}

func dataURL(f Format, b64 string, html bool) string {
	var b strings.Builder
	if html {
		b.WriteString(`<img src="`)
	}
	b.WriteString(`data:`)
	b.WriteString(mime.TypeByExtension(f.String()))
	b.WriteString(`;base64,`)
	b.WriteString(b64)
	if html {
		b.WriteString(`"></img>`)
	}
	return b.String()
}
