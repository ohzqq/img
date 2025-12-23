package img

import (
	"encoding/xml"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/bep/imagemeta"
	"github.com/evanoberholster/imagemeta/xmp"
	"github.com/spf13/cast"
)

func (i *Decoder) DecodeMeta(r io.ReadSeeker) error {
	i.withMeta = true
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

func (meta *Decoder) DublinCore() xmp.DublinCore {
	return meta.xmp.DC
}

func (meta *Decoder) SetField(f string, val string) *Decoder {
	meta.meta[f] = val
	return meta
}

func (meta *Decoder) GetField(f string) string {
	if v, ok := meta.meta[f]; ok {
		return v
	}
	return ""
}

func (i *Decoder) EncodeXMP(w io.Writer) error {
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(i.xmp)
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
