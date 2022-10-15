package ffmpeg

type Args struct {
	LogLevel      []string
	PreInput      []string
	Input         string
	PostInput     []string
	VideoCodec    []string
	VideoParams   []string
	VideoFilters  Filters
	AudioCodec    []string
	AudioParams   []string
	AudioFilters  Filters
	FilterComplex Filters
	MiscParams    []string
	Output        string
}

func (c *Args) AppendPreInput(flag, val string) *Args {
	c.PreInput = append(c.PreInput, "-"+flag, val)
	return c
}

func (c *Args) SetInput(i string) *Args {
	c.Input = i
	return c
}

func (c Args) HasPreInput() bool {
	return len(c.PreInput) > 0
}

func (c Args) HasPostInput() bool {
	return len(c.PreInput) > 0
}

func (c Args) HasMiscParams() bool {
	return len(c.MiscParams) > 0
}

func (c Args) HasAudioParams() bool {
	return len(c.AudioParams) > 0
}

func (c Args) HasVideoParams() bool {
	return len(c.VideoParams) > 0
}

func (c Args) HasAudioFilters() bool {
	return len(c.AudioFilters) > 1
}

func (c Args) HasVideoFilters() bool {
	return len(c.VideoFilters) > 0
}

func (c Args) HasInput() bool {
	return c.Input != ""
}

func (c Args) HasFilters() bool {
	return len(c.FilterComplex) > 1
}

func (c Args) HasAudioCodec() bool {
	return len(c.AudioCodec) > 1
}

func (c Args) HasVideoCodec() bool {
	return len(c.VideoCodec) > 1
}

func (c Args) HasLogLevel() bool {
	return len(c.LogLevel) > 1
}

func (c *Args) AppendPostInput(flag, val string) *Args {
	c.PostInput = append(c.PostInput, "-"+flag, val)
	return c
}

func (c *Args) AppendVideoParam(flag, val string) *Args {
	c.VideoParams = append(c.VideoParams, "-"+flag, val)
	return c
}

func (c *Args) AppendAudioParam(flag, val string) *Args {
	c.AudioParams = append(c.AudioParams, "-"+flag, val)
	return c
}

func (c *Args) AppendVideoFilter(f Filter) *Args {
	c.VideoFilters = append(c.VideoFilters, f)
	return c
}

func (c *Args) AppendAudioFilter(f Filter) *Args {
	c.AudioFilters = append(c.AudioFilters, f)
	return c
}

func (c *Args) AppendFilter(f Filter) *Args {
	c.FilterComplex = append(c.FilterComplex, f)
	return c
}

func (c *Args) SetVideoCodec(codec string) *Args {
	c.VideoCodec = append(c.VideoCodec, codec)
	return c
}

func (c *Args) SetAudioCodec(codec string) *Args {
	c.AudioCodec = append(c.AudioCodec, codec)
	return c
}

func (c *Args) SetLogLevel(l string) *Args {
	c.LogLevel = append(c.LogLevel, l)
	return c
}

func (c *Args) SetOutput(o string) *Args {
	c.Output = o
	return c
}
