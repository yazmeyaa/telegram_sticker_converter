package webm

import (
	"archive/zip"
	"context"
	"errors"
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

	stream := ffmpeg_go.
		Input("pipe:0", ffmpeg_go.KwArgs{
			"f": "webm",
		}).
		Silent(true).
		Output("pipe:1", preset).
		WithInput(rIn).
		WithOutput(wOut)

	if (opts.Frame == converter.FrameAll || opts.Frame == converter.FrameRange) &&
		(opts.Format == converter.FormatWEBP || opts.Format == converter.FormatJPEG || opts.Format == converter.FormatPNG) {

		go func() {
			defer wOut.Close()
			_ = stream.Run()
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

	defer wOut.Close()
	if err := stream.WithOutput(out).Run(); err != nil {
		return err
	}

	return nil
}
