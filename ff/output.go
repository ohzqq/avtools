package ff

import (
	"fmt"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Output struct {
	Args ffmpeg.KwArgs
}

func NewOutput(args ...ffmpeg.KwArgs) Output {
	return Output{
		Args: ffmpeg.MergeKwArgs(args),
	}
}

func (out Output) KwArgs() ffmpeg.KwArgs {
	args := make(ffmpeg.KwArgs)
	for key, val := range out.Args {
		switch key {
		case "ext", "padding", "num", "name":
		default:
			args[key] = val
		}
	}
	return args
}

func (out Output) String() string {
	var name string
	if n := out.Get("name"); n != nil {
		name = n.(string)
	}

	var padding string
	if pad := out.Get("padding"); pad != nil {
		padding = pad.(string)
	}

	var ext string
	if e := out.Get("ext"); e != nil {
		ext = e.(string)
	}

	var num int
	if n := out.Get("num"); n != nil {
		num = n.(int)
	}

	if padding == "" {
		return fmt.Sprintf("%s%s", name, ext)
	}

	name = fmt.Sprintf("%s"+padding+"%s", name, num, ext)

	return name
}

func (out Output) Compile(s *ffmpeg.Stream) *ffmpeg.Stream {
	return s.Output(out.String(), out.KwArgs())
}

func (out *Output) Name(n string) *Output {
	out.Set("name", n)
	return out
}

func (out *Output) Pad(p string) *Output {
	out.Set("padding", p)
	return out
}

func (out *Output) Num(n int) *Output {
	out.Set("num", n)
	return out
}

func (out *Output) Ext(ext string) *Output {
	out.Set("ext", ext)
	return out
}

func (out *Output) Copy() *Output {
	out.Set("c", "copy")
	out.Del("c:a")
	out.Del("c:v")
	return out
}

func (out *Output) IsStreamCopy() bool {
	var aCopy bool
	if ac, ok := out.Args["c:a"]; ok {
		if ac == "copy" {
			aCopy = true
		}
	}

	var vCopy bool
	if vc, ok := out.Args["c:v"]; ok {
		if vc == "copy" {
			vCopy = true
		}
	}

	return aCopy || vCopy
}

func (out *Output) Get(key string) any {
	if val, ok := out.Args[key]; ok {
		return val
	}
	return nil
}

func (out *Output) Set(key string, val any) *Output {
	out.Args[key] = val
	return out
}

func (out *Output) Del(key string) *Output {
	delete(out.Args, key)
	return out
}

func (o *Output) VideoCodec(val string) *Output {
	o.Set("c:v", val)
	return o
}

func (o *Output) VideoParams(args map[string]any) *Output {
	for key, val := range args {
		o.Set(key, val)
	}
	return o
}

func (o *Output) AudioCodec(val string) *Output {
	o.Set("c:a", val)
	return o
}

func (o *Output) AudioParams(args map[string]any) *Output {
	for key, val := range args {
		o.Set(key, val)
	}
	return o
}
