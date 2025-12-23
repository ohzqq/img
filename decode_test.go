package img

import (
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestDecoder(t *testing.T) {
	c := qt.New(t)
	tstImg := `testdata/test.webp`
	dec := NewDecoder()
	f, err := os.Open(tstImg)
	c.Assert(err, qt.IsNil)
	dec.Fmt = WEBP

	_, err = dec.Decode(f)
	c.Assert(err, qt.IsNil)
}

func TestDecodeAnimatedWebp(t *testing.T) {
	c := qt.New(t)
	tstImg := `testdata/test.webp`
	//tstImg := `testdata/video-001.png`
	f, err := os.Open(tstImg)
	c.Assert(err, qt.IsNil)
	defer f.Close()

	dec := NewDecoder()
	_, err = dec.Fmt.DecodeAnimatedWebP(f)
	c.Assert(err, qt.IsNil)
}
