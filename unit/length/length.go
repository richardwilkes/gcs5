package length

import (
	"fmt"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Units holds the length unit type. Note that conversions to/from metric are done using the simplified GURPS metric
// conversion of 1 yd = 1 meter. For consistency, all metric lengths are converted to meters, then to yards, rather than
// the variations at different lengths that the GURPS rules suggest.
type Units string

// Possible Units values.
const (
	Centimeter    = Units("cm")
	Inch          = Units("in")
	Feet          = Units("ft")
	Yard          = Units("yd")
	Meter         = Units("m")
	Kilometer     = Units("km")
	Mile          = Units("mi")
	FeetAndInches = Units(`'"`) // This one is special and not a suffix
)

// Length contains a fixed-point value in inches.
type Length fixed.F64d4

// FromInt64 creates a new Length.
func FromInt64(value int64, unit Units) Length {
	return convertToInches(fixed.F64d4FromInt64(value), unit)
}

// FromFloat64 creates a new Length.
func FromFloat64(value float64, unit Units) Length {
	return convertToInches(fixed.F64d4FromFloat64(value), unit)
}

// FromString creates a new Length. May have any of the known Units suffixes, a feet and inches format (e.g. 6'2"), or no
// notation at all, in which case defaultUnits is used.
func FromString(text string, defaultUnits Units) (Length, error) {
	text = strings.TrimLeft(strings.TrimSpace(text), "+")
	for _, unit := range []Units{
		Centimeter,
		Inch,
		Feet,
		Yard,
		Meter,
		Kilometer,
		Mile,
	} {
		if strings.HasSuffix(text, string(unit)) {
			value, err := fixed.F64d4FromString(strings.TrimSpace(strings.TrimSuffix(text, string(unit))))
			if err != nil {
				return 0, err
			}
			return convertToInches(value, unit), nil
		}
	}
	// Didn't match any of the Units types, let's try feet & inches
	feetIndex := strings.Index(text, "'")
	inchIndex := strings.Index(text, `"`)
	if feetIndex == -1 && inchIndex == -1 {
		// Nope, so let's use our passed-in default units
		value, err := fixed.F64d4FromString(strings.TrimSpace(text))
		if err != nil {
			return 0, err
		}
		return convertToInches(value, defaultUnits), nil
	}
	var feet, inches fixed.F64d4
	var err error
	if feetIndex != -1 {
		s := strings.TrimSpace(text[:feetIndex])
		feet, err = fixed.F64d4FromString(s)
		if err != nil {
			return 0, err
		}
	}
	if inchIndex != -1 {
		if feetIndex > inchIndex {
			return 0, errs.New(fmt.Sprintf("invalid format: %s", text))
		}
		s := strings.TrimSpace(text[feetIndex+1 : inchIndex])
		inches, err = fixed.F64d4FromString(s)
		if err != nil {
			return 0, err
		}
	}
	return Length(feet.Mul(fixed.F64d4FromInt64(12)) + inches), nil
}

func convertToInches(value fixed.F64d4, unit Units) Length {
	switch unit {
	case Centimeter:
		value = value.Mul(fixed.F64d4FromFloat64(36)).Div(fixed.F64d4FromInt64(100))
	case Feet:
		value = value.Mul(fixed.F64d4FromInt64(12))
	case Yard, Meter:
		value = value.Mul(fixed.F64d4FromInt64(36))
	case Kilometer:
		value = value.Mul(fixed.F64d4FromInt64(36000))
	case Mile:
		value = value.Mul(fixed.F64d4FromInt64(63360))
	default: // Same as Inch
	}
	return Length(value)
}

func (l Length) String() string {
	return l.Format(FeetAndInches)
}

// Format the length as the given units.
func (l Length) Format(unit Units) string {
	inches := fixed.F64d4(l)
	switch unit {
	case Centimeter:
		return inches.Div(fixed.F64d4FromInt64(36)).Mul(fixed.F64d4FromInt64(100)).String() + " " + string(unit)
	case Feet:
		return inches.Div(fixed.F64d4FromInt64(12)).String() + " " + string(unit)
	case Yard, Meter:
		return inches.Div(fixed.F64d4FromInt64(36)).String() + " " + string(unit)
	case Kilometer:
		return inches.Div(fixed.F64d4FromInt64(36000)).String() + " " + string(unit)
	case Mile:
		return inches.Div(fixed.F64d4FromInt64(5280)).String() + " " + string(unit)
	case FeetAndInches:
		oneFoot := fixed.F64d4FromInt64(12)
		feet := inches.Div(oneFoot).Trunc()
		inches -= feet.Mul(oneFoot)
		if feet == 0 && inches == 0 {
			return "0'"
		}
		var buffer strings.Builder
		if feet > 0 {
			buffer.WriteString(feet.String())
			buffer.WriteByte('\'')
		}
		if inches > 0 {
			buffer.WriteString(inches.String())
			buffer.WriteByte('"')
		}
		return buffer.String()
	default: // Same as Inch
		return inches.String() + " " + string(Inch)
	}
}
