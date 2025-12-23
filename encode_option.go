package img

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"mime"

	"github.com/HugoSmits86/nativewebp"
	"github.com/spf13/cast"
)

// EncodeOption sets an optional parameter for the Encode and Save functions.
// https://github.com/disintegration/imaging
type EncodeOption func(*Encoder)

var defaultEncodeConfig = &Encoder{
	Quality:             75,
	gifAnimation:        &gif.GIF{},
	gifNumColors:        256,
	gifQuantizer:        nil,
	gifDrawer:           nil,
	pngCompressionLevel: png.DefaultCompression,
	tiffCompressionType: TIFFDeflate,
	padding:             `%02d`,
	webpAnimation:       &nativewebp.Animation{},
	webpDisposal:        1,
	webpDuration:        10,
	//background:          color.Transparent,
}

// Padding returns an EncodeOption that writes a batch of images. Arguments are a
// list of images and a fmt string for zero padding: eg, %04d.
func Padding(padding string) EncodeOption {
	return func(c *Encoder) {
		c.batch = true
		c.padding = padding
	}
}

// Quality returns an EncodeOption that sets the output JPEG or PDF quality.
// Quality ranges from 1 to 100 inclusive, higher is better.
func Quality(quality int) EncodeOption {
	return func(c *Encoder) {
		c.Quality = quality
	}
}

// GIFNumColors returns an EncodeOption that sets the maximum number of colors
// used in the GIF-encoded image. It ranges from 1 to 256.  Default is 256.
func GIFNumColors(numColors int) EncodeOption {
	return func(c *Encoder) {
		c.gifNumColors = numColors
	}
}

// GIFQuantizer returns an EncodeOption that sets the quantizer that is used to produce
// a palette of the GIF-encoded image.
func GIFQuantizer(quantizer draw.Quantizer) EncodeOption {
	return func(c *Encoder) {
		c.gifQuantizer = quantizer
	}
}

// GIFDrawer returns an EncodeOption that sets the drawer that is used to convert
// the source image to the desired palette of the GIF-encoded image.
func GIFDrawer(drawer draw.Drawer) EncodeOption {
	return func(c *Encoder) {
		c.gifDrawer = drawer
	}
}

// GIFDelay returns an EncodeOption that sets the delay for gif frames. This is
// a convenience function to set the same delay for all frames.
func GIFDelay(d int) EncodeOption {
	return func(c *Encoder) {
		c.gifDelay = d
	}
}

// GIFDelays returns an EncodeOption that sets the delay for gif frames.
func GIFDelays(d []int) EncodeOption {
	return func(c *Encoder) {
		c.gifAnimation.Delay = d
	}
}

// GIFDisposal returns an EncodeOption that sets the delay for gif frames.
func GIFDisposal(d []byte) EncodeOption {
	return func(c *Encoder) {
		c.gifAnimation.Disposal = d
	}
}

// GIFBackgroundIndex returns an EncodeOption that sets the delay for gif frames.
func GIFBackgroundIndex(d byte) EncodeOption {
	return func(c *Encoder) {
		c.gifAnimation.BackgroundIndex = d
	}
}

// GIFLoopCount returns an EncodeOption that sets the delay for gif frames.
func GIFLoopCount(d int) EncodeOption {
	return func(c *Encoder) {
		c.gifAnimation.LoopCount = d
	}
}

// PNGCompressionLevel returns an EncodeOption that sets the compression level
// of the PNG-encoded image. Default is png.DefaultCompression.
func PNGCompressionLevel(level png.CompressionLevel) EncodeOption {
	return func(c *Encoder) {
		c.pngCompressionLevel = level
	}
}

// TIFFCompressionType returns an EncodeOption that sets the compression type
// of the TIFF-encoded image. Default is tiff.Deflate.
func TIFFCompressionType(compressionType TIFFCompression) EncodeOption {
	return func(c *Encoder) {
		c.tiffCompressionType = compressionType
	}
}

// WEBPUseExtendedFormat returns EncodeOption that determines whether to use extended format
// of the WEBP-encoded image. Default is false.
func WEBPUseExtendedFormat(b bool) EncodeOption {
	return func(c *Encoder) {
		c.webpUseExtendedFormat = b
	}
}

// WEBPAnimationFrames returns an EncodeOption that sets the webp animation
// frames.
func WEBPAnimationFrames(frames []image.Image) EncodeOption {
	return func(c *Encoder) {
		c.webpAnimation.Images = frames
	}
}

// WEBPAnimationDurations returns an EncodeOption that sets the webp animation
// durations.
func WEBPAnimationDurations(dur []int) EncodeOption {
	return func(c *Encoder) {
		c.webpAnimation.Durations = cast.ToUintSlice(dur)
	}
}

// WEBPAnimationDurations returns an EncodeOption that sets the webp animation
// durations.
func WEBPAnimationDuration(dur int) EncodeOption {
	return func(c *Encoder) {
		c.webpDuration = cast.ToUint(dur)
	}
}

// WEBPAnimationDisposals returns an EncodeOption that sets the webp animation
// durations.
func WEBPAnimationDisposals(disposals []int) EncodeOption {
	return func(c *Encoder) {
		c.webpAnimation.Disposals = cast.ToUintSlice(disposals)
	}
}

// WEBPAnimationDisposal returns an EncodeOption that sets the webp animation
// durations.
func WEBPAnimationDisposal(disposal int) EncodeOption {
	return func(c *Encoder) {
		c.webpDisposal = cast.ToUint(disposal)
	}
}

// WEBPAnimationLoopCount returns an EncodeOption that sets the webp animation
// durations.
func WEBPAnimationLoopCount(loops int) EncodeOption {
	return func(c *Encoder) {
		c.webpAnimation.LoopCount = cast.ToUint16(loops)
	}
}

// WEBPAnimationBackgroundColor returns an EncodeOption that sets the webp animation
// durations.
// Canvas background color in BGRA order, used for clear operations.
func WEBPAnimationBackgroundColor(color uint32) EncodeOption {
	return func(c *Encoder) {
		c.webpAnimation.BackgroundColor = color
	}
}

// BackgroundColor returns an EncodeOption that sets the background color.
func BackgroundColor(color color.Color) EncodeOption {
	return func(c *Encoder) {
		c.background = color
	}
}

// PDFPages returns an EncodeOption that sets multiple pages for pdf conversion.
func PDFPages(pages []image.Image) EncodeOption {
	return func(c *Encoder) {
		c.pages = pages
	}
}

// Base64 returns an EncodeOption that encodes the format to Base64.
func Base64(outFmt Format) EncodeOption {
	return func(c *Encoder) {
		c.toBase64 = true
		c.base64Fmt = outFmt
	}
}

func init() {
	for _, ext := range []string{".tif", ".tiff"} {
		mime.AddExtensionType(ext, `image/tiff`)
	}
	for _, ext := range []string{".b64", "uue"} {
		mime.AddExtensionType(ext, "text/plain")
	}
}
