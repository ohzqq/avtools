package main

import (
	"log"

	"github.com/ohzqq/avtools/cmd/yygif/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cmd.Execute()
}
