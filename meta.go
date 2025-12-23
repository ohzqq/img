package img

import (
	"encoding/xml"
	"io"
	"strings"

	"github.com/evanoberholster/imagemeta/xmp"
)

func (dec *Img) DublinCore() xmp.DublinCore {
	return dec.xmp.DC
}

func (dec *Img) SetField(f string, val string) *Img {
	dec.meta[f] = val
	return dec
}

func (dec *Img) GetField(f string) string {
	if v, ok := dec.meta[f]; ok {
		return v
	}
	return ""
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
