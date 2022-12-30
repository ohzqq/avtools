package avtools

import (
	"encoding/json"
	"log"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var probeArgs = []ffmpeg.KwArgs{
	ffmpeg.KwArgs{"show_chapters": ""},
	ffmpeg.KwArgs{"select_streams": "a"},
	ffmpeg.KwArgs{"show_entries": "stream=codec_type,codec_name:format=filename, start_time, duration, size, bit_rate:format_tags"},
	ffmpeg.KwArgs{"of": "json"},
}

func Probe(input string) []byte {
	args := ffmpeg.MergeKwArgs(probeArgs)
	info, err := ffmpeg.ProbeWithTimeoutExec(input, 0, args)
	if err != nil {
		log.Fatal(err)
	}

	return []byte(info)
}

func ReadEmbeddedMeta(input string) *ProbeMeta {
	info := Probe(input)

	var meta ProbeMeta
	err := json.Unmarshal(info, &meta)
	if err != nil {
		log.Fatal(err)
	}

	return &meta
}
