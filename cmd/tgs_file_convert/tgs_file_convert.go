package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	converter "github.com/yazmeyaa/telegram_sticker_converter"
	"github.com/yazmeyaa/telegram_sticker_converter/tgs"
)

func main() {
	input := flag.String("input", "stream:stdin", "path to input .tgs file")
	outPath := flag.String("output", "", "path to output file")
	format := flag.String("format", "png", "output format (png|jpeg|webp|gif|webm|mp4|lottie)")
	frame := flag.String("frame", "all", "frame selector (first|all|n)")
	frameIndex := flag.Int("frame-index", 0, "frame index (used only with frame=n)")
	quality := flag.Int("quality", 100, "output quality (0-100)")
	resizeW := flag.Int("resize-width", 0, "resize width (0 = keep original)")
	resizeH := flag.Int("resize-height", 0, "resize height (0 = keep original)")
	flag.Parse()

	inputParts := strings.Split(*input, ":")
	if len(inputParts) < 2 {
		fmt.Fprintf(os.Stderr, "wrong input signature. valid is \"{type}:{input}\"\nExample: \"stream:stdin\"; \"file:input.tgs\"\n")
		os.Exit(1)
	}

	inputType := inputParts[0]
	var in *os.File
	if inputType == "file" {
		input, err := os.Open(inputParts[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening file: %v\n", err)
			os.Exit(1)
		}
		in = input
	}
	if inputType == "stream" {
		in = os.Stdin
	}
	if in == nil {
		fmt.Fprintf(os.Stderr, "Unexpected input type\n")
		os.Exit(1)
	}
	defer in.Close()

	if *outPath == "" {
		fmt.Fprintln(os.Stderr, "error: -output is required")
		os.Exit(1)
	}

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

	opts := converter.TGSTransformOptions{
		Format:       converter.OutputFormat(*format),
		Frame:        frameSel,
		FrameIndex:   *frameIndex,
		Qualtity:     *quality,
		ResizeWidth:  *resizeW,
		ResizeHeight: *resizeH,
	}

	converter := tgs.NewConverter()
	start := time.Now()
	if err := converter.Transform(context.Background(), in, out, opts); err != nil {
		fmt.Fprintf(os.Stderr, "transform failed: %v\n", err)
		os.Exit(1)
	}
	duration := time.Since(start)
	fmt.Printf("Complete in %d ms\n", duration.Milliseconds())

	fmt.Println("success:", *outPath)
}
