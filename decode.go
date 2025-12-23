package img

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"

	"github.com/bep/imagemeta"
	"github.com/evanoberholster/imagemeta/xmp"
)

type Decoder struct {
	file     string
	withMeta bool
	Fmt      Format
	img      image.Image
	hTags    *Tags
	meta     map[string]string
	opts     imagemeta.Options
	xmp      xmp.XMP
}

func NewDecoder(name string, opts ...DecodeOption) (*Decoder, error) {
	dec := &Decoder{
		file:  name,
		meta:  make(map[string]string),
		opts:  imagemeta.Options{},
		hTags: NewTags(),
	}
	ext := filepath.Ext(name)
	imgFmt, err := FormatFromExtension(ext)
	if err != nil {
		return nil, fmt.Errorf("error initializing new decoder %w", err)
	}
	dec.Fmt = imgFmt
	for _, opt := range opts {
		opt(dec)
	}
	if dec.withMeta {
		dec.xmp = xmp.XMP{
			DC: xmp.DublinCore{
				Identifier: name,
				Format:     dec.Fmt.ImageType(),
			},
		}
		dec.opts.ImageFormat = dec.Fmt.metaFmt()
	}
	return dec, nil
}

func (dec *Decoder) Open() error {
	f, err := os.Open(dec.file)
	if err != nil {
		return fmt.Errorf("decoder.DecodeImage %w\n", err)
	}
	defer f.Close()
	return dec.Decode(f)
}

func (dec *Decoder) Decode(r io.Reader) error {
	if dec.withMeta {
		return dec.DecodeWithMeta(r)
	}
	img, err := dec.Fmt.Decode(r)
	if err != nil {
		return err
	}
	dec.img = img
	return nil
}

func (dec *Decoder) DecodeWithMeta(r io.Reader) error {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		return err
	}
	rs := bytes.NewReader(buf.Bytes())
	img, err := dec.Fmt.Decode(&buf)
	if err != nil {
		return err
	}
	dec.img = img
	err = dec.DecodeMeta(rs)
	if err != nil {
		return err
	}
	return nil
}

// Open loads an image from file.
// https://github.com/sunshineplan/imgconv
func Open(file string) (image.Image, error) {
	dec, err := NewDecoder(file)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return dec.Fmt.Decode(f)
}

// OpenAll loads images from files.
func OpenAll(files []string) ([]image.Image, error) {
	imgs := make([]image.Image, len(files))
	for i, file := range files {
		img, err := Open(file)
		if err != nil {
			return imgs, err
		}
		imgs[i] = img
	}
	return imgs, nil
}
