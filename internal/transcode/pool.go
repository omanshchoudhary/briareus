package transcode

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
)

func All(files []string, outDir string, workers int) error {
	var wg sync.WaitGroup
	errs := make(chan error, len(files))
	jobs := make(chan string)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobs {
				output := filepath.Join(outDir, filepath.Base(file))
				if err := Segment(file, output); err != nil {
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

	var failed []error
	for err := range errs {
		failed = append(failed, err)
	}
	if len(failed) > 0 {
		return fmt.Errorf("%d of %d segments failed: %w", len(failed), len(files), errors.Join(failed...))
	}
	return nil
}
