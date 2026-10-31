package transcode

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Segment(input, output string) error {
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	cmd := exec.Command("ffmpeg",
		"-v", "error",
		"-y",
		"-i", input,
		"-c:v", "libx264",
		"-preset", "medium",
		"-crf", "23",
		"-threads", "1",
		output,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("transcode %s: %w: %s", input, err, out)
	}
	return nil
}
