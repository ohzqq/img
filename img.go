package img

import (
	"encoding/xml"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanoberholster/imagemeta/xmp"
)

type Img struct {
	xmp      xmp.XMP
	Fmt      Format
	img      image.Image
	file     string
	withMeta bool
}

func NewImg(name string) (*Img, error) {
	img := &Img{
		file: name,
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
	f, err := os.Open(img.file)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := NewDecoder(f)
	dec.opts.ImageFormat = img.Fmt.metaFmt()
	x, err := dec.DecodeXMP(f)
	if err != nil {
		return err
	}
	img.xmp = x
	img.xmp.DC.Identifier = img.file
	img.xmp.DC.Format = img.Fmt.ImageType()
	return nil
}

func (img *Img) Open() error {
	f, err := os.Open(img.file)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := NewDecoder(f)
	i, err := dec.Decode(img.Fmt)
	if err != nil {
		return err
	}
	img.img = i
	if img.withMeta {
		err := img.ReadMeta()
		if err != nil {
			return err
		}
	}
	return nil
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

func (dec *Img) DublinCore() xmp.DublinCore {
	return dec.xmp.DC
}

func (dec *Img) EncodeXMP(w io.Writer) error {
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(dec.xmp)
}

var (
	imgMetaFields = []ExifField{
		Categories,
		Caption,
		Credit,
		ImageDescription,
	}
	imgMetaFieldsStr = []string{
		strings.ToLower(Categories.String()),
		strings.ToLower(Caption.String()),
		Credit.String(),
		ImageDescription.String(),
	}
)
