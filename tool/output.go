package tool

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

type Output struct {
	num     int
	name    string
	Ext     string
	Abs     string
	path    string
	Padding string
	Pad     bool
}

func NewOutput(o string) *Output {
	abs, err := filepath.Abs(o)
	if err != nil {
		log.Fatal(err)
	}
	dir, file := filepath.Split(abs)
	ext := filepath.Ext(file)
	return &Output{
		num:     1,
		path:    dir,
		Abs:     abs,
		Pad:     Cfg().Defaults.HasPadding(),
		Padding: Cfg().Defaults.Padding,
		Ext:     ext,
		name:    strings.TrimSuffix(file, ext),
	}
}

func (o Output) Name() string {
	return strings.Join([]string{o.name, o.Ext}, "")
}

func (o Output) String() string {
	name := o.name
	if o.Pad {
		name = o.name + fmt.Sprintf(o.Padding, o.num)
	}
	return filepath.Join(o.path, name+o.Ext)
}
