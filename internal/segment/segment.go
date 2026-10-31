package segment

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// keyframeTimes and returned cuts are keyframe pts in seconds
func CutPoints(keyframeTimes []float64, target float64) []float64 {
	if len(keyframeTimes) == 0 {
		return nil
	}

	cuts := []float64{keyframeTimes[0]}
	for _, time := range keyframeTimes {
		last := cuts[len(cuts)-1]
		if time-last >= target {
			cuts = append(cuts, time)
		}
	}
	return cuts
}

func Split(path string, cuts []float64, outDir string) ([]string, error) {
	if len(cuts) < 2 {
		return nil, fmt.Errorf("need at least 2 cut points, got %d", len(cuts))
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, fmt.Errorf("create %s: %w", outDir, err)
	}
	times := make([]string, 0, len(cuts)-1)
	for _, c := range cuts[1:] {
		times = append(times, strconv.FormatFloat(c, 'f', -1, 64))
	}
	cmd := exec.Command("ffmpeg",
		"-v", "error",
		"-y",
		"-i", path,
		"-map", "0:v:0",
		"-c", "copy",
		"-f", "segment",
		"-segment_times", strings.Join(times, ","),
		"-segment_time_delta", "0.001",
		"-reset_timestamps", "1",
		filepath.Join(outDir, "seg_%03d.mp4"),
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg split: %w: %s", err, out)
	}

	return filepath.Glob(filepath.Join(outDir, "seg_*.mp4"))
}
