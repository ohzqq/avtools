package ffmpeg

type Cmd struct {
	*Args
	args []string
}

func New() *Cmd {
	return &Cmd{
		Args: &Args{
			VideoCodec:    []string{"-c:v"},
			AudioCodec:    []string{"-c:a"},
			VideoFilters:  []string{"-vf"},
			AudioFilters:  []string{"-af"},
			FilterComplex: []string{"-filter_complex"},
			LogLevel:      []string{"-loglevel"},
		},
	}
}

func (c *Cmd) ParseArgs() *Cmd {
	if c.HasLogLevel() {
		c.args = append(c.args, c.LogLevel...)
	}

	if c.HasPreInput() {
		c.args = append(c.args, c.PreInput...)
	}

	if c.HasInput() {
		c.args = append(c.args, c.Input.Map()...)
	}

	if c.HasPostInput() {
		c.args = append(c.args, c.PostInput...)
	}

	if c.HasVideoCodec() {
		c.args = append(c.args, c.VideoCodec...)
	}

	if c.HasVideoParams() {
		c.args = append(c.args, c.VideoParams...)
	}

	if c.HasVideoFilters() && !c.HasFilterComplex() {
		c.args = append(c.args, c.VideoFilters.String())
	}

	if c.HasAudioCodec() {
		c.args = append(c.args, c.AudioCodec...)
	}

	if c.HasAudioFilters() && !c.HasFilterComplex() {
		c.args = append(c.args, c.AudioFilters.String())
	}

	if c.HasAudioParams() {
		c.args = append(c.args, c.AudioParams...)
	}

	if c.HasFilterComplex() {
		c.args = append(c.args, c.FilterComplex.String())
	}

	if c.HasMiscParams() {
		c.args = append(c.args, c.MiscParams...)
	}

	if c.Output != "" {
		c.args = append(c.args, c.Output)
	}

	return c
}
