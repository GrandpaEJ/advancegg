package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/color"
	"math"
)

// ICC Color Profile support for accurate color conversion

// ICCProfile represents an ICC color profile
type ICCProfile struct {
	Header        ICCHeader
	TagTable      []ICCTag
	Data          []byte
	ColorSpace    ColorSpace
	PCS           ColorSpace // Profile Connection Space
	Intent        RenderingIntent
	WhitePoint    XYZColor
	BlackPoint    XYZColor
	RedColorant   XYZColor
	GreenColorant XYZColor
	BlueColorant  XYZColor
	Curves        []ToneCurve
}

// ICCHeader represents the ICC profile header
type ICCHeader struct {
	ProfileSize        uint32
	PreferredCMM       [4]byte
	ProfileVersion     uint32
	DeviceClass        DeviceClass
	DataColorSpace     ColorSpace
	PCS                ColorSpace
	CreationDateTime   [12]byte
	PlatformSignature  [4]byte
	ProfileFlags       uint32
	DeviceManufacturer [4]byte
	DeviceModel        [4]byte
	DeviceAttributes   uint64
	RenderingIntent    RenderingIntent
	PCSIlluminant      XYZColor
	ProfileCreator     [4]byte
	Reserved           [44]byte
}

// ICCTag represents an ICC tag
type ICCTag struct {
	Signature [4]byte
	Offset    uint32
	Size      uint32
}

// DeviceClass represents ICC device class
type DeviceClass uint32

const (
	DeviceClassInput      DeviceClass = 0x73636E72 // 'scnr'
	DeviceClassDisplay    DeviceClass = 0x6D6E7472 // 'mntr'
	DeviceClassOutput     DeviceClass = 0x70727472 // 'prtr'
	DeviceClassLink       DeviceClass = 0x6C696E6B // 'link'
	DeviceClassAbstract   DeviceClass = 0x61627374 // 'abst'
	DeviceClassNamedColor DeviceClass = 0x6E6D636C // 'nmcl'
)

// RenderingIntent represents ICC rendering intent
type RenderingIntent uint32

const (
	IntentPerceptual           RenderingIntent = 0
	IntentRelativeColorimetric RenderingIntent = 1
	IntentSaturation           RenderingIntent = 2
	IntentAbsoluteColorimetric RenderingIntent = 3
)

// XYZColor represents a color in CIE XYZ color space
type XYZColor struct {
	X, Y, Z float64
}

// ToneCurve represents a tone reproduction curve
type ToneCurve struct {
	Type   uint32
	Points []float64
}

// ColorSpace represents different color spaces
type ColorSpace uint32

const (
	ColorSpaceXYZ  ColorSpace = 0x58595A20 // 'XYZ '
	ColorSpaceLab  ColorSpace = 0x4C616220 // 'Lab '
	ColorSpaceRGB  ColorSpace = 0x52474220 // 'RGB '
	ColorSpaceCMYK ColorSpace = 0x434D594B // 'CMYK'
	ColorSpaceGray ColorSpace = 0x47524159 // 'GRAY'
)

// Common ICC tag signatures
const (
	TagRedColorant   = 0x7258595A // 'rXYZ'
	TagGreenColorant = 0x6758595A // 'gXYZ'
	TagBlueColorant  = 0x6258595A // 'bXYZ'
	TagWhitePoint    = 0x77747074 // 'wtpt'
	TagBlackPoint    = 0x626B7074 // 'bkpt'
	TagRedTRC        = 0x72545243 // 'rTRC'
	TagGreenTRC      = 0x67545243 // 'gTRC'
	TagBlueTRC       = 0x62545243 // 'bTRC'
	TagGrayTRC       = 0x6B545243 // 'kTRC'
	TagDescription   = 0x64657363 // 'desc'
	TagCopyright     = 0x63707274 // 'cprt'
)

// Standard illuminants
var (
	IlluminantD50 = XYZColor{X: 0.9642, Y: 1.0000, Z: 0.8249}
	IlluminantD65 = XYZColor{X: 0.9505, Y: 1.0000, Z: 1.0890}
)

// NewICCProfile creates a new ICC profile
func NewICCProfile() *ICCProfile {
	return &ICCProfile{
		ColorSpace: ColorSpaceRGB,
		PCS:        ColorSpaceXYZ,
		Intent:     IntentPerceptual,
		WhitePoint: IlluminantD65,
		Curves:     make([]ToneCurve, 0),
	}
}

// LoadICCProfile loads an ICC profile from data
func LoadICCProfile(data []byte) (*ICCProfile, error) {
	if len(data) < 128 {
		return nil, fmt.Errorf("ICC profile too small")
	}

	profile := NewICCProfile()
	reader := bytes.NewReader(data)

	// Read header
	if err := binary.Read(reader, binary.BigEndian, &profile.Header); err != nil {
		return nil, fmt.Errorf("failed to read ICC header: %v", err)
	}

	// Validate profile
	if profile.Header.ProfileSize != uint32(len(data)) {
		return nil, fmt.Errorf("profile size mismatch")
	}

	// Read tag table
	var tagCount uint32
	if err := binary.Read(reader, binary.BigEndian, &tagCount); err != nil {
		return nil, fmt.Errorf("failed to read tag count: %v", err)
	}

	profile.TagTable = make([]ICCTag, tagCount)
	for i := uint32(0); i < tagCount; i++ {
		if err := binary.Read(reader, binary.BigEndian, &profile.TagTable[i]); err != nil {
			return nil, fmt.Errorf("failed to read tag %d: %v", i, err)
		}
	}

	// Parse important tags
	profile.parseColorants(data)
	profile.parseToneCurves(data)

	profile.Data = data
	return profile, nil
}

// parseColorants parses colorant tags
func (p *ICCProfile) parseColorants(data []byte) {
	for _, tag := range p.TagTable {
		signature := binary.BigEndian.Uint32(tag.Signature[:])

		switch signature {
		case TagRedColorant:
			p.RedColorant = p.parseXYZTag(data, tag)
		case TagGreenColorant:
			p.GreenColorant = p.parseXYZTag(data, tag)
		case TagBlueColorant:
			p.BlueColorant = p.parseXYZTag(data, tag)
		case TagWhitePoint:
			p.WhitePoint = p.parseXYZTag(data, tag)
		case TagBlackPoint:
			p.BlackPoint = p.parseXYZTag(data, tag)
		}
	}
}

// parseXYZTag parses an XYZ tag
func (p *ICCProfile) parseXYZTag(data []byte, tag ICCTag) XYZColor {
	if tag.Offset+tag.Size > uint32(len(data)) {
		return XYZColor{}
	}

	tagData := data[tag.Offset : tag.Offset+tag.Size]
	if len(tagData) < 20 {
		return XYZColor{}
	}

	// Skip type signature and reserved bytes
	reader := bytes.NewReader(tagData[8:])

	var x, y, z uint32
	binary.Read(reader, binary.BigEndian, &x)
	binary.Read(reader, binary.BigEndian, &y)
	binary.Read(reader, binary.BigEndian, &z)

	// Convert from fixed point to float
	return XYZColor{
		X: float64(x) / 65536.0,
		Y: float64(y) / 65536.0,
		Z: float64(z) / 65536.0,
	}
}

// parseToneCurves parses tone reproduction curves
func (p *ICCProfile) parseToneCurves(data []byte) {
	curves := make([]ToneCurve, 0, 3)

	for _, tag := range p.TagTable {
		signature := binary.BigEndian.Uint32(tag.Signature[:])

		switch signature {
		case TagRedTRC, TagGreenTRC, TagBlueTRC, TagGrayTRC:
			curve := p.parseCurveTag(data, tag)
			curves = append(curves, curve)
		}
	}

	p.Curves = curves
}

// parseCurveTag parses a curve tag
func (p *ICCProfile) parseCurveTag(data []byte, tag ICCTag) ToneCurve {
	if tag.Offset+tag.Size > uint32(len(data)) {
		return ToneCurve{}
	}

	tagData := data[tag.Offset : tag.Offset+tag.Size]
	if len(tagData) < 12 {
		return ToneCurve{}
	}

	reader := bytes.NewReader(tagData)

	var typeSignature uint32
	var reserved uint32
	var count uint32

	binary.Read(reader, binary.BigEndian, &typeSignature)
	binary.Read(reader, binary.BigEndian, &reserved)
	binary.Read(reader, binary.BigEndian, &count)

	curve := ToneCurve{
		Type:   typeSignature,
		Points: make([]float64, count),
	}

	if count == 0 {
		// Linear curve
		curve.Points = []float64{0.0, 1.0}
	} else if count == 1 {
		// Gamma curve
		var gamma uint16
		binary.Read(reader, binary.BigEndian, &gamma)
		curve.Points = []float64{float64(gamma) / 256.0}
	} else {
		// Table curve
		for i := uint32(0); i < count; i++ {
			var value uint16
			binary.Read(reader, binary.BigEndian, &value)
			curve.Points[i] = float64(value) / 65535.0
		}
	}

	return curve
}

// ColorConverter handles color space conversions using ICC profiles
type ColorConverter struct {
	SourceProfile *ICCProfile
	DestProfile   *ICCProfile
	Intent        RenderingIntent
}

// NewColorConverter creates a new color converter
func NewColorConverter(source, dest *ICCProfile) *ColorConverter {
	return &ColorConverter{
		SourceProfile: source,
		DestProfile:   dest,
		Intent:        IntentPerceptual,
	}
}

// ConvertColor converts a color from source to destination profile
func (cc *ColorConverter) ConvertColor(c color.Color) color.Color {
	// Convert to RGBA first
	r, g, b, a := c.RGBA()
	rf := float64(r) / 65535.0
	gf := float64(g) / 65535.0
	bf := float64(b) / 65535.0
	af := float64(a) / 65535.0

	// Convert source RGB to XYZ
	xyz := cc.rgbToXYZ(rf, gf, bf, cc.SourceProfile)

	// Apply chromatic adaptation if needed
	if cc.SourceProfile.WhitePoint != cc.DestProfile.WhitePoint {
		xyz = cc.chromaticAdaptation(xyz, cc.SourceProfile.WhitePoint, cc.DestProfile.WhitePoint)
	}

	// Convert XYZ to destination RGB
	rf, gf, bf = cc.xyzToRGB(xyz, cc.DestProfile)

	// Clamp values
	rf = math.Max(0, math.Min(1, rf))
	gf = math.Max(0, math.Min(1, gf))
	bf = math.Max(0, math.Min(1, bf))

	return color.RGBA{
		R: uint8(rf * 255),
		G: uint8(gf * 255),
		B: uint8(bf * 255),
		A: uint8(af * 255),
	}
}

// rgbToXYZ converts RGB to XYZ using the profile's colorants and curves
func (cc *ColorConverter) rgbToXYZ(r, g, b float64, profile *ICCProfile) XYZColor {
	// Apply inverse tone curves
	if len(profile.Curves) >= 3 {
		r = cc.applyInverseCurve(r, profile.Curves[0])
		g = cc.applyInverseCurve(g, profile.Curves[1])
		b = cc.applyInverseCurve(b, profile.Curves[2])
	}

	// Matrix transformation using colorants
	x := r*profile.RedColorant.X + g*profile.GreenColorant.X + b*profile.BlueColorant.X
	y := r*profile.RedColorant.Y + g*profile.GreenColorant.Y + b*profile.BlueColorant.Y
	z := r*profile.RedColorant.Z + g*profile.GreenColorant.Z + b*profile.BlueColorant.Z

	return XYZColor{X: x, Y: y, Z: z}
}

// xyzToRGB converts XYZ to RGB using the profile's colorants and curves
func (cc *ColorConverter) xyzToRGB(xyz XYZColor, profile *ICCProfile) (float64, float64, float64) {
	// Inverse matrix transformation
	// This is a simplified implementation - real ICC would use proper matrix inversion
	r := xyz.X/profile.RedColorant.X + xyz.Y/profile.RedColorant.Y + xyz.Z/profile.RedColorant.Z
	g := xyz.X/profile.GreenColorant.X + xyz.Y/profile.GreenColorant.Y + xyz.Z/profile.GreenColorant.Z
	b := xyz.X/profile.BlueColorant.X + xyz.Y/profile.BlueColorant.Y + xyz.Z/profile.BlueColorant.Z

	// Apply tone curves
	if len(profile.Curves) >= 3 {
		r = cc.applyCurve(r, profile.Curves[0])
		g = cc.applyCurve(g, profile.Curves[1])
		b = cc.applyCurve(b, profile.Curves[2])
	}

	return r, g, b
}

// applyCurve applies a tone curve to a value
func (cc *ColorConverter) applyCurve(value float64, curve ToneCurve) float64 {
	if len(curve.Points) == 0 {
		return value // Linear
	} else if len(curve.Points) == 1 {
		// Gamma curve
		gamma := curve.Points[0]
		return math.Pow(value, gamma)
	} else {
		// Table interpolation
		if value <= 0 {
			return curve.Points[0]
		}
		if value >= 1 {
			return curve.Points[len(curve.Points)-1]
		}

		// Linear interpolation
		index := value * float64(len(curve.Points)-1)
		i := int(index)
		if i >= len(curve.Points)-1 {
			return curve.Points[len(curve.Points)-1]
		}

		t := index - float64(i)
		return curve.Points[i]*(1-t) + curve.Points[i+1]*t
	}
}

// applyInverseCurve applies the inverse of a tone curve
func (cc *ColorConverter) applyInverseCurve(value float64, curve ToneCurve) float64 {
	if len(curve.Points) == 0 {
		return value // Linear
	} else if len(curve.Points) == 1 {
		// Inverse gamma
		gamma := curve.Points[0]
		if gamma != 0 {
			return math.Pow(value, 1.0/gamma)
		}
		return value
	} else {
		// Inverse table lookup (simplified)
		for i, point := range curve.Points {
			if point >= value {
				if i == 0 {
					return 0
				}
				// Linear interpolation
				t := (value - curve.Points[i-1]) / (point - curve.Points[i-1])
				return (float64(i-1) + t) / float64(len(curve.Points)-1)
			}
		}
		return 1.0
	}
}

// chromaticAdaptation performs chromatic adaptation between white points
func (cc *ColorConverter) chromaticAdaptation(xyz XYZColor, sourceWP, destWP XYZColor) XYZColor {
	// Bradford chromatic adaptation matrix (simplified)
	// In a full implementation, this would use proper Bradford transformation

	scaleX := destWP.X / sourceWP.X
	scaleY := destWP.Y / sourceWP.Y
	scaleZ := destWP.Z / sourceWP.Z

	return XYZColor{
		X: xyz.X * scaleX,
		Y: xyz.Y * scaleY,
		Z: xyz.Z * scaleZ,
	}
}

// Standard ICC profiles

// CreateSRGBProfile creates a standard sRGB ICC profile
func CreateSRGBProfile() *ICCProfile {
	profile := NewICCProfile()

	// sRGB colorants (ITU-R BT.709 primaries)
	profile.RedColorant = XYZColor{X: 0.4124, Y: 0.2126, Z: 0.0193}
	profile.GreenColorant = XYZColor{X: 0.3576, Y: 0.7152, Z: 0.1192}
	profile.BlueColorant = XYZColor{X: 0.1805, Y: 0.0722, Z: 0.9505}
	profile.WhitePoint = IlluminantD65

	// sRGB gamma curves (simplified)
	gamma := 2.4
	profile.Curves = []ToneCurve{
		{Type: 0, Points: []float64{gamma}},
		{Type: 0, Points: []float64{gamma}},
		{Type: 0, Points: []float64{gamma}},
	}

	return profile
}

// CreateAdobeRGBProfile creates an Adobe RGB ICC profile
func CreateAdobeRGBProfile() *ICCProfile {
	profile := NewICCProfile()

	// Adobe RGB colorants
	profile.RedColorant = XYZColor{X: 0.5767, Y: 0.2973, Z: 0.0270}
	profile.GreenColorant = XYZColor{X: 0.1856, Y: 0.6274, Z: 0.0707}
	profile.BlueColorant = XYZColor{X: 0.1882, Y: 0.0753, Z: 0.9911}
	profile.WhitePoint = IlluminantD65

	// Adobe RGB gamma (2.2)
	gamma := 2.2
	profile.Curves = []ToneCurve{
		{Type: 0, Points: []float64{gamma}},
		{Type: 0, Points: []float64{gamma}},
		{Type: 0, Points: []float64{gamma}},
	}

	return profile
}

// Context integration

// SetColorProfile sets the color profile for the context
func (dc *Context) SetColorProfile(profile *ICCProfile) {
	dc.colorProfile = profile
}

// GetColorProfile returns the current color profile
func (dc *Context) GetColorProfile() *ICCProfile {
	if dc.colorProfile == nil {
		dc.colorProfile = CreateSRGBProfile()
	}
	return dc.colorProfile
}

// SetColorConverter sets the color converter for the context
func (dc *Context) SetColorConverter(converter *ColorConverter) {
	dc.colorConverter = converter
}

// ConvertToColorSpace converts the current image to a different color space
func (dc *Context) ConvertToColorSpace(targetProfile *ICCProfile) {
	if dc.colorConverter == nil {
		sourceProfile := dc.GetColorProfile()
		dc.colorConverter = NewColorConverter(sourceProfile, targetProfile)
	}

	bounds := dc.im.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := dc.im.RGBAAt(x, y)
			convertedColor := dc.colorConverter.ConvertColor(originalColor)
			if rgba, ok := convertedColor.(color.RGBA); ok {
				dc.im.SetRGBA(x, y, rgba)
			}
		}
	}
}
