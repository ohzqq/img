//go:build ignore

package img

import (
	"os"
	"testing"

	"github.com/bep/imagemeta"
	"github.com/evanoberholster/imagemeta/imagetype"
	qt "github.com/frankban/quicktest"
)

func TestImgMeta(t *testing.T) {
	c := qt.New(t)
	f, err := os.Open(`testdata/test.webp`)
	c.Assert(err, qt.IsNil)
	defer f.Close()

	i, err := NewDecoder(`testdata/test.webp`, WithMeta())
	c.Assert(err, qt.IsNil)

	err = i.DecodeMeta(f)
	c.Assert(err, qt.IsNil)
}

func TestDecodeImgMeta(t *testing.T) {
	c := qt.New(t)
	f, err := os.Open(`testdata/test.webp`)
	c.Assert(err, qt.IsNil)
	defer f.Close()

	i, err := NewDecoder(`testdata/test.webp`, WithMeta())
	c.Assert(err, qt.IsNil)

	err = i.DecodeMeta(f)
	c.Assert(err, qt.IsNil)

	//err = i.EncodeXMP(os.Stdout)
	//c.Assert(err, qt.IsNil)
}

func TestImageFmtConvert(t *testing.T) {
	c := qt.New(t)
	i, err := NewDecoder(`testdata/test.webp`, WithMeta())
	c.Assert(err, qt.IsNil)

	it := imagetype.ImageWebP
	mt := it.String()
	metaT := imagemeta.WebP
	c.Assert(it, qt.Equals, i.xmp.DC.Format)
	c.Assert(mt, qt.Equals, "image/webp")
	c.Assert(metaT, qt.Equals, i.opts.ImageFormat)
}
