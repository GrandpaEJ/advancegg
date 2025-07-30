package core

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"

	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func Degrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	im, _, err := image.Decode(file)
	return im, err
}

func LoadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(file)
}

func SavePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, im)
}

// SaveJPEG encodes the image as a JPEG and writes it to disk.
func SaveJPEG(path string, im image.Image, quality int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return jpeg.Encode(file, im, &jpeg.Options{Quality: quality})
}

// SaveGIF encodes the image as a GIF and writes it to disk.
func SaveGIF(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return gif.Encode(file, im, nil)
}

// SaveBMP encodes the image as a BMP and writes it to disk.
func SaveBMP(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return bmp.Encode(file, im)
}

// SaveTIFF encodes the image as a TIFF and writes it to disk.
func SaveTIFF(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return tiff.Encode(file, im, nil)
}

func LoadJPG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return jpeg.Decode(file)
}

// LoadGIF loads a GIF image from the specified file path.
func LoadGIF(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return gif.Decode(file)
}

// LoadBMP loads a BMP image from the specified file path.
func LoadBMP(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return bmp.Decode(file)
}

// LoadTIFF loads a TIFF image from the specified file path.
func LoadTIFF(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return tiff.Decode(file)
}

// LoadWebP loads a WebP image from the specified file path.
func LoadWebP(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return webp.Decode(file)
}

func SaveJPG(path string, im image.Image, quality int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var opt jpeg.Options
	opt.Quality = quality

	return jpeg.Encode(file, im, &opt)
}

func imageToRGBA(src image.Image) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

func parseHexColor(x string) (r, g, b, a int) {
	x = strings.TrimPrefix(x, "#")
	a = 255
	if len(x) == 3 {
		format := "%1x%1x%1x"
		fmt.Sscanf(x, format, &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
	}
	if len(x) == 6 {
		format := "%02x%02x%02x"
		fmt.Sscanf(x, format, &r, &g, &b)
	}
	if len(x) == 8 {
		format := "%02x%02x%02x%02x"
		fmt.Sscanf(x, format, &r, &g, &b, &a)
	}
	return
}

func fixp(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{fix(x), fix(y)}
}

func fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}

func unfix(x fixed.Int26_6) float64 {
	const shift, mask = 6, 1<<6 - 1
	if x >= 0 {
		return float64(x>>shift) + float64(x&mask)/64
	}
	x = -x
	if x >= 0 {
		return -(float64(x>>shift) + float64(x&mask)/64)
	}
	return 0
}

// LoadFontFace is a helper function to load the specified font file with
// the specified point size. Supports both TTF and OTF font formats.
// Note that the returned `font.Face` objects are not thread safe and
// cannot be used in parallel across goroutines.
// You can usually just use the Context.LoadFontFace function instead of
// this package-level function.
func LoadFontFace(path string, points float64) (font.Face, error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseFontFace(fontBytes, points)
}

// LoadTTFFace loads a TTF font file with the specified point size.
func LoadTTFFace(path string, points float64) (font.Face, error) {
	return LoadFontFace(path, points) // TTF and OTF use the same parser
}

// LoadOTFFace loads an OTF font file with the specified point size.
func LoadOTFFace(path string, points float64) (font.Face, error) {
	return LoadFontFace(path, points) // TTF and OTF use the same parser
}

// ParseFontFace parses font data from bytes and creates a font face.
// Supports both TTF and OTF formats.
func ParseFontFace(fontBytes []byte, points float64) (font.Face, error) {
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return face, nil
}

// ParseFontFaceWithOptions parses font data with custom options.
func ParseFontFaceWithOptions(fontBytes []byte, options *truetype.Options) (font.Face, error) {
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, options)
	return face, nil
}

// GetFontFormat attempts to detect the font format from the file header.
func GetFontFormat(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the first 4 bytes to check the signature
	header := make([]byte, 4)
	_, err = file.Read(header)
	if err != nil {
		return "", err
	}

	// Check font format signatures
	switch {
	case header[0] == 0x00 && header[1] == 0x01 && header[2] == 0x00 && header[3] == 0x00:
		return "TTF", nil
	case string(header) == "OTTO":
		return "OTF", nil
	case string(header) == "true" || string(header) == "typ1":
		return "TTF", nil // Some TTF variants
	default:
		return "UNKNOWN", nil
	}
}
