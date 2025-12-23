package img

import (
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
		return nil, err
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

func (dec *Decoder) DecodeImage() error {
	f, err := os.Open(dec.file)
	if err != nil {
		return err
	}
	defer f.Close()
	dec.img, err = Decode(f)
	if err != nil {
		return err
	}
	return nil
}

// Decode reads an image from r.
// If want to use custom image format packages which were registered in image package, please
// make sure these custom packages imported before importing imgconv package.
// https://github.com/disintegration/imaging
func Decode(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// DecodeConfig decodes the color model and dimensions of an image that has been encoded in a
// registered format. The string returned is the format name used during format registration.
// https://github.com/sunshineplan/imgconv
func DecodeConfig(r io.Reader) (image.Config, string, error) {
	return image.DecodeConfig(r)
}

// Open loads an image from file.
// https://github.com/sunshineplan/imgconv
func Open(file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Decode(f)
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
