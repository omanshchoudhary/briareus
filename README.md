# Briareus

Segment-parallel distributed video transcoder in Go, built on FFmpeg.

Briareus is under active development.

## Overview

Briareus divides a video into independently encodable segments at keyframe
boundaries. FFmpeg performs the encoding, while Briareus handles probing,
segmentation and job execution.

| Stage | Description |
| --- | --- |
| Probe | Reads keyframe timestamps from the source with `ffprobe`. |
| Segment | Selects keyframe-aligned cut points and splits the source with stream copy. |
| Transcode | Encodes each segment with `libx264`. |

## Requirements

- Go 1.27 or newer
- FFmpeg and ffprobe on your `PATH`, with `libx264` support

## Usage

Place the source video at `testdata/test.mp4`, then run the following from the
project root:

```sh
go run ./cmd/briareus
```

| Path | Contents |
| --- | --- |
| `out/segments/` | Source segments |
| `out/transcoded/` | Transcoded segments |

## Known Limitations

### Open-GOP sources

Briareus splits input at keyframes using stream copy and assumes **closed GOPs**,
where no frame references a frame outside its own GOP.

Some sources, commonly broadcast and Blu-ray content, use **open GOPs**, where
B-frames at the end of one GOP reference the I-frame of the next. Cutting at that
I-frame separates those B-frames from their reference, which can cause visual
artifacts at segment boundaries.

Content encoded with x264's default settings produces closed GOPs and is
unaffected.

### Output

Output is written as individual video-only segments. Source audio is not included.

### Input duration

Inputs must be long enough to produce at least two segments at the 10-second
target duration.
