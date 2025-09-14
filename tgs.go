package converter

import (
	"context"
	"io"
)

type TGSTransformOptions struct {
	Format       OutputFormat
	Frame        FrameSelector
	FrameIndex   int
	FrameOffset  int
	Qualtity     int
	ResizeWidth  int
	ResizeHeight int
	CacheKey     string
}

type TGSConverter interface {
	Transform(ctx context.Context, data io.Reader, out io.Writer, opts TGSTransformOptions) error
}
