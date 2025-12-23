package img

import (
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestDecoder(t *testing.T) {
	c := qt.New(t)
	tstImg := `testdata/test.webp`
	dec, err := NewDecoder(tstImg)
	c.Assert(err, qt.IsNil)

	err = dec.Open()
	c.Assert(err, qt.IsNil)
}

func TestDecodeWithMeta(t *testing.T) {
	c := qt.New(t)
	tstImg := `testdata/test.webp`
	dec, err := NewDecoder(tstImg, WithMeta())
	c.Assert(err, qt.IsNil)

	err = dec.Open()
	c.Assert(err, qt.IsNil)

	dc := dec.DublinCore()
	c.Assert(dc.Description[0], qt.Equals, "insert")

	//err = dec.EncodeXMP(os.Stdout)
	//c.Assert(err, qt.IsNil)
}

func TestDecodeAnimatedWebp(t *testing.T) {
	c := qt.New(t)
	tstImg := `testdata/test.webp`
	//tstImg := `testdata/video-001.png`
	f, err := os.Open(tstImg)
	c.Assert(err, qt.IsNil)
	defer f.Close()

	dec, err := NewDecoder(tstImg)
	c.Assert(err, qt.IsNil)
	_, err = dec.Fmt.DecodeAnimatedWebP(f)
	c.Assert(err, qt.IsNil)
}
