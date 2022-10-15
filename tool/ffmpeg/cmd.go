package ffmpeg

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const ffmpegBin = `ffmpeg`

type Cmd struct {
	*Args
	args []string
}

func New() *Cmd {
	return &Cmd{
		Args: NewArgs(),
	}
}

func (c Cmd) Build() (*exec.Cmd, error) {
	args, err := c.ParseArgs()
	if err != nil {
		return nil, err
	}
	return exec.Command(ffmpegBin, args...), nil
}

func (c Cmd) String() string {
	args, err := c.ParseArgs()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Join(args, " ")
}

func (c *Cmd) ParseArgs() ([]string, error) {
	var args []string
	if c.HasLogLevel() {
		args = append(args, c.logLevel...)
	}

	if c.HasPreInput() {
		args = append(args, c.PreInput...)
	}

	if c.HasInput() {
		args = append(args, "-i", c.input)
	} else {
		return args, fmt.Errorf("no input file specified")
	}

	if c.HasPostInput() {
		args = append(args, c.PostInput...)
	}

	if c.HasVideoCodec() {
		args = append(args, c.videoCodec...)
	}

	if c.HasVideoParams() {
		args = append(args, c.VideoParams...)
	}

	if c.HasVideoFilters() && !c.HasFilters() && !c.HasFilterGraph() {
		args = append(args, "-vf", c.VideoFilters.String())
	}

	if c.HasAudioCodec() {
		args = append(args, c.audioCodec...)
	}

	if c.HasAudioFilters() && !c.HasFilters() && !c.HasFilterGraph() {
		args = append(args, "-af", c.AudioFilters.String())
	}

	if c.HasAudioParams() {
		args = append(args, c.AudioParams...)
	}

	if c.HasFilters() && !c.HasFilterGraph() {
		args = append(args, "-filter_complex", c.FilterComplex.String())
	}

	if c.HasMetadata() {
		for key, val := range c.Metadata {
			args = append(args, "-metadata", key+"="+val)
		}
	}

	if c.HasMiscParams() {
		args = append(args, c.MiscParams...)
	}

	if c.HasFilterGraph() {
		args = append(args, "-filter_complex", c.filters)
	}

	if c.output != "" {
		args = append(args, c.output)
	} else {
		return args, fmt.Errorf("no output file specified")
	}

	return args, nil
}
