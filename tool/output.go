package tool

import (
	"github.com/ohzqq/avtools/file"
)

type Output struct {
	num     int
	Padding string
	padName bool
	file.File
}

func NewOutput(o string) *Output {
	return &Output{
		num:     1,
		padName: Cfg().Defaults.HasPadding(),
		Padding: Cfg().Defaults.Padding,
		File:    file.New(o),
	}
}

func (o Output) String() string {
	if o.padName {
		return o.Pad(o.num)
	}
	return o.Abs
}
