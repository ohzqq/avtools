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

func (ff *FFmpegCmd) GetChapters() *Chapters {
	var (
		meta *MediaMeta
	)
	if ff.args.CueSheet != "" {
		meta = ReadCueSheet(ff.args.CueSheet)
	} else if ff.args.Metadata != "" {
		meta = ReadFFmetadata(ff.args.Metadata)
	} else {
		meta = ff.Meta()
	}
	return meta.Chapters
}
