package avtools

type ProbeMeta struct {
	name        string
	Streams     []*StreamEntry `json:"streams"`
	FormatEntry `json:"format"`
	Chaps       []ChapterEntry `json:"chapters"`
}

type StreamEntry struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
}

type FormatEntry struct {
	Filename string            `json:"filename"`
	Dur      string            `json:"duration"`
	Size     string            `json:"size"`
	BitRate  string            `json:"bit_rate"`
	Tags     map[string]string `json:"tags"`
}

type ChapterEntry struct {
	Base         string            `json:"time_base",ini:"timebase"`
	StartTime    int               `json:"start",ini:"start"`
	EndTime      int               `json:"end",ini:"end"`
	ChapterTitle string            `json:"title", ini:"title"`
	Tags         map[string]string `json:"tags"`
}
