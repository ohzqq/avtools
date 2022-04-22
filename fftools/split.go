package fftools

import (
	"fmt"
)
var _ = fmt.Printf

func (cmd *FFmpegCmd) Split() {
	ch := cmd.GetChapters()
	for i, chap := range *ch {
		cmd.Input[0].Cut(chap.Start, chap.End, i)
	}
}

