package storage

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

var (
	_fltDig uint = 15
)

func fltDig(precision uint) uint {
	if precision > _fltDig {
		precision = _fltDig
	}

	return precision
}

type sizeType struct {
	unit  string
	value float64
}

var _units = []sizeType{
	{
		unit:  "B",
		value: 1024 * 1024 * 1024 * 1024 * 1024,
	},

	{
		unit:  "KiB",
		value: 1024 * 1024 * 1024 * 1024,
	},

	{
		unit:  "MiB",
		value: 1024 * 1024 * 1024,
	},

	{
		unit:  "GiB",
		value: 1024 * 1024,
	},

	{
		unit:  "TiB",
		value: 1024,
	},

	{
		unit:  "PiB",
		value: 1,
	},
}

func indexOfUnits(unit string) (int, error) {
	for i := 0; i < len(_units); i++ {
		if strings.EqualFold(unit, _units[i].unit) {
			return i, nil
		}
	}

	return 0, fmt.Errorf("invalid size unit '%s'", unit)
}

func unitOfUnits(index int) (string, error) {
	if index < 0 || index > len(_units)-1 {
		return "", fmt.Errorf("invalid size unit index '%d'", index)
	}

	return _units[index].unit, nil
}

// FormatSize supports user friendly formatting, comparing,
// and also converting to kb, mb, gb etc.
type FormatSize struct {
	size  float64
	index int
}

func Format(size float64, unit string) (*FormatSize, error) {
	if size < 0 {
		return nil, errors.New("invalid size number")
	}

	index, err := indexOfUnits(unit)
	if err != nil {
		return nil, err
	}

	return &FormatSize{size: size, index: index}, nil
}

func FormatByte(size float64) *FormatSize {
	f, _ := Format(size, "B")
	return f
}

func (f *FormatSize) Unit() string {
	unit, _ := unitOfUnits(f.index)
	return unit
}

func (f *FormatSize) Truncate(precision uint) float64 {
	precision = fltDig(precision)
	return float64(uint64(f.size*math.Pow(10,
		float64(precision)))) / math.Pow(10, float64(precision))
}

func (f *FormatSize) Round(precision uint) float64 {
	precision = fltDig(precision)
	return float64(uint64(f.size*math.Pow(10,
		float64(precision))+0.5)) / math.Pow(10, float64(precision))
}

func (f *FormatSize) convert(unit string) (*FormatSize, error) {
	index, err := indexOfUnits(unit)
	if err != nil {
		return nil, err
	}

	if index == f.index {
		return f, nil
	}

	var size float64
	if index > f.index {
		size = f.size / math.Pow(1024, float64(index-f.index))
	} else {
		size = f.size * math.Pow(1024, float64(f.index-index))
	}

	return &FormatSize{size: size, index: index}, nil
}

func (f *FormatSize) Convert(unit string, precision uint, doTruncate bool) (float64, string, error) {
	ff, err := f.convert(unit)
	if err != nil {
		return 0, unit, err
	}

	if doTruncate {
		return ff.Truncate(precision), unit, nil
	}

	return ff.Round(precision), unit, nil
}

func (f *FormatSize) FriendlyConvert(precision uint, doTruncate bool) (float64, string, error) {
	i := len(_units) - 1
	for ; i >= 0; i-- {
		if f.size*(_units[i].value/_units[f.index].value) >= 1 {
			break
		}
	}

	if i < 0 {
		return f.Convert(f.Unit(), precision, doTruncate)
	}

	return f.Convert(_units[i].unit, precision, doTruncate)
}

func (f *FormatSize) Compare(x *FormatSize) int {
	ff, _ := f.convert("B")
	xf, _ := x.convert("B")

	if ff.size == xf.size {
		return 0
	} else if ff.size > xf.size {
		return 1
	}

	return -1
}

func (f *FormatSize) Add(x *FormatSize) (float64, string) {
	unit, _ := minUnit(f.index, x.index)
	ff, _ := f.convert(unit)
	xf, _ := x.convert(unit)

	return ff.size + xf.size, unit
}

func (f *FormatSize) Sub(x *FormatSize) (float64, string) {
	unit, _ := minUnit(f.index, x.index)
	ff, _ := f.convert(unit)
	xf, _ := x.convert(unit)

	if ff.size > xf.size {
		return ff.size - xf.size, unit
	}

	return xf.size - ff.size, unit
}

func (f *FormatSize) Show() string {
	size, unit, _ := f.FriendlyConvert(2, false)
	return fmt.Sprintf("%.2f %s", size, unit)
}

func minUnit(x, y int) (string, error) {
	var min int
	if x > y {
		min = y
	} else {
		min = x
	}

	return unitOfUnits(min)
}
