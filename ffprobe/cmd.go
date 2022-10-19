package ffprobe

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const (
	ffprobeBin    = "ffprobe"
	hideBanner    = "-hide_banner"
	verbose       = "-v"
	showEntries   = "-show_entries"
	selectStreams = "-select_streams"
	writer        = "-of"
	pretty        = "-pretty"
	showChapters  = "-show_chapters"
)

type Cmd struct {
	*Args
}

func New() *Cmd {
	return &Cmd{
		Args: NewArgs(),
	}
}

func (c Cmd) Run() ([]byte, error) {
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)

	cmd, err := c.Build()
	if err != nil {
		return stderr.Bytes(), fmt.Errorf("Cmd failed to build: %v\n", cmd.String())
	}

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err = cmd.Run()

	if err != nil {
		return stderr.Bytes(), fmt.Errorf("%v\n", stderr.String())
	}

	return stdout.Bytes(), nil
}

func (c Cmd) Build() (*exec.Cmd, error) {
	args, err := c.ParseArgs()
	if err != nil {
		return nil, err
	}
	return exec.Command(ffprobeBin, args...), nil
}

func (c Cmd) String() string {
	args, err := c.ParseArgs()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Join(args, " ")
}

func (c Cmd) ParseArgs() ([]string, error) {
	args := []string{hideBanner}

	if c.HasLogLevel() {
		args = append(args, c.logLevel...)
	}

	if c.pretty {
		args = append(args, pretty)
	}

	if c.showChaps {
		args = append(args, showChapters)
	}

	if c.HasStreams() {
		for _, s := range c.streams {
			args = append(args, selectStreams, s)
		}
	}

	if c.HasEntries() {
		args = append(args, showEntries, c.entries.String())
	}

	if c.HasFormat() {
		args = append(args, c.format...)
	}

	if c.HasInput() {
		args = append(args, c.input)
	} else {
		return args, fmt.Errorf("no input file specified")
	}

	return args, nil
}
