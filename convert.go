package img

import (
	"image"
	"path/filepath"
)

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
