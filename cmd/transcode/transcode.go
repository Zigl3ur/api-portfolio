//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type res struct {
	width        int
	height       int
	bitrate      int
	maxrate      int
	bufsize      int
	audioBitrate int
}

var formats = map[string]res{
	"1080p": {width: 1920, height: 1080, bitrate: 5000, maxrate: 5350, bufsize: 7500, audioBitrate: 128},
	"720p":  {width: 1280, height: 720, bitrate: 2800, maxrate: 2996, bufsize: 4200, audioBitrate: 128},
	"480p":  {width: 854, height: 480, bitrate: 1400, maxrate: 1498, bufsize: 2100, audioBitrate: 96},
	"240p":  {width: 426, height: 240, bitrate: 400, maxrate: 428, bufsize: 642, audioBitrate: 64},
}

func generatePlaylist(outputPath string, videoName string) error {
	var content bytes.Buffer
	content.WriteString("#EXTM3U\n")

	for key, specs := range formats {
		fmt.Fprintf(&content, "#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d\n", specs.bitrate*1000, specs.width, specs.height)
		vidResName := fmt.Sprintf("%s-%s.m3u8", videoName, key)
		fmt.Fprintf(&content, "%s/%s.m3u8\n", vidResName, vidResName)
	}

	return os.WriteFile(filepath.Join(outputPath, "playlist.m3u8"), content.Bytes(), 0644)
}

func transcodeVideo(inputFile, outputDir, name, res string) (*exec.Cmd, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return nil, fmt.Errorf("ffmpeg not found in PATH: %v", err)
	}

	specs, ok := formats[res]
	if !ok {
		return nil, fmt.Errorf("unsupported resolution: %s", res)
	}

	cmd := exec.Command("ffmpeg",
		"-hwaccel", "vaapi", "-hwaccel_device", "/dev/dri/renderD128", "-hwaccel_output_format", "vaapi",
		"-i", inputFile,
		"-vf", fmt.Sprintf("scale_vaapi=w=%d:h=%d", specs.width, specs.height),
		"-c:a", "aac", "-b:a", fmt.Sprintf("%dk", specs.audioBitrate),
		"-c:v", "h264_vaapi", "-profile:v", "main",
		"-g", "48", "-keyint_min", "48",
		"-hls_time", "6", "-hls_playlist_type", "vod",
		"-b:v", fmt.Sprintf("%dk", specs.bitrate),
		"-maxrate", fmt.Sprintf("%dk", specs.maxrate),
		"-bufsize", fmt.Sprintf("%dk", specs.bufsize),
		"-hls_segment_filename", fmt.Sprintf("%s/%s-%s_%%03d.ts", outputDir, name, res),
		fmt.Sprintf("%s/%s-%s.m3u8", outputDir, name, res),
	)

	return cmd, nil
}

type videos struct {
	path string
	name string
}

func transcodeVideos(vids []videos, outputDir string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(vids)*len(formats))

	for _, video := range vids {
		vidDir := filepath.Join(outputDir, video.name)
		if err := os.MkdirAll(vidDir, 0755); err != nil {
			return fmt.Errorf("Error creating video directory %s: %v", vidDir, err)
		}

		for res := range formats {
			outDir := filepath.Join(vidDir, fmt.Sprintf("%s-%s", video.name, res))
			fmt.Println("Transcoding video:", video.name, "to resolution:", res, "in directory:", outDir)
			if err := os.MkdirAll(outDir, 0755); err != nil {
				return fmt.Errorf("Error creating output directory %s: %v", outDir, err)
			}

			wg.Go(func() {
				cmd, err := transcodeVideo(video.path, outDir, video.name, res)
				if err != nil {
					errChan <- fmt.Errorf("Error transcoding video %s to %s: %v with path %s", video.name, res, err, video.path)
					return
				}

				if err := cmd.Run(); err != nil {
					errChan <- fmt.Errorf("Error running ffmpeg for video %s to %s: %v with path %s", video.name, res, err, video.path)
					fmt.Println(err)
					return
				}
			})
		}
		generatePlaylist(vidDir, video.name)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		var errMsg strings.Builder
		for err := range errChan {
			errMsg.WriteString(err.Error())
			errMsg.WriteString("\n")
		}
		return fmt.Errorf("Errors occurred during transcoding:\n%s", errMsg.String())
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Provide the output dir as the first arg please")
		os.Exit(1)
	}

	outDir := os.Args[1]

	videos := []videos{
		{path: "/home/eden/Downloads/PXL_20260616_194105268.mp4", name: "lp_1"},
	}

	if err := transcodeVideos(videos, outDir); err != nil {
		fmt.Printf("Error transcoding videos: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
