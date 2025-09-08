package webp

import (
	"context"
	"image/jpeg"
	"image/png"
	"io"

	converter "github.com/yazmeyaa/telegram_sticker_converter"
	"golang.org/x/image/webp"
)

type webpConverterService struct{}

// Transform implements converter.WebpConverterService.
func (w webpConverterService) Transform(ctx context.Context, in io.Reader, out io.Writer, opts converter.WebpTransformOptions) error {
	i, err := webp.Decode(in)
	if err != nil {
		return err
	}

	if opts.Format == converter.FormatPNG {
		return png.Encode(out, i)
	}
	if opts.Format == converter.FormatJPEG {
		return jpeg.Encode(out, i, nil)
	}

	return converter.ErrUnknownFormat
}

var _ converter.WebpConverterService = webpConverterService{}

func NewService() *webpConverterService {
	return &webpConverterService{}
}
