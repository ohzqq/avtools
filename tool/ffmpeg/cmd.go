package ffmpeg

import (
	"strconv"
)

type Cmd struct {
	Input         Input
	PreInput      []string
	PostInput     []string
	VideoCodec    []string
	VideoParams   []string
	VideoFilters  Filters
	AudioCodec    []string
	AudioParams   []string
	AudioFilters  Filters
	FilterComplex Filters
	MiscParams    []string
	LogLevel      []string
	Output        string
	Name          string
	Padding       string
	Ext           string
}

type Input struct {
	input []string
	Meta  string
}

func (i Input) Map() []string {
	total := len(i.input)

	var input []string
	for _, in := range i.input {
		input = append(input, "-i", in)
	}

	if i.Meta != "" {
		input = append(input, "-i", i.Meta)
	}

	if total > 1 || i.Meta != "" {
		for idx, _ := range i.input {
			input = append(input, "-map", strconv.Itoa(idx)+":0")
		}
	}

	if i.Meta != "" {
		input = append(input, "-map_metadata", strconv.Itoa(total))
	}

	return input
}

func New() *Cmd {
	return &Cmd{
		VideoCodec:    []string{"-c:v"},
		AudioCodec:    []string{"-c:a"},
		VideoFilters:  []string{"-vf"},
		AudioFilters:  []string{"-af"},
		FilterComplex: []string{"-filter_complex"},
		LogLevel:      []string{"-loglevel"},
	}
}

func (c *Cmd) AppendPreInput(flag, val string) *Cmd {
	c.PreInput = append(c.PreInput, flag, val)
	return c
}

func (c Cmd) HasPreInput() bool {
	return len(c.PreInput) > 0
}

func (c Cmd) HasPostInput() bool {
	return len(c.PreInput) > 0
}

func (c Cmd) HasMiscParams() bool {
	return len(c.MiscParams) > 0
}

func (c Cmd) HasAudioParams() bool {
	return len(c.AudioParams) > 0
}

func (c Cmd) HasVideoParams() bool {
	return len(c.VideoParams) > 0
}

func (c Cmd) HasAudioFilters() bool {
	return len(c.AudioFilters) > 1
}

func (c Cmd) HasVideoFilters() bool {
	return len(c.VideoFilters) > 1
}

func (c Cmd) HasFilters() bool {
	return len(c.Filters) > 1
}

func (c Cmd) HasAudioCodec() bool {
	return len(c.AudioCodec) > 1
}

func (c Cmd) HasVideoCodec() bool {
	return len(c.VideoCodec) > 1
}

func (c *Cmd) AppendPostInput(flag, val string) *Cmd {
	c.PostInput = append(c.PostInput, flag, val)
	return c
}

func (c *Cmd) AppendVideoParam(flag, val string) *Cmd {
	c.VideoParams = append(c.VideoParams, flag, val)
	return c
}

func (c *Cmd) AppendAudioParam(flag, val string) *Cmd {
	c.AudioParams = append(c.AudioParams, flag, val)
	return c
}

func (c *cmd) AppendVideoFilter(f Filter) *cmd {
	c.VideoFilters = append(c.VideoFilters, f)
	return c
}

func (c *Cmd) AppendAudioFilter(f Filter) *Cmd {
	c.AudioFilters = append(c.AudioFilters, f)
	return c
}

func (c *Cmd) AppendFilter(f Filter) *Cmd {
	c.FilterComplex = append(c.FilterComplex, f)
	return c
}

func (c *Cmd) SetVideoCodec(codec string) *Cmd {
	c.VideoCodec = append(c.VideoCodec, codec)
	return c
}

func (c *Cmd) SetVideoCodec(codec string) *Cmd {
	c.AudioCodec = append(c.AudioCodec, codec)
	return c
}

func (c *Cmd) AppendVideoParam(flag, val string) *Cmd {
	c.VideoParams = append(c.VideoParams, flag, val)
	return c
}

func (c *Cmd) AppendAudioParam(flag, val string) *Cmd {
	c.AudioParams = append(c.AudioParams, flag, val)
	return c
}

func (c *Cmd) SetLogLevel(l string) *Cmd {
	c.LogLevel = append(c.LogLevel, l)
	return c
}

func (c *Cmd) SetOutput(o string) *Cmd {
	c.Output = o
	return c
}
