package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/yazmeyaa/telegram_sticker_converter/internal/converter"
	"github.com/yazmeyaa/telegram_sticker_converter/internal/tgs"
)

func main() {
	filePath := flag.String("file", "", "path to input .tgs file")
	outPath := flag.String("out", "", "path to output file")
	format := flag.String("format", "png", "output format (png|jpeg|webp|gif|webm|mp4|lottie)")
	frame := flag.String("frame", "all", "frame selector (first|all|n)")
	frameIndex := flag.Int("frame-index", 0, "frame index (used only with frame=n)")
	quality := flag.Int("quality", 90, "output quality (0-100)")
	resizeW := flag.Int("resize-width", 0, "resize width (0 = keep original)")
	resizeH := flag.Int("resize-height", 0, "resize height (0 = keep original)")
	flag.Parse()

	if *filePath == "" {
		fmt.Fprintln(os.Stderr, "error: -file is required")
		os.Exit(1)
	}

	if *outPath == "" {
		fmt.Fprintln(os.Stderr, "error: -out is required")
		os.Exit(1)
	}

	in, err := os.Open(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v\n", err)
		os.Exit(1)
	}
	defer in.Close()

	out, err := os.Create(*outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating output: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	var frameSel converter.FrameSelector
	switch *frame {
	case "first":
		frameSel = converter.FrameFirst
	case "all":
		frameSel = converter.FrameAll
	case "n":
		frameSel = converter.FrameN
	default:
		fmt.Fprintf(os.Stderr, "unknown frame selector: %s\n", *frame)
		os.Exit(1)
	}

	opts := converter.TransformOptions{
		Format:       converter.OutputFormat(*format),
		Frame:        frameSel,
		FrameIndex:   *frameIndex,
		Qualtity:     *quality,
		ResizeWidth:  *resizeW,
		ResizeHeight: *resizeH,
	}
	fmt.Fprintf(os.Stdout, "opts: %+v\n", opts)

	service := converter.TGSConverterService(tgs.NewService())
	if err := service.Transform(context.Background(), in, out, opts); err != nil {
		fmt.Fprintf(os.Stderr, "transform failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("success:", *outPath)
}
