package ffmeta

import (
	"encoding/json"
	"log"

	"github.com/ohzqq/avtools/chap"
)

func LoadJson(d []byte) FFmeta {
	var meta FFmeta
	err := json.Unmarshal(d, &meta)
	if err != nil {
		log.Fatal(err)
	}

	if len(meta.Chaps) > 0 {
		for _, c := range meta.Chaps {
			ch := chap.NewChapter().SetMeta(c)
			meta.Chapters.Chapters = append(meta.Chapters.Chapters, ch)
		}
	}

	//meta.Tags = meta.Format.Tags

	return meta
}

func (ff FFmeta) DumpJson() []byte {
	data, err := json.Marshal(ff)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
