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
		// {Path: "/home/eden/Videos/portfolio_vids/fzwt_2026/PXL_20260616_205907809.mp4", Name: "lp_1", Rotate: 90},
		// {Path: "/home/eden/Videos/portfolio_vids/fzwt_2026/PXL_20260616_191949415.mp4", Name: "lp_2"},
		// {Path: "/home/eden/Videos/portfolio_vids/fzwt_2026/PXL_20260616_204918862.mp4", Name: "lp_3", Rotate: 90},
		// {Path: "/home/eden/Videos/portfolio_vids/fzwt_2026/PXL_20260616_205510062.mp4", Name: "lp_4", Rotate: 90},
		// {Path: "/home/eden/Videos/portfolio_vids/fzwt_2026/PXL_20260616_210820168.mp4", Name: "lp_5", Rotate: 90},
		// {Path: "/home/eden/Videos/portfolio_vids/fzwt_2026/PXL_20260616_211407090.mp4", Name: "lp_6", Rotate: 90},
		// {Path: "/home/eden/Videos/portfolio_vids/fzwt_2026/PXL_20260616_194156019.mp4", Name: "lp_7"},
		// {Path: "/home/eden/Videos/portfolio_vids/fzwt_2026/PXL_20260616_195455845.mp4", Name: "lp_8", Rotate: 90},

		// {Path: "/home/eden/Videos/portfolio_vids/mayhem_ball_2025/IMG_4752.MP4", Name: "gaga_1", Rotate: 90},
		// {Path: "/home/eden/Videos/portfolio_vids/mayhem_ball_2025/IMG_4755.MP4", Name: "gaga_2", Rotate: 90},
		// {Path: "/home/eden/Videos/portfolio_vids/mayhem_ball_2025/IMG_4758.MP4", Name: "gaga_3"},
		// {Path: "/home/eden/Videos/portfolio_vids/mayhem_ball_2025/IMG_4787.MP4", Name: "gaga_4", Rotate: 90},
		// {Path: "/home/eden/Videos/portfolio_vids/mayhem_ball_2025/IMG_4992.MP4", Name: "gaga_5"},

		// {Path: "/home/eden/Videos/portfolio_vids/until_youre_statisfied_tour/PXL_20260609_192018811.mp4", Name: "yonaka_1", Rotate: 90},
		// {Path: "/home/eden/Videos/portfolio_vids/until_youre_statisfied_tour/PXL_20260609_193102097.mp4", Name: "yonaka_2", Rotate: 90},
		// {Path: "/home/eden/Videos/portfolio_vids/until_youre_statisfied_tour/PXL_20260609_193952590.mp4", Name: "yonaka_3", Rotate: 90},

		// {Path: "/home/eden/Videos/portfolio_vids/manic_dream_tour/IMG_5129.MOV", Name: "deathbyromy_1", Rotate: 90},
		// {Path: "/home/eden/Videos/portfolio_vids/manic_dream_tour/IMG_5132.MOV", Name: "deathbyromy_2", Rotate: 90},
		// {Path: "/home/eden/Videos/portfolio_vids/manic_dream_tour/IMG_5133.MOV", Name: "deathbyromy_3", Rotate: 90},
		// {Path: "/home/eden/Videos/portfolio_vids/manic_dream_tour/IMG_5135.MOV", Name: "deathbyromy_4", Rotate: 90},

		{Path: "/home/eden/Videos/portfolio_vids/fzwt_2024/IMG_4121.MOV", Name: "lp_1_2024", Rotate: 90},
		{Path: "/home/eden/Videos/portfolio_vids/fzwt_2024/IMG_4122.MOV", Name: "lp_2_2024", Rotate: 90},
	}

	t := transcoder.NewTranscoder(outDir, videos, formats)

	if err := t.TranscodeAll(); err != nil {
		fmt.Printf("Error transcoding videos: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
