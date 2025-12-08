package img

import (
	"image"
	"io"
	"os"
)

// Decode reads an image from r.
// If want to use custom image format packages which were registered in image package, please
// make sure these custom packages imported before importing imgconv package.
// https://github.com/disintegration/imaging
func Decode(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return img, nil
}

//
// resize.go
//

// DecodeConfig decodes the color model and dimensions of an image that has been encoded in a
// registered format. The string returned is the format name used during format registration.
// https://github.com/sunshineplan/imgconv
func DecodeConfig(r io.Reader) (image.Config, string, error) {
	return image.DecodeConfig(r)
}

// Open loads an image from file.
// https://github.com/sunshineplan/imgconv
func Open(file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Decode(f)
}

// Write image according format option
// https://github.com/sunshineplan/imgconv
func Write(w io.Writer, base image.Image, option *Encoder) error {
	return option.Encode(w, base)
}

// Save saves image according format option
// https://github.com/sunshineplan/imgconv
func Save(output string, base image.Image, option *Encoder) error {
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	return option.Encode(f, base)
}
