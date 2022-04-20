package fftools

import (
	"bufio"
	"log"
	"os"
)

func ReadFile(file string) *bufio.Scanner {
	contents, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer contents.Close()

	return bufio.NewScanner(contents)
}
