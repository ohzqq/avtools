package cue

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"os"
)

func (c CueSheet) Dump() []byte {
	var (
		tmpl = template.Must(template.New("cue").Funcs(tmplFuncs).Parse(cueTmpl))
		buf  bytes.Buffer
	)

	if c.Audio == "" {
		c.Audio = "tmp"
	}

	err := tmpl.Execute(&buf, c)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func (c CueSheet) Write(wr io.Writer) error {
	_, err := wr.Write(c.Dump())
	if err != nil {
		return err
	}
	return nil
}

func (c CueSheet) Save() error {
	return c.SaveAs(c.Audio)
}

func (c CueSheet) SaveAs(name string) error {
	if name == "" || c.Audio == "" {
		name = "tmp"
	}

	file, err := os.Create(name + ".cue")
	if err != nil {
		return err
	}
	defer file.Close()

	err = c.Write(file)
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
