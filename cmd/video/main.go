package main

import (
	"fmt"
	"os"

	"github.com/Zigl3ur/api-portfolio/internal/transcoder"
)

var formats = map[string]transcoder.FormatSpec{
	"source": {Bitrate: 12000, Maxrate: 16000, Bufsize: 24000, AudioBitrate: 192},
	"1080":   {Width: 1920, Height: 1080, Bitrate: 5000, Maxrate: 5350, Bufsize: 7500, AudioBitrate: 128},
	"720":    {Width: 1280, Height: 720, Bitrate: 2800, Maxrate: 2996, Bufsize: 4200, AudioBitrate: 128},
	"480":    {Width: 854, Height: 480, Bitrate: 1400, Maxrate: 1498, Bufsize: 2100, AudioBitrate: 96},
	"240":    {Width: 426, Height: 240, Bitrate: 400, Maxrate: 428, Bufsize: 642, AudioBitrate: 64},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Provide the output dir as the first arg please")
		os.Exit(1)
	}

	outDir := os.Args[1]

	videos := []transcoder.VideoData{
		{Path: "/home/eden/Downloads/PXL_20260616_205907809.mp4", Name: "lp_1", Rotate: 90},
		{Path: "/home/eden/Downloads/IMG_4992.MP4", Name: "gaga_1"},
	}

	t := transcoder.NewTranscoder(outDir, videos, formats)

	if err := t.TranscodeAll(); err != nil {
		fmt.Printf("Error transcoding videos: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
