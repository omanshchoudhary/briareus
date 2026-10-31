package probe

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type packet struct {
	PtsTime string `json:"pts_time"`
	Flags   string `json:"flags"`
}

type probeOutput struct {
	Packets []packet `json:"packets"`
}

func KeyframeTimes(path string) ([]float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "packet=pts_time,flags",
		"-of", "json",
		path,
	)

	out, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("ffprobe %s: %w", path, err)
	}
	var result probeOutput

	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("parse ffprobe output: %w", err)
	}

	var keyframeTimes []float64
	for _, p := range result.Packets {
		if !strings.Contains(p.Flags, "K") {
			continue
		}

		t, err := strconv.ParseFloat(p.PtsTime, 64)
		if err != nil {
			return nil, fmt.Errorf("parse pts_time %q: %w", p.PtsTime, err)
		}
		keyframeTimes = append(keyframeTimes, t)
	}

	return keyframeTimes, nil
}
