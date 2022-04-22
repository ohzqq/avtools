package fftools

import (
	"fmt"
)

func (m *Media) Split() {
	cmd := NewCmd().In(m)
	ch := cmd.GetChapters()
	for i, chap := range *ch {
		m.Cut(chap.Start, chap.End, i)
	}
}

func (m *Media) Cut(ss, to string, no int) {
	count := fmt.Sprintf("%06d", no + 1)
	cmd := NewCmd().In(m)
	timestamps := make(map[string]string)
	if ss != "" {
		timestamps["ss"] = ss
	}
	if to != "" {
		timestamps["to"] = to
	}
	cmd.Args().PostInput(timestamps).Out("tmp" + count).Ext(m.Ext)
	cmd.Run()
}

