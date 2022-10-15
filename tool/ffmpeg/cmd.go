package ffmpeg

import (
	"fmt"
	"log"
	"strings"
)

type Cmd struct {
	*Args
	args []string
}

func New() *Cmd {
	return &Cmd{
		Args: &Args{
			VideoCodec: []string{"-c:v"},
			AudioCodec: []string{"-c:a"},
			//VideoFilters:  []string{"-vf"},
			//AudioFilters:  []string{"-af"},
			//FilterComplex: []string{"-filter_complex"},
			LogLevel: []string{"-loglevel"},
		},
	}
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
		args = append(args, c.LogLevel...)
	}

	if c.HasPreInput() {
		args = append(args, c.PreInput...)
	}

	if c.HasInput() {
		args = append(args, "-i", c.Input)
	} else {
		return args, fmt.Errorf("no input file specified")
	}

	if c.HasPostInput() {
		args = append(args, c.PostInput...)
	}

	if c.HasVideoCodec() {
		args = append(args, c.VideoCodec...)
	}

	if c.HasVideoParams() {
		args = append(args, c.VideoParams...)
	}

	if c.HasVideoFilters() && !c.HasFilters() {
		args = append(args, "-vf", c.VideoFilters.String())
	}

	if c.HasAudioCodec() {
		args = append(args, c.AudioCodec...)
	}

	if c.HasAudioFilters() && !c.HasFilters() {
		args = append(args, "-af", c.AudioFilters.String())
	}

	if c.HasAudioParams() {
		args = append(args, c.AudioParams...)
	}

	if c.HasFilters() {
		args = append(args, "-filter_complex", c.FilterComplex.String())
	}

	if c.HasMiscParams() {
		args = append(args, c.MiscParams...)
	}

	if c.Output != "" {
		args = append(args, c.Output)
	}

	return args, nil
}
