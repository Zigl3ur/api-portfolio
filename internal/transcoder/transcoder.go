package transcoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type VideoData struct {
	Path   string
	Name   string
	Rotate int
}

type FormatSpec struct {
	Width        int
	Height       int
	Bitrate      int
	Maxrate      int
	Bufsize      int
	AudioBitrate int
}

type Transcoder struct {
	destinationDir string
	Videos         []VideoData
	Formats        map[string]*FormatSpec
}

func NewTranscoder(destinationDir string, videos []VideoData, formats map[string]*FormatSpec) *Transcoder {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		panic(fmt.Sprintf("ffmpeg not found, %v", err))
	}

	if _, err := exec.LookPath("ffprobe"); err != nil {
		panic(fmt.Sprintf("ffprobe not found, %v", err))
	}

	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		panic(fmt.Sprintf("Error creating destination directory %s: %v", destinationDir, err))
	}

	return &Transcoder{
		destinationDir: destinationDir,
		Videos:         videos,
		Formats:        formats,
	}
}

func (t *Transcoder) GeneratePlaylistFile(path string, video VideoData) error {
	var content bytes.Buffer
	content.WriteString("#EXTM3U\n")

	for key, value := range t.Formats {
		var bitrate int
		var res string
		if key == "source" {
			width, height, br, err := GetVideoResolution(video.Path)
			if err != nil {
				return fmt.Errorf("error getting video resolution for source: %v", err)
			}
			bitrate = br / 1000
			res = fmt.Sprintf("%dx%d", width, height)
		} else {
			bitrate = value.Bitrate
			res = fmt.Sprintf("%dx%d", value.Width, value.Height)
		}

		fmt.Fprintf(&content, "#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s\n", bitrate, res)
		fmt.Fprintf(&content, "%s/playlist.m3u8\n", key)
	}

	return os.WriteFile(filepath.Join(t.destinationDir, video.Name, "playlist.m3u8"), content.Bytes(), 0644)
}

func rotateSide(angle int) string {
	switch angle {
	case 90:
		return "clock"
	case 180:
		return "clock,clock"
	case 270:
		return "cclock"
	default:
		return ""
	}
}

type probeResult struct {
	Streams []struct {
		Width   int `json:"width"`
		Height  int `json:"height"`
		BitRate int `json:"bit_rate,string"`
	} `json:"streams"`
}

func GetVideoResolution(path string) (width, height, bitrate int, err error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height,bit_rate",
		"-show_entries", "format=bit_rate",
		"-of", "json",
		path,
	)

	out, err := cmd.Output()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	var result probeResult
	if err := json.Unmarshal(out, &result); err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	if len(result.Streams) == 0 {
		return 0, 0, 0, fmt.Errorf("no video stream found in %s", path)
	}

	return result.Streams[0].Width, result.Streams[0].Height, result.Streams[0].BitRate, nil
}

func (t *Transcoder) BuildCmd(video VideoData, specs *FormatSpec, rotate int) *exec.Cmd {
	var vf strings.Builder

	rotation := rotateSide(rotate)
	if rotation != "" {
		fmt.Fprintf(&vf, "transpose_vaapi=dir=%s", rotation)
	}

	resLabel := "source"
	var audioArgs, videoArgs []string

	if specs != nil {
		if vf.Len() > 0 {
			vf.WriteString(",")
		}
		fmt.Fprintf(&vf, "scale_vaapi=w=%d:h=%d", specs.Width, specs.Height)
		resLabel = strconv.Itoa(specs.Height)

		audioArgs = []string{"-b:a", fmt.Sprintf("%dk", specs.AudioBitrate)}
		videoArgs = []string{
			"-b:v", fmt.Sprintf("%dk", specs.Bitrate),
			"-maxrate", fmt.Sprintf("%dk", specs.Maxrate),
			"-bufsize", fmt.Sprintf("%dk", specs.Bufsize),
		}
	}

	args := []string{
		"-hwaccel", "vaapi", "-hwaccel_device", "/dev/dri/renderD128", "-hwaccel_output_format", "vaapi",
		"-i", video.Path,
	}

	if vf.Len() > 0 {
		args = append(args, "-vf", vf.String())
	}

	args = append(args, "-c:a", "aac")
	args = append(args, audioArgs...)
	args = append(args, videoArgs...)
	args = append(args,
		"-c:v", "h264_vaapi", "-profile:v", "main",
		"-g", "48", "-keyint_min", "48",
		"-hls_time", "6", "-hls_playlist_type", "vod",
	)
	args = append(args,
		"-hls_segment_filename", fmt.Sprintf("%s/chunk_%%03d", resLabel),
		fmt.Sprintf("%s/playlist.m3u8", resLabel),
	)

	return exec.Command("ffmpeg", args...)
}

func (t *Transcoder) TranscodeAll() error {
	start := time.Now()

	for _, video := range t.Videos {
		videoDir := filepath.Join(t.destinationDir, video.Name)
		if err := os.MkdirAll(videoDir, 0755); err != nil {
			return fmt.Errorf("error creating video directory %s: %v", videoDir, err)
		}

		for res, specs := range t.Formats {
			outDir := filepath.Join(videoDir, res)
			fmt.Println("Transcoding video:", video.Name, "to resolution:", res+"p", "in directory:", outDir)
			if err := os.MkdirAll(outDir, 0755); err != nil {
				return fmt.Errorf("error creating output directory %s: %v", outDir, err)
			}

			cmd := t.BuildCmd(video, specs, video.Rotate)
			cmd.Dir = videoDir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("error transcoding video %s to resolution %s: %v", video.Name, res, err)
			}
		}

		if err := t.GeneratePlaylistFile(videoDir, video); err != nil {
			return fmt.Errorf("error generating playlist for video %s: %v", video.Name, err)
		}
	}

	fmt.Printf("Transcoding completed in %v\n", time.Since(start))

	return nil
}
