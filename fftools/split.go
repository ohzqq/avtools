package fftools

import (
	"fmt"
)

func (ff *FFmpegCmd) Split() {
	//var chapters Chapters
	if ff.args.CueSheet == "" {
	}
	fmt.Println("split")
}

func (ff *FFmpegCmd) getChapMeta() *Chapters {
	var (
		meta *MediaMeta
	)
	if ff.args.CueSheet != "" {
		meta = ReadCueSheet(ff.args.CueSheet)
	}
	return meta.Chapters
}
