package file

import (
	"fmt"
	"log"
	"mime"
	"path/filepath"
	"strings"
)

type File struct {
	Abs      string
	Path     string
	Base     string
	Ext      string
	Name     string
	File     string
	Padding  string
	Mimetype string
	name     string
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
	f.Mimetype = mime.TypeByExtension(f.Ext)
	f.Name = strings.TrimSuffix(abs, f.Ext)
	f.name = strings.TrimSuffix(f.Base, f.Ext)

	f.Path, f.File = filepath.Split(abs)

	return f
}

func (f File) AddSuffix(s string) string {
	name := f.name + s + f.Ext
	return filepath.Join(f.Path, name)
}

func (f File) Pad(i int) string {
	p := fmt.Sprintf(f.Padding, i)
	return f.Name + p + f.Ext
}
