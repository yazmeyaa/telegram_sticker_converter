package converter

import (
	"context"
	"io"
)

type WebpTransformOptions struct {
	Format OutputFormat
}

type WebpConverter interface {
	Transform(ctx context.Context, in io.Reader, out io.Writer, opts WebpTransformOptions) error
}
