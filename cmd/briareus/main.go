package main

import (
	"log"
	"runtime"

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

	if err := transcode.All(files, "out/transcoded", runtime.NumCPU()); err != nil {
		log.Fatal(err)
	}
}
