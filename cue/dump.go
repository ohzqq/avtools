package cue

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"os"
)

func (s Sheet) Dump() []byte {
	var (
		tmpl = template.Must(template.New("cue").Funcs(tmplFuncs).Parse(cueTmpl))
		buf  bytes.Buffer
	)

	if s.Audio == "" {
		s.Audio = "tmp"
	}

	err := tmpl.Execute(&buf, s)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func (s Sheet) Write(wr io.Writer) error {
	_, err := wr.Write(s.Dump())
	if err != nil {
		return err
	}
	return nil
}

func (s Sheet) Save() error {
	return s.SaveAs(s.Audio)
}

func (s Sheet) SaveAs(name string) error {
	if name == "" || s.Audio == "" {
		name = "tmp"
	}

	file, err := os.Create(name + ".cue")
	if err != nil {
		return err
	}
	defer file.Close()

	err = s.Write(file)
	if err != nil {
		return err
	}

	return nil
}

var tmplFuncs = template.FuncMap{
	"inc": Inc,
}

func Inc(n int) int {
	return n + 1
}

const cueTmpl = `FILE "{{.File}}" {{.Ext -}}
{{range $idx, $ch := .Tracks}}
TRACK {{inc $idx}} AUDIO
  TITLE "{{$ch.Title}}"
  INDEX 01 {{$ch.Stamp}}
{{- end -}}`
