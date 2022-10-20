package ffmpeg

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const (
	ffmpegBin  = `ffmpeg`
	hideBanner = "-hide_banner"
)

type Cmd struct {
	*Args
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

func (c Cmd) String() string {
	args, err := c.ParseArgs()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Join(args, " ")
}

func (c *Cmd) ParseArgs() ([]string, error) {
	args := []string{hideBanner}

	if c.HasLogLevel() {
		args = append(args, c.logLevel...)
	}

	if c.HasPreInput() {
		args = append(args, c.PreInput...)
	}

	if c.HasInput() {
		args = append(args, c.input.Parse()...)
	} else {
		return args, fmt.Errorf("no input file specified")
	}

	if c.HasPostInput() {
		args = append(args, c.PostInput...)
	}

	if c.HasVideoCodec() && !c.streamCopy {
		args = append(args, vCodecFlag, c.videoCodec)
	}

	if c.noVideo {
		args = append(args, noVideo)
	}

	if c.HasVideoParams() {
		args = append(args, c.VideoParams...)
	}

	if c.HasVideoFilters() && !c.HasFilters() && !c.HasFilterGraph() {
		args = append(args, "-vf", c.VideoFilters.String())
	}

	if c.HasAudioCodec() && !c.streamCopy {
		args = append(args, aCodecFlag, c.audioCodec)
	}

	if c.noAudio {
		args = append(args, noAudio)
	}

	//if c.streamCopy && !c.noAudio || !c.noVideo {
	//  args = append(args, "-c", "copy")
	//}

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
