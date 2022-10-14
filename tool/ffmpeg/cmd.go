package ffmpeg

type Cmd struct {
	Input         []string
	PreInput      mapArgs
	PostInput     mapArgs
	VideoCodec    string
	VideoParams   mapArgs
	VideoFilters  stringArgs
	AudioCodec    string
	AudioParams   mapArgs
	AudioFilters  stringArgs
	FilterComplex stringArgs
	MiscParams    stringArgs
	LogLevel      string
	Name          string
	Padding       string
	Ext           string
}

type mapArgs []map[string]string

type stringArgs []string

func New() *Cmd {
	return &Cmd{}
}
