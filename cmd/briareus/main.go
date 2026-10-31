package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/omanshchoudhary/briareus/internal/probe"
	"github.com/omanshchoudhary/briareus/internal/segment"
	"github.com/omanshchoudhary/briareus/internal/transcode"
)

func main() {
	keyframeTimes, err := probe.KeyframeTimes("testdata/test.mp4")

	if err != nil {
		log.Fatal(err)
	}
	cuts := segment.CutPoints(keyframeTimes, 10)
	files, err := segment.Split("testdata/test.mp4", cuts, "out/segments")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		output := filepath.Join("out/transcoded", filepath.Base(file))
		if err := transcode.Segment(file, output); err != nil {
			log.Fatal(err)
		}
		fmt.Println("done", output)
	}
}
