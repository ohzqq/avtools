package ffmpeg

const (
	aCodecFlag   = `-c:a`
	vCodecFlag   = `-c:v`
	logLevelFlag = `-loglevel`
	noAudio      = `-an`
	noVideo      = `-vn`
)

type Args struct {
	logLevel      string
	PreInput      []string
	input         Input
	PostInput     []string
	streamCopy    bool
	noAudio       bool
	noVideo       bool
	videoCodec    string
	VideoParams   []string
	VideoFilters  Filters
	audioCodec    string
	AudioParams   []string
	AudioFilters  Filters
	FilterComplex Filters
	MiscParams    []string
	Metadata      map[string]string
	output        string
	filters       string
}

func NewArgs() *Args {
	return &Args{
		Metadata: make(map[string]string),
	}
}

func (c Args) HasLogLevel() bool {
	return c.logLevel != ""
}

func (c *Args) LogLevel(l string) *Args {
	c.logLevel = l
	return c
}

func (c *Args) Overwrite() *Args {
	c.AppendPreInput("y")
	return c
}

func (c Args) HasPreInput() bool {
	return len(c.PreInput) > 0
}

func (c *Args) AppendPreInput(flag string, val ...string) *Args {
	c.PreInput = append(c.PreInput, "-"+flag)
	c.PreInput = append(c.PreInput, val...)
	return c
}

func (c *Args) SS(t string) *Args {
	c.AppendPreInput("ss", t)
	return c
}

func (c *Args) To(t string) *Args {
	c.AppendPreInput("to", t)
	return c
}

func (c Args) HasInput() bool {
	return len(c.input.files) > 0
}

func (c *Args) Input(i string) *Args {
	c.input.Add(i)
	return c
}

func (c *Args) FFmeta(i string) *Args {
	c.input.FFmetadata = i
	return c
}

func (c *Args) HasChapters() *Args {
	c.input.HasChapters = true
	return c
}

func (c Args) HasPostInput() bool {
	return len(c.PostInput) > 0
}

func (c *Args) AppendPostInput(flag string, val ...string) *Args {
	c.PostInput = append(c.PostInput, "-"+flag)
	c.PostInput = append(c.PostInput, val...)
	return c
}

func (c *Args) Stream() *Args {
	c.streamCopy = true
	return c
}

func (c Args) HasVideoCodec() bool {
	if c.noVideo {
		return false
	}
	return c.videoCodec != ""
}

func (c *Args) SetVideoCodec(codec string) *Args {
	c.videoCodec = codec
	return c
}

func (c *Args) CV(codec string) *Args {
	c.videoCodec = codec
	return c
}

func (c *Args) VN() *Args {
	c.videoCodec = ""
	c.noVideo = true
	return c
}

func (c Args) HasVideoParams() bool {
	return len(c.VideoParams) > 0
}

func (c *Args) AppendVideoParam(flag string, val ...string) *Args {
	c.VideoParams = append(c.VideoParams, "-"+flag)
	c.VideoParams = append(c.VideoParams, val...)
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
	if c.noAudio {
		return false
	}
	return c.audioCodec != ""
}

func (c *Args) SetAudioCodec(codec string) *Args {
	c.audioCodec = codec
	return c
}

func (c *Args) CA(codec string) *Args {
	c.audioCodec = codec
	return c
}

func (c *Args) AN() *Args {
	c.audioCodec = ""
	c.noAudio = true
	return c
}
func (c Args) HasAudioParams() bool {
	return len(c.AudioParams) > 0
}

func (c *Args) AppendAudioParam(flag string, val ...string) *Args {
	c.AudioParams = append(c.AudioParams, "-"+flag)
	c.AudioParams = append(c.AudioParams, val...)
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

func (c *Args) AppendMiscParam(flag string, val ...string) *Args {
	c.MiscParams = append(c.MiscParams, "-"+flag)
	c.MiscParams = append(c.MiscParams, val...)
	return c
}

func (c Args) HasOutput() bool {
	return c.output != ""
}

func (c *Args) Output(o string) *Args {
	c.output = o
	return c
}
