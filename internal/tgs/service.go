package tgs

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/arugaz/go-rlottie"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"github.com/yazmeyaa/telegram_sticker_converter/internal/converter"
)

func (t tgsServiceImpl) Transform(ctx context.Context, in io.Reader, out io.Writer, opts converter.TGSTransformOptions) error {
	data, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer gr.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, gr); err != nil {
		return err
	}
	if opts.Format == converter.FormatLottie {
		_, err := out.Write(buf.Bytes())
		return err
	}

	anim := rlottie.LottieAnimationFromData(buf.String(), "", "")
	defer rlottie.LottieAnimationDestroy(anim)

	width, height := rlottie.LottieAnimationGetSize(anim)
	if opts.ResizeWidth == 0 && opts.ResizeHeight == 0 {
		opts.ResizeWidth = int(width)
		opts.ResizeHeight = int(height)
	}

	if opts.Format == converter.FormatGIF || opts.Format == converter.FormatWEBM || opts.Format == converter.FormatMP4 {
		return t.processVideo(ctx, anim, out, opts)
	}
	if opts.Format == converter.FormatPNG || opts.Format == converter.FormatWEBP || opts.Format == converter.FormatJPEG {
		return t.processFrames(ctx, anim, out, opts)
	}

	return converter.ErrUnknownFormat
}

func (t tgsServiceImpl) processFrames(ctx context.Context, anim rlottie.Lottie_Animation, out io.Writer, opts converter.TGSTransformOptions) error {
	width := uint(opts.ResizeWidth)
	height := uint(opts.ResizeHeight)

	if width == 0 || height == 0 {
		width, height = rlottie.LottieAnimationGetSize(anim)
	}

	frameBuffer := make([]byte, width*height*4)

	if opts.Frame == converter.FrameFirst {
		if err := t.processFrame(ctx, anim, 0, uint(opts.ResizeWidth), uint(opts.ResizeHeight), frameBuffer); err != nil {
			return err
		}

		return t.makeSingleImage(ctx, frameBuffer, out, opts)
	}

	if opts.Frame == converter.FrameN {
		if err := t.processFrame(ctx, anim, uint(opts.FrameIndex), uint(opts.ResizeWidth), uint(opts.ResizeHeight), frameBuffer); err != nil {
			return err
		}
		return t.makeSingleImage(ctx, frameBuffer, out, opts)
	}

	if opts.Frame == converter.FrameAll {
		return t.makeAllImages(ctx, anim, frameBuffer, out, opts)
	}

	return nil
}

var (
	PresetGif = ffmpeg_go.KwArgs{
		"f":       "gif",
		"pix_fmt": "rgb24",
	}
	PresetWebm = ffmpeg_go.KwArgs{
		"vcodec":   "libvpx",
		"format":   "webm",
		"pix_fmt":  "yuv420p",
		"b:v":      "0",
		"deadline": "realtime",
	}
	PresetMP4 = ffmpeg_go.KwArgs{
		"vcodec":   "libx264",
		"format":   "mp4",
		"pix_fmt":  "yuv420p",
		"movflags": "frag_keyframe+empty_moov",
		"preset":   "ultrafast",
		"tune":     "zerolatency",
	}
)

func (t tgsServiceImpl) processVideo(ctx context.Context, anim rlottie.Lottie_Animation, out io.Writer, opts converter.TGSTransformOptions) error {
	totalFrames := rlottie.LottieAnimationGetTotalframe(anim)
	frameRate := rlottie.LottieAnimationGetFramerate(anim)
	var preset ffmpeg_go.KwArgs

	switch opts.Format {
	case converter.FormatGIF:
		preset = PresetGif
	case converter.FormatWEBM:
		preset = PresetWebm
	case converter.FormatMP4:
		preset = PresetMP4
	default:
		return converter.ErrUnknownFormat
	}

	frameBuffer := make([]byte, opts.ResizeWidth*opts.ResizeHeight*4)
	r, w := io.Pipe()

	go func() {
		for frameIdx := range totalFrames {
			err := t.processFrame(ctx, anim, frameIdx, uint(opts.ResizeWidth), uint(opts.ResizeHeight), frameBuffer)
			if err != nil {
				w.CloseWithError(err)
				return
			}

			w.Write(frameBuffer)
		}
		w.Close()
	}()

	err := ffmpeg_go.
		Input("pipe:0", ffmpeg_go.KwArgs{
			"format":  "rawvideo",
			"pix_fmt": "bgra",
			"s":       fmt.Sprintf("%dx%d", opts.ResizeWidth, opts.ResizeHeight),
			"r":       frameRate,
		}).
		Silent(true).
		Output("pipe:1", preset).
		WithInput(r).
		WithOutput(out).
		Run()

	if err != nil {
		return err
	}

	return nil
}

func BGRAtoRGBA(buf []byte) {
	for i := 0; i < len(buf); i += 4 {
		buf[i], buf[i+2] = buf[i+2], buf[i]
	}
}

func (t tgsServiceImpl) makeSingleImage(ctx context.Context, frameBuffer []byte, out io.Writer, opts converter.TGSTransformOptions) error {
	BGRAtoRGBA(frameBuffer)
	img := &image.RGBA{
		Pix:    frameBuffer,
		Stride: int(opts.ResizeWidth) * 4,
		Rect:   image.Rect(0, 0, opts.ResizeWidth, opts.ResizeHeight),
	}

	if opts.Format == converter.FormatPNG {
		return png.Encode(out, img)
	}

	if opts.Format == converter.FormatJPEG {
		return jpeg.Encode(out, img, nil)
	}

	// TODO
	// if opts.Format == converter.FormatWEBP {
	// }
	return converter.ErrUnknownFormat
}

func (t tgsServiceImpl) makeAllImages(ctx context.Context, anim rlottie.Lottie_Animation, frameBuffer []byte, out io.Writer, opts converter.TGSTransformOptions) error {
	totalFrames := rlottie.LottieAnimationGetTotalframe(anim)
	archive := zip.NewWriter(out)
	defer archive.Close()

	for i := range totalFrames {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		file, err := archive.Create(fmt.Sprintf("frame_%d.%s", i, opts.Format))
		if err != nil {
			return err
		}

		t.processFrame(ctx, anim, i, uint(opts.ResizeWidth), uint(opts.ResizeHeight), frameBuffer)

		if err := t.makeSingleImage(ctx, frameBuffer, file, opts); err != nil {
			return err
		}
	}

	if err := archive.Close(); err != nil {
		return err
	}

	return nil
}

func (t tgsServiceImpl) processFrame(ctx context.Context, anim rlottie.Lottie_Animation, frameN uint, sx uint, sy uint, out []byte) error {
	expectedSize := int(sx * sy * 4)
	rlottie.LottieAnimationRender(anim, frameN, out, sx, sy, sx*4)
	if len(out) < expectedSize {
		return errors.New("not valid buffer size")
	}
	return nil
}

type tgsServiceImpl struct{}

func NewService() *tgsServiceImpl {
	return &tgsServiceImpl{}
}
