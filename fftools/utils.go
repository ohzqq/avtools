package fftools

import (
	"time"
	"strings"
	"log"
)

func secsToHHMMSS(sec string) string {
	d := strings.ReplaceAll(sec, ".", "s") + "us"
	dur, err := time.ParseDuration(d)
	if err != nil { log.Fatal(err) }
	duration := dur.String()
	h := strings.ReplaceAll(duration, "h", "")
	return duration
}
