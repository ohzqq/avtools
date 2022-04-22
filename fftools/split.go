package fftools

import (
	"fmt"
)
var _ = fmt.Printf

func Split(input string, profile string) {
	media := NewMedia(input)
	cmd := NewCmd().In(media).Profile(profile)
	ch := cmd.GetChapters()
	for i, chap := range *ch {
		media.Cut(chap.Start, chap.End, i)
	}
}

