package avtools

import (
	"path/filepath"
	"log"
	"fmt"
	//"os"
	//"bytes"
	//"encoding/json"
	//"strconv"
	//"strings"
)
var _ = fmt.Printf

type Media struct {
	Overwrite bool
	File string
	Path string
	Dir string
	Ext string
	//Meta *MediaMeta
}

func NewMedia(input string) *Media {
	m := Media{}

	abs, err := filepath.Abs(input)
	if err != nil {
		log.Fatal(err)
	}

	m.Path = abs
	m.File = filepath.Base(input)
	m.Dir = filepath.Dir(input)
	m.Ext = filepath.Ext(input)
	//m.ParseJsonMeta()

	return &m
}

