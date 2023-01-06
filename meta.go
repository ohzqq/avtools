package avtools

import (
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FFmeta struct {
	Filename string            `json:"filename"`
	Dur      string            `json:"duration"`
	Size     string            `json:"size"`
	BitRate  string            `json:"bit_rate"`
	Tags     map[string]string `json:"tags"`
	Streams  []*Stream
	Chapters []*Chapter
}

type Stream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
	stream    *ffmpeg.Stream
}

type chapterEntry struct {
	Base  string            `json:"time_base",ini:"timebase"`
	Start int               `json:"start",ini:"start"`
	End   int               `json:"end",ini:"end"`
	Title string            `json:"title", ini:"title"`
	Tags  map[string]string `json:"tags"`
}
