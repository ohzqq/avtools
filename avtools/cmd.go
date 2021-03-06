package avtools

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Cmd struct {
	verbose bool
	cwd     string
	exec    *exec.Cmd
	tmpFile *os.File
}

func NewCmd(cmd *exec.Cmd, verbose bool) *Cmd {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Cmd{
		cwd:     cwd,
		exec:    cmd,
		verbose: verbose,
	}
}

func (cmd *Cmd) tmp(f *os.File) *Cmd {
	cmd.tmpFile = f
	return cmd
}

func (cmd Cmd) Run() []byte {
	if cmd.tmpFile != nil {
		defer os.Remove(cmd.tmpFile.Name())
	}

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmd.exec.Stderr = &stderr
	cmd.exec.Stdout = &stdout

	err := cmd.exec.Run()
	if err != nil {
		log.Fatal("Command finished with error: %v\n", cmd.exec.String())
		fmt.Printf("%v\n", stderr.String())
	}

	if len(stdout.Bytes()) > 0 {
		return stdout.Bytes()
	}

	if cmd.verbose {
		fmt.Println(cmd.exec.String())
	}
	return nil
}
