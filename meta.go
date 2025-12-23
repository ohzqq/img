package imgtag

import (
	"encoding/xml"
	"io"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bep/imagemeta"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/xmp"
	"github.com/ohzqq/img"
	"github.com/spf13/cast"
)

type Img struct {
	hTags *Tags
	Fmt   img.Format
	meta  map[string]string
	opts  imagemeta.Options
	xmp   xmp.XMP
}

func NewImg(name string) (*Img, error) {
	i := &Img{
		meta: make(map[string]string),
		opts: imagemeta.Options{},
	}
	ext := filepath.Ext(name)
	f := imagetype.FromString(ext)
	i.xmp = xmp.XMP{
		DC: xmp.DublinCore{
			Identifier: name,
			Format:     f,
		},
	}
	for _, ifmt := range metaFmts {
		it := imagetype.FromString("." + ifmt.String())
		if it == f {
			i.opts.ImageFormat = ifmt
		}
	}
	return i, nil
}

func (meta *Img) DublinCore() xmp.DublinCore {
	return meta.xmp.DC
}

func (meta *Img) SetField(f string, val string) *Img {
	meta.meta[f] = val
	return meta
}

func (meta *Img) GetField(f string) string {
	if v, ok := meta.meta[f]; ok {
		return v
	}
	return ""
}

func (i *Img) DecodeMeta(r io.ReadSeeker) error {
	var tags imagemeta.Tags
	i.opts.HandleTag = func(ti imagemeta.TagInfo) error {
		if slices.Contains(imgMetaFieldsStr, ti.Tag) {
			tags.Add(ti)
			return nil
		}
		return nil
	}
	i.opts.R = r
	i.opts.Sources = imagemeta.EXIF | imagemeta.XMP

	err := imagemeta.Decode(i.opts)
	if err != nil {
		return err
	}

	for n, ti := range tags.All() {
		switch n {
		case strings.ToLower(Categories.String()):
			t, err := UnmarshalHTags([]byte(cast.ToString(ti.Value)))
			if err != nil {
				return err
			}
			i.hTags = t
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

func (i *Img) EncodeXMP(w io.Writer) error {
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(i.xmp)
}

var metaFmts = []imagemeta.ImageFormat{
	imagemeta.JPEG,
	imagemeta.TIFF,
	imagemeta.PNG,
	imagemeta.WebP,
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
