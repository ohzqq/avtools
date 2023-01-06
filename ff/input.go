package ff

import ffmpeg "github.com/u2takey/ffmpeg-go"

type Input struct {
	File string
	Args ffmpeg.KwArgs
}

func NewInput(args ...ffmpeg.KwArgs) Input {
	init := ffmpeg.KwArgs{
		"loglevel":    "error",
		"hide_banner": "",
	}
	in := []ffmpeg.KwArgs{init}
	in = append(in, args...)
	return Input{
		Args: ffmpeg.MergeKwArgs(in),
	}
}

func (i *Input) Merge(kwargs ffmpeg.KwArgs) *Input {
	args := []ffmpeg.KwArgs{i.Args, kwargs}
	i.Args = ffmpeg.MergeKwArgs(args)
	return i
}

func (i *Input) Compile(file string) *ffmpeg.Stream {
	return ffmpeg.Input(file, i.Args)
}

func (i *Input) Verbose() *Input {
	i.Set("loglevel", "info")
	return i
}

func (i *Input) Overwrite() *Input {
	i.Set("y", "")
	return i
}

func (i *Input) Start(ss string) *Input {
	i.Set("ss", ss)
	return i
}

func (i *Input) End(to string) *Input {
	i.Set("to", to)
	return i
}

func (i *Input) Set(key string, val any) {
	i.Args[key] = val
}
