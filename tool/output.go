package tool

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Output struct {
	num     int
	Name    string
	Ext     string
	path    string
	Padding string
	Pad     bool
}

func NewOutput(o string) *Output {
	dir, file := filepath.Split(o)
	ext := filepath.Ext(file)
	return &Output{
		num:     1,
		path:    dir,
		Pad:     Cfg().Defaults.HasPadding(),
		Padding: Cfg().Defaults.Padding,
		Ext:     ext,
		Name:    strings.TrimSuffix(file, ext),
	}
}

func (o Output) String() string {
	name := o.Name
	if o.Pad {
		name = o.Name + fmt.Sprintf(o.Padding, o.num)
	}
	return filepath.Join(o.path, name+o.Ext)
}
