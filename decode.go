package img

import (
	"fmt"
	"image"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/bep/imagemeta"
	"github.com/spf13/cast"
)

type Decoder struct {
	Fmt      Format
	withMeta bool
	opts     imagemeta.Options
}

func NewDecoder(opts ...DecodeOption) *Decoder {
	dec := &Decoder{
		opts: imagemeta.Options{},
	}
	for _, opt := range opts {
		opt(dec)
	}
	return dec
}

func (dec *Decoder) Decode(r io.Reader) (image.Image, error) {
	img, err := dec.Fmt.Decode(r)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func Open(file string, withMeta bool) (*Img, error) {
	img, err := NewImg(file)
	if err != nil {
		return nil, err
	}
	dec := NewDecoder()
	dec.Fmt = img.Fmt
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	i, err := dec.Decode(f)
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

func (dec *Decoder) DecodeMeta(r io.ReadSeeker, i *Img) error {
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
		return fmt.Errorf("imagemeta decode err %w\n", err)
	}

	for n, ti := range tags.All() {
		switch n {
		case strings.ToLower(Categories.String()):
			err := i.hTags.UnmarshalXMP([]byte(cast.ToString(ti.Value)))
			if err != nil {
				return err
			}
			i.xmp.DC.Subject = i.hTags.StringSlice()
		case strings.ToLower(Caption.String()):
			i.xmp.DC.Title = []string{cast.ToString(ti.Value)}
		case Credit.String():
			i.xmp.DC.Creator = []string{cast.ToString(ti.Value)}
		case ImageDescription.String():
			i.xmp.DC.Description = []string{cast.ToString(ti.Value)}
		}
	}
	return nil
}

// open loads an image from file.
// https://github.com/sunshineplan/imgconv
func open(file string) (image.Image, error) {
	dec := NewDecoder()
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
