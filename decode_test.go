package img

import (
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestDecoder(t *testing.T) {
	c := qt.New(t)
	tstImg := `testdata/test.webp`
	f, err := os.Open(tstImg)
	c.Assert(err, qt.IsNil)
	dec := NewDecoder(f)

	_, err = dec.Decode(WEBP)
	c.Assert(err, qt.IsNil)
}

func TestDecodeAnimatedWebp(t *testing.T) {
	c := qt.New(t)
	tstImg := `testdata/test.webp`
	//tstImg := `testdata/video-001.png`
	f, err := os.Open(tstImg)
	c.Assert(err, qt.IsNil)
	defer f.Close()

	dec := newDecoder(WEBP)
	_, err = dec.Fmt.DecodeAnimatedWebP(f)
	c.Assert(err, qt.IsNil)
}
