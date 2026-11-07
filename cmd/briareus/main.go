package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sync"

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

	var wg sync.WaitGroup
	errs := make(chan error, len(files))
	jobs := make(chan string)

	workers := runtime.NumCPU()
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobs {
				output := filepath.Join("out/transcoded", filepath.Base(file))
				if err := transcode.Segment(file, output); err != nil {
					errs <- err
					continue
				}
				fmt.Println("done", output)
			}
		}()
	}

	for _, file := range files {
		jobs <- file
	}
	close(jobs)
	wg.Wait()
	close(errs)

	failed := 0
	for err := range errs {
		log.Println(err)
		failed++
	}
	if failed > 0 {
		log.Fatalf("%d of %d segments failed", failed, len(files))
	}
}
