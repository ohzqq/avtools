package ffmpeg

import "strconv"

type Input struct {
	input []string
	Meta  string
	Cover string
}

func (i Input) Map() []string {
	total := len(i.input)

	var input []string
	for _, in := range i.input {
		input = append(input, "-i", in)
	}

	if i.Cover != "" {
		input = append(input, "-i", i.Cover)
	}

	if i.Meta != "" {
		input = append(input, "-i", i.Meta)
	}

	if i.Cover != "" || i.Meta != "" {
		for idx, _ := range i.input {
			input = append(input, "-map", strconv.Itoa(idx)+":0")
		}
	}

	if i.Cover != "" {
		input = append(input, "-map", strconv.Itoa(total)+":0")
		total++
	}

	if i.Meta != "" {
		input = append(input, "-map_metadata", strconv.Itoa(total)+":0")
		total++
	}

	return input
}
