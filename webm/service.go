package webm

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	converter "github.com/yazmeyaa/telegram_sticker_converter"
)

type webmConverter struct{}

func NewConverter() *webmConverter {
	return &webmConverter{}
}

func (ws webmConverter) Transform(ctx context.Context, in io.Reader, out io.Writer, opts converter.WEBMTransformOptions) error {
	rIn, wIn := io.Pipe()
	rOut, wOut := io.Pipe()

	go func() {
		defer wIn.Close()
		_, _ = io.Copy(wIn, in)
	}()

	preset, err := buildPreset(opts)
	if err != nil {
		return err
	}

	inputVideo := ffmpeg_go.
		Input("pipe:0").
		Silent(true)

	var output *ffmpeg_go.Stream

	switch opts.Format {
	case converter.FormatMP4:
		inputAudio := ffmpeg_go.Input(
			"anullsrc=channel_layout=stereo:sample_rate=44100",
			ffmpeg_go.KwArgs{"f": "lavfi"},
		)

		output = ffmpeg_go.Output(
			[]*ffmpeg_go.Stream{inputVideo, inputAudio},
			"pipe:1",
			preset,
		)

	default:
		output = inputVideo.Output("pipe:1", preset)
	}

	output = output.Silent(true).WithInput(rIn).WithOutput(wOut)

	if (opts.Frame == converter.FrameAll || opts.Frame == converter.FrameRange) &&
		(opts.Format == converter.FormatWEBP || opts.Format == converter.FormatJPEG || opts.Format == converter.FormatPNG) {

		output = output.Silent(true).WithInput(rIn).WithOutput(wOut)

		go func() {
			defer wOut.Close()
			_ = output.Run()
		}()

		zw := zip.NewWriter(out)
		defer zw.Close()

		frameIdx := 0
		scanner := newScanner(rOut, opts)

		for {
			frame, err := scanner.Next()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return err
			}

			f, err := zw.Create(fmt.Sprintf("frame_%d.%s", frameIdx, string(opts.Format)))
			if err != nil {
				return err
			}

			if _, err := f.Write(frame); err != nil {
				return err
			}
			frameIdx++
		}

		return nil
	}

	if err := output.WithInput(rIn).WithOutput(out).Run(); err != nil {
		return err
	}

	return nil
}
