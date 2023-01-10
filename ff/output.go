package ff

import (
	"fmt"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Output struct {
	name    string
	ext     string
	padding string
	num     int
	Args    ffmpeg.KwArgs
}

func NewOutput(args ...ffmpeg.KwArgs) Output {
	var out Output

	for _, a := range args {
		if n, ok := a["name"]; ok {
			out.name = n.(string)
		}

		if pad, ok := a["padding"]; ok {
			out.padding = pad.(string)
		}

		if e, ok := a["ext"]; ok {
			out.ext = e.(string)
		}

		if n, ok := a["num"]; ok {
			out.num = n.(int)
		}
	}

	out.Args = ffmpeg.MergeKwArgs(args)

	return out
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
	if n := out.Get("name"); n != nil {
		out.name = n.(string)
	}

	if pad := out.Get("padding"); pad != nil {
		out.padding = pad.(string)
	}

	if e := out.Get("ext"); e != nil {
		out.ext = e.(string)
	}

	if n := out.Get("num"); n != nil {
		out.num = n.(int)
	}

	if out.padding == "" {
		return fmt.Sprintf("%s%s", out.name, out.ext)
	}

	name := fmt.Sprintf("%s"+out.padding+"%s", out.name, out.num, out.ext)

	return name
}

func (out Output) Compile(s *ffmpeg.Stream) *ffmpeg.Stream {
	return s.Output(out.String(), out.KwArgs())
}

func (out *Output) Name(n string) *Output {
	out.name = n
	return out
}

func (out *Output) Pad(p string) *Output {
	out.padding = p
	return out
}

func (out *Output) Num(n int) *Output {
	out.num = n
	return out
}

func (out *Output) Ext(ext string) *Output {
	out.ext = ext
	return out
}

func (out *Output) Copy() *Output {
	out.Set("c", "copy")
	out.Del("c:a")
	out.Del("c:v")
	return out
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
