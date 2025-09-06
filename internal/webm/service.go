package webm

import (
	"context"
	"fmt"
	"io"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"github.com/yazmeyaa/telegram_sticker_converter/internal/converter"
)

type webmService struct{}

func NewService() *webmService {
	return &webmService{}
}

func (ws webmService) Transform(ctx context.Context, in io.Reader, out io.Writer, opts converter.WEBMTransformOptions) error {
	return ws.process(ctx, in, out, opts)
}

var (
	PresetMP4 = ffmpeg_go.KwArgs{
		"vcodec":   "libx264",
		"format":   "mp4",
		"pix_fmt":  "yuv420p",
		"movflags": "frag_keyframe+empty_moov",
		"preset":   "ultrafast",
		"tune":     "zerolatency",
	}

	PresetPNG = ffmpeg_go.KwArgs{
		"f":     "image2pipe",
		"c:v":   "png",
		"vsync": "0",
	}

	PresetJPEG = ffmpeg_go.KwArgs{
		"f":     "image2pipe",
		"c:v":   "mjpeg",
		"vsync": "0",
	}

	PresetWEBP = ffmpeg_go.KwArgs{
		"f":   "image2pipe",
		"c:v": "libwebp",
	}
	PresetGIF = ffmpeg_go.KwArgs{"f": "gif"}
)

func buildPreset(opts converter.WEBMTransformOptions) (ffmpeg_go.KwArgs, error) {
	switch opts.Format {
	case converter.FormatGIF:
		return PresetGIF, nil
	case converter.FormatMP4:
		return PresetMP4, nil
	case converter.FormatPNG:
		switch opts.Frame {
		case converter.FrameAll:
			return PresetPNG, nil
		case converter.FrameFirst:
			return ffmpeg_go.MergeKwArgs([]ffmpeg_go.KwArgs{
				PresetPNG,
				{
					"vframes": "1",
				},
			}), nil
		case converter.FrameN:
			return ffmpeg_go.MergeKwArgs([]ffmpeg_go.KwArgs{
				PresetPNG,
				{
					"vframes": "1",
					"vf":      fmt.Sprintf("select=eq(n\\,%d)", opts.FrameIndex),
				},
			}), nil
		case converter.FrameRange:
			return ffmpeg_go.MergeKwArgs([]ffmpeg_go.KwArgs{
				PresetPNG,
				{
					"vf": fmt.Sprintf("select=between(n\\,%d\\,%d)", opts.FrameIndex, opts.FrameIndex+opts.FrameOffset),
				},
			}), nil
		}
	case converter.FormatJPEG:
		switch opts.Frame {
		case converter.FrameAll:
			return PresetJPEG, nil
		case converter.FrameFirst:
			return ffmpeg_go.MergeKwArgs([]ffmpeg_go.KwArgs{
				PresetJPEG,
				{
					"frames": "1",
				},
			}), nil
		case converter.FrameN:
			return ffmpeg_go.MergeKwArgs([]ffmpeg_go.KwArgs{
				PresetJPEG,
				{
					"vframes": "1",
					"vf":      fmt.Sprintf("select=eq(n\\,%d)", opts.FrameIndex),
				},
			}), nil
		case converter.FrameRange:
			return ffmpeg_go.MergeKwArgs([]ffmpeg_go.KwArgs{
				PresetJPEG,
				{
					"vf": fmt.Sprintf("select=between(n\\,%d\\,%d)", opts.FrameIndex, opts.FrameIndex+opts.FrameOffset),
				},
			}), nil
		}
	case converter.FormatWEBP:
		switch opts.Frame {
		case converter.FrameAll:
			return PresetWEBP, nil
		case converter.FrameFirst:
			return ffmpeg_go.MergeKwArgs([]ffmpeg_go.KwArgs{
				PresetWEBP,
				{
					"frames": "1",
				},
			}), nil
		case converter.FrameN:
			return ffmpeg_go.MergeKwArgs([]ffmpeg_go.KwArgs{
				PresetWEBP,
				{
					"vframes": "1",
					"vf":      fmt.Sprintf("select=eq(n\\,%d)", opts.FrameIndex),
				},
			}), nil
		case converter.FrameRange:
			return ffmpeg_go.MergeKwArgs([]ffmpeg_go.KwArgs{
				PresetWEBP,
				{
					"vf": fmt.Sprintf("select=between(n\\,%d\\,%d)", opts.FrameIndex, opts.FrameIndex+opts.FrameOffset),
				},
			}), nil
		}
	}
	return ffmpeg_go.KwArgs{}, converter.ErrUnknownFormat
}

func (ws webmService) process(ctx context.Context, in io.Reader, out io.Writer, opts converter.WEBMTransformOptions) error {
	r, w := io.Pipe()

	go func() {
		defer w.Close()
		io.Copy(w, in)
	}()

	preset, err := buildPreset(opts)
	if err != nil {
		return err
	}

	err = ffmpeg_go.
		Input("pipe:0", ffmpeg_go.KwArgs{
			"f": "webm",
		}).
		Output("pipe:1", preset).
		WithInput(r).
		WithOutput(out).
		Run()

	if err != nil {
		return err
	}

	return nil
}
