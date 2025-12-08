package img

import (
	"image"
	"io"
	"os"
	"path/filepath"
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

// OpenAll loads images from files.
func OpenAll(files []string) ([]image.Image, error) {
	imgs := make([]image.Image, len(files))
	for i, file := range files {
		img, err := Open(file)
		if err != nil {
			return imgs, err
		}
		imgs[i] = img
	}
	return imgs, nil
}

// Save saves image according to the encoder
// https://github.com/sunshineplan/imgconv
func Save(output string, base image.Image, opts ...EncodeOption) error {
	ext := filepath.Ext(output)
	imgFmt, err := FormatFromExtension(ext)
	if err != nil {
		return err
	}
	return NewEncoder(imgFmt, opts...).Save(output, base)
}

// SaveAll saves images according to the encoder
// https://github.com/sunshineplan/imgconv
func SaveAll(output string, images []image.Image, opts ...EncodeOption) error {
	ext := filepath.Ext(output)
	imgFmt, err := FormatFromExtension(ext)
	if err != nil {
		return err
	}
	return NewEncoder(imgFmt, opts...).SaveAll(output, images)
}
