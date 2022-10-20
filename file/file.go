package file

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

type File struct {
	Abs     string
	Path    string
	Base    string
	Ext     string
	Name    string
	File    string
	Padding string
}

func New(n string) File {
	abs, err := filepath.Abs(n)
	if err != nil {
		log.Fatal(err)
	}

	f := File{
		Base:    filepath.Base(abs),
		Ext:     filepath.Ext(abs),
		Abs:     abs,
		Padding: "%03d",
	}
	f.Name = strings.TrimSuffix(abs, f.Ext)

	f.Path, f.File = filepath.Split(abs)

	return f
}

func (f File) Pad(i int) string {
	p := fmt.Sprintf(f.Padding, i)
	return f.Name + p + f.Ext
}
