package img

import (
	"image"
	"io"
	"path/filepath"
)

var defaultFormat = &FormatOption{Format: JPEG}

// Options represents options that can be used to configure a image operation.
type Options struct {
	Format *FormatOption
}

// NewOptions creates a new option with default setting.
func NewOptions() *Options {
	return &Options{Format: defaultFormat}
}

// SetFormat sets the value for the Format field.
func (opts *Options) SetFormat(f Format, options ...EncodeOption) *Options {
	opts.Format = &FormatOption{f, options}
	return opts
}

// Convert image according options opts.
func (opts *Options) Convert(w io.Writer, base image.Image) error {

	if opts.Format == nil {
		opts.Format = defaultFormat
	}

	return opts.Format.Encode(w, base)
}

// ConvertExt convert filename's ext according image format.
func (opts *Options) ConvertExt(filename string) string {
	return filename[0:len(filename)-len(filepath.Ext(filename))] + formatExts[opts.Format.Format][0]
}
