package img

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/evanoberholster/imagemeta/imagetype"
	qt "github.com/frankban/quicktest"
)

func TestImgMeta(t *testing.T) {
	c := qt.New(t)
	f, err := os.Open(`testdata/test.webp`)
	c.Assert(err, qt.IsNil)
	defer f.Close()

	i, err := NewImg(`testdata/test.webp`)
	c.Assert(err, qt.IsNil)

	err = i.DecodeMeta(f)
	c.Assert(err, qt.IsNil)
}

func TestEncodeImgMeta(t *testing.T) {
	c := qt.New(t)
	f, err := os.Open(`testdata/test.webp`)
	c.Assert(err, qt.IsNil)
	defer f.Close()

	i, err := NewImg(`testdata/test.webp`)
	c.Assert(err, qt.IsNil)

	err = i.DecodeMeta(f)
	c.Assert(err, qt.IsNil)

	err = i.EncodeXMP(os.Stdout)
	c.Assert(err, qt.IsNil)
}

func TestImageFmtConvert(t *testing.T) {
	c := qt.New(t)
	i, err := NewImg(`testdata/test.webp`)
	c.Assert(err, qt.IsNil)

	it := imagetype.FromString(filepath.Ext(`testdata/test.webp`))
	c.Assert(it, qt.Equals, i.xmp.DC.Format)
}
