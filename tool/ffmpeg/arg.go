package ffmpeg

type Args struct {
	logLevel      []string
	PreInput      []string
	input         string
	PostInput     []string
	videoCodec    []string
	VideoParams   []string
	VideoFilters  Filters
	audioCodec    []string
	AudioParams   []string
	AudioFilters  Filters
	FilterComplex Filters
	MiscParams    []string
	Metadata      map[string]string
	output        string
	filters       string
}

func (c Args) HasLogLevel() bool {
	return len(c.logLevel) > 1
}

func (c *Args) LogLevel(l string) *Args {
	c.logLevel = append(c.logLevel, l)
	return c
}

func (c Args) HasPreInput() bool {
	return len(c.PreInput) > 0
}

func (c *Args) AppendPreInput(flag, val string) *Args {
	c.PreInput = append(c.PreInput, "-"+flag, val)
	return c
}

func (c Args) HasInput() bool {
	return c.input != ""
}

func (c *Args) Input(i string) *Args {
	c.input = i
	return c
}

func (c Args) HasPostInput() bool {
	return len(c.PreInput) > 0
}

func (c *Args) AppendPostInput(flag, val string) *Args {
	c.PostInput = append(c.PostInput, "-"+flag, val)
	return c
}

func (c Args) HasVideoCodec() bool {
	return len(c.videoCodec) > 1
}

func (c *Args) SetVideoCodec(codec string) *Args {
	c.videoCodec = append(c.videoCodec, codec)
	return c
}

func (c *Args) CV(codec string) *Args {
	c.videoCodec = append(c.videoCodec, codec)
	return c
}

func (c Args) HasVideoParams() bool {
	return len(c.VideoParams) > 0
}

func (c *Args) AppendVideoParam(flag, val string) *Args {
	c.VideoParams = append(c.VideoParams, "-"+flag, val)
	return c
}

func (c Args) HasVideoFilters() bool {
	return len(c.VideoFilters) > 0
}

func (c *Args) AppendVideoFilter(f Filter) *Args {
	c.VideoFilters = append(c.VideoFilters, f)
	return c
}

func (c *Args) VF(f Filter) *Args {
	c.VideoFilters = append(c.VideoFilters, f)
	return c
}

func (c Args) HasAudioCodec() bool {
	return len(c.audioCodec) > 1
}

func (c *Args) SetAudioCodec(codec string) *Args {
	c.audioCodec = append(c.audioCodec, codec)
	return c
}

func (c *Args) CA(codec string) *Args {
	c.audioCodec = append(c.audioCodec, codec)
	return c
}

func (c Args) HasAudioParams() bool {
	return len(c.AudioParams) > 0
}

func (c *Args) AppendAudioParam(flag, val string) *Args {
	c.AudioParams = append(c.AudioParams, "-"+flag, val)
	return c
}

func (c Args) HasAudioFilters() bool {
	return len(c.AudioFilters) > 0
}

func (c *Args) AppendAudioFilter(f Filter) *Args {
	c.AudioFilters = append(c.AudioFilters, f)
	return c
}

func (c *Args) AF(f Filter) *Args {
	c.AudioFilters = append(c.AudioFilters, f)
	return c
}

func (c Args) HasFilters() bool {
	return len(c.FilterComplex) > 0
}

func (c *Args) AppendFilter(f Filter) *Args {
	c.FilterComplex = append(c.FilterComplex, f)
	return c
}

func (c *Args) Filter(f Filter) *Args {
	c.FilterComplex = append(c.FilterComplex, f)
	return c
}

func (c Args) HasFilterGraph() bool {
	return c.filters != ""
}

func (c *Args) SetFilterGraph(f string) *Args {
	c.filters = f
	return c
}

func (c Args) HasMetadata() bool {
	return len(c.Metadata) > 0
}

func (c *Args) SetMetadata(key, val string) *Args {
	c.Metadata[key] = val
	return c
}

func (c Args) HasMiscParams() bool {
	return len(c.MiscParams) > 0
}

func (c *Args) AppendMiscParam(flag, val string) *Args {
	c.MiscParams = append(c.MiscParams, "-"+flag, val)
	return c
}

func (c Args) HasOutput() bool {
	return c.output != ""
}

func (c *Args) Output(o string) *Args {
	c.output = o
	return c
}
