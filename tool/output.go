package tool

import (
	"fmt"
	"path/filepath"
	"strings"
)

type output struct {
	num     int
	name    string
	ext     string
	path    string
	padding string
	Pad     bool
}

func NewOutput() *output {
	return &output{
		num:     1,
		Pad:     true,
		padding: "%03d",
		name:    "tmp",
	}
}

func OutputFromInput(i string) *output {
	out := NewOutput()
	dir, file := filepath.Split(i)
	out.path = dir
	out.ext = filepath.Ext(file)
	out.name = strings.TrimSuffix(file, out.ext)
	return out
}

func (o output) String() string {
	if o.Pad {
		o.name = o.name + fmt.Sprintf(o.padding, o.num)
	}

	return filepath.Join(o.path, o.name+o.ext)
}
