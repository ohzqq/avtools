package avtools

import "fmt"

type Media struct {
	files map[string]*FileFormat
	//Input         *FileFormat
	Cue      *FileFormat
	FFmeta   *FileFormat
	Json     *FileFormat
	Cover    *FileFormat
	CueChaps bool
	outMeta  *FileFormat
	json     []byte
}

func NewMedia(input string) *Media {
	media := Media{
		files: make(map[string]*FileFormat),
		//Input: NewFormat(input),
	}
	media.SetFile("input", input)
	//media.files["input"] = NewFormat(input)
	return &media
}

func (m Media) Meta() *MediaMeta {
	meta := m.GetFile("input").meta

	fmt.Printf("%+V\n", meta.Format.Tags)
	if m.HasFFmeta() {
		ff := m.GetFile("ffmeta")
		meta.SetChapters(ff.meta.Chapters)
		meta.SetTags(ff.meta.Format.Tags)
	}

	if m.HasCue() {
		meta.SetChapters(m.GetFile("cue").meta.Chapters)
	}

	return meta
}

func (m *Media) GetFile(file string) *FileFormat {
	return m.files[file]
}

func (m *Media) SetFile(name, f string) *Media {
	file := NewFormat(name)
	if f != "" {
		file.SetFile(f)
		if name != "cover" {
			file.Parse()
		}
	}
	m.files[name] = file
	return m
}

//func (m *Media) MarshalMetaTo(format string) *Media {
//  f := NewFormat(format)
//  m.outMeta = f.render(m.Meta())
//  return m
//  //f.SetMeta(f.Meta())
//}

//func (m Media) StringMeta() string {
//  if len(m.convertedData) > 0 {
//    return string(m.convertedData)
//  }
//  return ""
//}

//func (m Media) PrintMeta() {
//  println(m.StringMeta())
//}

//func (m Media) WriteMeta() {
//}

//func (m *Media) MarshalTo(format string) *Media {
//  switch format {
//  case "json", ".json":
//    m.data = MarshalJson(m)
//  case "ffmeta", "ini", ".ini":
//    m.data = RenderFFmetaTmpl(m)
//  case "cue", ".cue":
//    if !m.HasChapters() {
//      log.Fatal("No chapters")
//    }
//    m.data = RenderCueTmpl(m)
//  }
//  return m
//}

//func (m *Media) ConvertTo(kind string) *FileFormat {
//  f := m.GetFormat(kind)
//  f.render(f)

//  switch kind {
//  case "json", ".json":
//    f.data = MarshalJson(f)
//  case "ffmeta", "ini", ".ini":
//    f.data = RenderTmpl(f)
//  case "cue", ".cue":
//    if len(f.meta.Chapters) == 0 {
//      log.Fatal("No chapters")
//    }
//    f.data = RenderTmpl(f)
//        return fmt.Render("cue")
//  }
//  return f
//}

//func (m *Media) FFmetaChapsToCue() {
//  if !m.HasChapters() {
//    log.Fatal("No chapters")
//  }

//  f, err := os.Create("chapters.cue")
//  if err != nil {
//    log.Fatal(err)
//  }

//  tmpl, err := GetTmpl("cue")
//  if err != nil {
//    log.Println(err)
//  }

//  err = tmpl.Execute(f, m.Meta)
//  if err != nil {
//    log.Println("executing template:", err)
//  }
//}

//func (m *Media) SetChapters(ch Chapters) {
//  m.Meta.Chapters = ch
//}

func (m Media) HasCue() bool {
	return m.Cue != nil
}

func (m Media) HasCover() bool {
	return m.Cover != nil
}

func (m Media) HasFFmeta() bool {
	return m.FFmeta != nil
}

func (m Media) HasChapters() bool {
	if len(m.Meta().Chapters) != 0 {
		return true
	}
	return false
}

func (m Media) HasVideo() bool {
	for _, stream := range m.Meta().Streams {
		if stream.CodecType == "video" {
			return true
		}
	}
	return false
}

func (m Media) HasAudio() bool {
	for _, stream := range m.Meta().Streams {
		if stream.CodecType == "audio" {
			return true
		}
	}
	return false
}

func (m *Media) VideoCodec() string {
	for _, stream := range m.Meta().Streams {
		if stream.CodecType == "video" {
			return stream.CodecName
		}
	}
	return ""
}

func (m *Media) AudioCodec() string {
	for _, stream := range m.Meta().Streams {
		if stream.CodecType == "audio" {
			return stream.CodecName
		}
	}
	return ""
}

//func (m *Media) HasStreams() bool {
//  if m.HasMeta() && m.Meta.Streams != nil {
//    return true
//  }
//  return false
//}

//func (m *Media) HasFormat() bool {
//  if m.HasMeta() && m.Meta.Format != nil {
//    return true
//  }
//  return false
//}

func (m *Media) Duration() string {
	return secsToHHMMSS(m.Meta().Format.Duration)
}
