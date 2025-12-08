package img

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime"
	"os"
	"path/filepath"
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

// Save saves image according to the encoder
// https://github.com/sunshineplan/imgconv
func (enc *Encoder) Save(output string, base image.Image) error {
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	return enc.Encode(f, base)
}

// SaveAll saves images according to the encoder
func (enc *Encoder) SaveAll(output string, padding string, images []image.Image) error {
	enc.padding = padding
	enc.batch = true
	enc.pages = images
	ext := filepath.Ext(output)
	dir, name := filepath.Split(output)
	base := strings.TrimSuffix(name, ext)
	if enc.batch {
		for i, img := range enc.pages {
			n := fmt.Sprintf(base+enc.padding+ext, i)
			f, err := os.Create(filepath.Join(dir, n))
			if err != nil {
				return err
			}
			defer f.Close()
			err = enc.Encode(f, img)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Animate creates an animated WEBP according to the encoder
func (enc *Encoder) Animate(output string, images []image.Image) error {
	enc.Format = WEBP
	enc.webpAnimation.Images = images
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	return enc.animatedWebp(f)
}

// Animate creates an animated WEBP according to the encoder
func (enc *Encoder) AnimatedWebp(output string, images []string) error {
	enc.Format = WEBP
	//enc.webpAnimation.Images = images
	for _, file := range images {
		img, err := Open(file)
		if err != nil {
			return err
		}
		enc.webpAnimation.Images = append(enc.webpAnimation.Images, img)
	}
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	return enc.animatedWebp(f)
}

// Encode writes the image img to w in the specified format (JPEG, PNG, GIF,
// TIFF, BMP, PDF, WEBP, HTML, or BASE64).
func (enc *Encoder) Encode(w io.Writer, img image.Image) error {
	if enc.background != nil {
		i := image.NewNRGBA(img.Bounds())
		draw.Draw(i, i.Bounds(), &image.Uniform{enc.background}, img.Bounds().Min, draw.Src)
		draw.Draw(i, i.Bounds(), img, img.Bounds().Min, draw.Over)
		img = i
	}

	if enc.toBase64 {
		var buf bytes.Buffer
		err := enc.encode(&buf, img)
		if err != nil {
			return err
		}
		b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
		switch enc.base64Fmt {
		case BASE64:
		case HTML:
			b64 = dataURL(enc.Format, b64, true)
		case URL:
			b64 = dataURL(enc.Format, b64, false)
		}
		_, err = w.Write([]byte(b64))
		if err != nil {
			return err
		}
		return nil
	}

	return enc.encode(w, img)
}

func (enc *Encoder) encode(w io.Writer, img image.Image) error {
	switch enc.Format {
	case JPEG:
		if nrgba, ok := img.(*image.NRGBA); ok && nrgba.Opaque() {
			rgba := &image.RGBA{
				Pix:    nrgba.Pix,
				Stride: nrgba.Stride,
				Rect:   nrgba.Rect,
			}
			return jpeg.Encode(w, rgba, &jpeg.Options{Quality: enc.Quality})
		}
		return jpeg.Encode(w, img, &jpeg.Options{Quality: enc.Quality})

	case PNG:
		encoder := png.Encoder{CompressionLevel: enc.pngCompressionLevel}
		return encoder.Encode(w, img)

	case GIF:
		return gif.Encode(w, img, &gif.Options{
			NumColors: enc.gifNumColors,
			Quantizer: enc.gifQuantizer,
			Drawer:    enc.gifDrawer,
		})

	case TIFF:
		return tiff.Encode(w, img, &tiff.Options{Compression: enc.tiffCompressionType.value(), Predictor: true})

	case BMP:
		return bmp.Encode(w, img)

	case PDF:
		pages := []image.Image{img}
		pages = append(pages, enc.pages...)
		return pdf.Encode(w, pages, &pdf.Options{Quality: enc.Quality})

	case WEBP:
		if len(enc.webpAnimation.Images) > 0 {
			return enc.animatedWebp(w)
		}
		webpOpts := &nativewebp.Options{UseExtendedFormat: enc.webpUseExtendedFormat}
		return nativewebp.Encode(w, img, webpOpts)
	}

	return image.ErrFormat
}

func (enc *Encoder) animatedWebp(w io.Writer) error {
	webpOpts := &nativewebp.Options{UseExtendedFormat: enc.webpUseExtendedFormat}
	return nativewebp.EncodeAll(w, enc.webpAnimation, webpOpts)
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
