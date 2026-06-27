package video

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type res struct {
	label        string
	width        int
	height       int
	bitrate      int
	maxrate      int
	bufsize      int
	audioBitrate int
}

var formats = []res{
	{label: "1080p", width: 1920, height: 1080, bitrate: 5000, maxrate: 5350, bufsize: 7500, audioBitrate: 128},
	{label: "720p", width: 1280, height: 720, bitrate: 2800, maxrate: 2996, bufsize: 4200, audioBitrate: 128},
	{label: "480p", width: 854, height: 480, bitrate: 1400, maxrate: 1498, bufsize: 2100, audioBitrate: 96},
	{label: "240p", width: 426, height: 240, bitrate: 400, maxrate: 428, bufsize: 642, audioBitrate: 64},
}

func GeneratePlaylist(outputPath string) error {
	var content bytes.Buffer
	content.WriteString("#EXTM3U\n")

	for _, format := range formats {
		fmt.Fprintf(&content, "#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d\n", format.bitrate*1000, format.width, format.height)
		fmt.Fprintf(&content, "%s.m3u8\n", format.label)
	}

	return os.WriteFile(outputPath, content.Bytes(), 0644)
}

func TranscodeVideo(inputFile, outputDir string, res res) error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found: %w", err)
	}

	cmd := exec.Command("ffmpeg", "-i", inputFile, `-vf`, fmt.Sprintf("scale=w=%d:h=%d", res.width, res.height),
		"-c:a", "aac", "-b:a", fmt.Sprintf("%dk", res.audioBitrate),
		"-c:v", "h264", "-profile:v", "main", "-crf", "20", "-sc_threshold", "0", "-g", "48", "-keyint_min", "48", "-hls_time", "6", "-hls_playlist_type", "vod",
		"-b:v", fmt.Sprintf("%dk", res.bitrate), "-maxrate", fmt.Sprintf("%dk", res.maxrate), "-bufsize", fmt.Sprintf("%dk", res.bufsize),
		"-hls_segment_filename", fmt.Sprintf("%s/%s_%%03d.ts", outputDir, res.label),
		fmt.Sprintf("%s/%s.m3u8", outputDir, res.label),
	)

	return cmd.Run()
}
