package img

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/evanoberholster/imagemeta/xmp"
)

type Img struct {
	xmp   xmp.XMP
	Fmt   Format
	toFmt Format
	img   image.Image
	hTags *Tags
	meta  map[string]string
	file  string
}

func NewImg(name string) (*Img, error) {
	img := &Img{
		file:  name,
		meta:  make(map[string]string),
		toFmt: Format(-1),
		hTags: NewTags(),
	}
	ext := filepath.Ext(name)
	imgFmt, err := FormatFromExtension(ext)
	if err != nil {
		return nil, fmt.Errorf("format error %w", err)
	}
	img.Fmt = imgFmt
	return img, nil
}

func (img *Img) ReadMeta() error {
	dec := NewDecoder()
	img.xmp = xmp.XMP{
		DC: xmp.DublinCore{
			Identifier: img.file,
			Format:     img.Fmt.ImageType(),
		},
	}
	dec.opts.ImageFormat = img.Fmt.metaFmt()
	f, err := os.Open(img.file)
	if err != nil {
		return err
	}
	defer f.Close()
	return dec.DecodeMeta(f, img)
}

func (img *Img) Save(opts ...EncodeOption) error {
	if img.img == nil {
		i, err := open(img.file)
		if err != nil {
			return err
		}
		img.img = i
	}
	enc := NewEncoder(img.Fmt, opts...)
	return enc.Save(img.file, img.img)
}

func (img *Img) SaveAs(name string, opts ...EncodeOption) error {
	to, err := FormatFromExtension(filepath.Ext(name))
	if err != nil {
		return fmt.Errorf("can't save as format %w", err)
	}
	if img.img == nil {
		i, err := open(img.file)
		if err != nil {
			return err
		}
		img.img = i
	}
	enc := NewEncoder(to, opts...)
	return enc.Save(name, img.img)
}
