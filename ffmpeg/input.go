package ffmpeg

import "strconv"

type Input struct {
	files       []string
	inputMap    []string
	FFmetadata  string
	HasChapters bool
}

func (i *Input) Add(input string, m ...string) *Input {
	total := len(i.files)
	i.files = append(i.files, input)
	switch mi := len(m); mi {
	case 0:
		i.inputMap = append(i.inputMap, strconv.Itoa(total))
	default:
		i.inputMap = append(i.inputMap, m[0])
	}
	return i
}

func (i Input) Parse() []string {
	var in []string
	var inMap []string
	total := len(i.files)
	for idx, input := range i.files {
		in = append(in, "-i", input)
		if total > 1 {
			inMap = append(inMap, "-map", i.inputMap[idx])
		}
	}

	if i.HasChapters {
		inMap = append(inMap, "-map_chapters", strconv.Itoa(total))
	}

	if i.FFmetadata != "" {
		in = append(in, "-i", i.FFmetadata)
		inMap = append(inMap, "-map_metadata", strconv.Itoa(total))
	}

	var input []string
	input = append(input, in...)
	input = append(input, inMap...)
	return input
}
