package converter

import "errors"

type OutputFormat string
type FrameSelector int

const (
	FrameFirst FrameSelector = iota
	FrameAll
	FrameN
	FrameRange
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
