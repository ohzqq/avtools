package media

import (
	"github.com/ohzqq/avtools/ff"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ExtractCover(m *Media) {
	var stream Stream
	for _, s := range m.VideoStreams() {
		if s.IsCover {
			stream = s
		}
		break
	}
	cmd := ff.New()
	cmd.In(m.Input.Abs, ffmpeg.KwArgs{"y": ""})
	cmd.Output.Pad("").Set("c", "copy").Set("an", "")
	switch stream.CodecName {
	case "mjpeg":
		cmd.Ext(".jpg")
	case "png":
		cmd.Ext(".png")
	}
	cmd.Compile().Run()
}
