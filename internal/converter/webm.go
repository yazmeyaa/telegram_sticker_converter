package converter

import (
	"context"
	"io"
)

type WEBMTransformOptions struct {
	Format      OutputFormat
	Width       uint
	Height      uint
	Frame       FrameSelector
	FrameIndex  uint
	FrameOffset uint
}

type WEBMConverter interface {
	Transform(ctx context.Context, in io.Reader, out io.Writer, opts WEBMTransformOptions) error
}
