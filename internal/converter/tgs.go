package converter

import (
	"context"
	"errors"
	"io"
)

type OutputFormat string
type FrameSelector int

const (
	FrameFirst FrameSelector = iota
	FrameAll
	FrameN
)

var (
	ErrUnknownFormat error = errors.New("unknown format")
)

const (
	FormatPNG    OutputFormat = "png"
	FormatJPEG   OutputFormat = "jpeg"
	FormatWEBP   OutputFormat = "webp"
	FormatGIF    OutputFormat = "gif"
	FormatWEBM   OutputFormat = "webm"
	FormatMP4    OutputFormat = "mp4"
	FormatLottie OutputFormat = "lottie"
)

type TransformOptions struct {
	Format       OutputFormat
	Frame        FrameSelector
	FrameIndex   int
	Qualtity     int
	ResizeWidth  int
	ResizeHeight int
}

type TGSConverterService interface {
	Transform(ctx context.Context, data io.Reader, out io.Writer, opts TransformOptions) error
}
