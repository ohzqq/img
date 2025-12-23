package img

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestAnimation(t *testing.T) {
	tstImgs := []string{
		`testdata/body0001.tif`,
		`testdata/body0002.tif`,
		`testdata/body0003.tif`,
		`testdata/body0004.tif`,
	}
	enc := NewEncoder(WEBP,
		WEBPAnimationDuration(1),
	)
	err := enc.AnimatedWEBP(`testdata/animated.webp`, tstImgs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestConvert(t *testing.T) {
	c := qt.New(t)
	tstImg := `testdata/video-001.png`
	img, err := New(tstImg)
	c.Assert(err, qt.IsNil)

	outImg := `testdata/convert.jpg`
	err = img.SaveAs(outImg)
	c.Assert(err, qt.IsNil)
}

func TestBatch(t *testing.T) {
	tstBatch := []string{
		"testdata/video-001.bmp",
		"testdata/video-001.gif",
		"testdata/video-001.jpg",
		"testdata/video-001.tif",
		"testdata/video-001.webp",
	}
	tstImgs, err := OpenAll(tstBatch)
	if err != nil {
		t.Fatal(err)
	}

	err = SaveAll(`testdata/batch.webp`, tstImgs,
		Padding("%07d"),
	)
	if err != nil {
		t.Fatal(err)
	}
}
