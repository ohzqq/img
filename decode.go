package img

import (
	"fmt"
	"image"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/bep/imagemeta"
	"github.com/evanoberholster/imagemeta/xmp"
	"github.com/spf13/cast"
)

type Decoder struct {
	r        io.Reader
	Fmt      Format
	withMeta bool
	opts     imagemeta.Options
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r:    r,
		opts: imagemeta.Options{},
	}
}

func newDecoder(f Format, opts ...DecodeOption) *Decoder {
	dec := &Decoder{
		opts: imagemeta.Options{},
		Fmt:  f,
	}
	for _, opt := range opts {
		opt(dec)
	}
	return dec
}

func (dec *Decoder) Decode(f Format) (image.Image, error) {
	return f.Decode(dec.r)
}

func Open(file string, withMeta bool) (*Img, error) {
	img, err := NewImg(file)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	dec := NewDecoder(f)
	defer f.Close()
	i, err := dec.Decode(img.Fmt)
	if err != nil {
		return nil, err
	}
	img.img = i
	if withMeta {
		err := img.ReadMeta()
		if err != nil {
			return img, err
		}
	}
	return img, nil
}

func (dec *Decoder) DecodeXMP(r io.ReadSeeker) (xmp.XMP, error) {
	dec.withMeta = true
	var tags imagemeta.Tags
	dec.opts.HandleTag = func(ti imagemeta.TagInfo) error {
		if slices.Contains(imgMetaFieldsStr, ti.Tag) {
			tags.Add(ti)
			return nil
		}
		return nil
	}
	dec.opts.R = r
	dec.opts.Sources = imagemeta.EXIF | imagemeta.XMP

	err := imagemeta.Decode(dec.opts)
	if err != nil {
		return xmp.XMP{}, fmt.Errorf("imagemeta decode err %w\n", err)
	}

	x := xmp.XMP{DC: xmp.DublinCore{}}
	for n, ti := range tags.All() {
		switch n {
		case strings.ToLower(Categories.String()):
			hTags, err := UnmarshalHTags([]byte(cast.ToString(ti.Value)))
			if err != nil {
				return x, err
			}
			x.DC.Subject = hTags.StringSlice()
		case strings.ToLower(Caption.String()):
			x.DC.Title = []string{cast.ToString(ti.Value)}
		case Credit.String():
			x.DC.Creator = []string{cast.ToString(ti.Value)}
		case ImageDescription.String():
			x.DC.Description = []string{cast.ToString(ti.Value)}
		}
	}
	return x, nil
}

// open loads an image from file.
// https://github.com/sunshineplan/imgconv
func open(file string) (image.Image, error) {
	img, err := NewImg(file)
	if err != nil {
		return nil, err
	}
	dec := newDecoder(img.Fmt)
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
		img, err := open(file)
		if err != nil {
			return imgs, err
		}
		imgs[i] = img
	}
	return imgs, nil
}
