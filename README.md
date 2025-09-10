# Telegram Sticker Converter

Converter for Telegram stickers and animations into various formats.\
Supports **TGS**, **WEBP**, **WEBM**, with the ability to extract single
frames, frame sequences, or convert to video.

## Features

-   Support for main Telegram sticker formats: `TGS`, `WEBP`, `WEBM`
-   Conversion to static images, animations, and videos
-   Flexible frame extraction: first, single, range, or all frames
-   Resizing on transformation
-   Easy to use as a Go library

## Requirements

-   **Go** ≥ 1.24\
-   **ffmpeg** ≥ 7.1.1 (for video/animations)

## Supported conversions

| Input    | Output                                                                                                       |
| -------- | ------------------------------------------------------------------------------------------------------------ |
| **TGS**  | PNG (first, all, N-th, range) <br> JPEG (same) <br> WEBP (same) <br> Lottie JSON <br> GIF <br> WEBM <br> MP4 |
| **WEBP** | PNG <br> JPEG                                                                                                |
| **WEBM** | MP4 <br> GIF <br> PNG (first, all, N-th, range) <br> JPEG (same) <br> WEBP (same)                            |

## Installation

``` bash
go get github.com/yazmeyaa/telegram_sticker_converter
```

## Usage examples

> General pattern: open input file, create output file, pass
> `TransformOptions`.

### Extract frame from TGS

``` go
conv := tgs.NewConverter()

r, _ := os.Open("./sticker.tgs")
defer r.Close()

w, _ := os.Create("./sticker.png")
defer w.Close()

opts := converter.TGSTransformOptions{
    Format:     converter.FormatPNG,
    Frame:      converter.FrameN,
    FrameIndex: 10,
    ResizeWidth:  1024,
    ResizeHeight: 1024,
}

_ = conv.Transform(context.Background(), r, w, opts)
```

### Convert TGS → MP4

``` go
opts := converter.TGSTransformOptions{
    Format:      converter.FormatMP4,
    ResizeWidth: 1024,
    ResizeHeight: 1024,
}
_ = conv.Transform(ctx, r, w, opts)
```

### Convert TGS → ZIP with frames

``` go
opts := converter.TGSTransformOptions{
    Format:      converter.FormatPNG,
    Frame:       converter.FrameAll,
    ResizeWidth: 1024,
    ResizeHeight: 1024,
}
_ = conv.Transform(ctx, r, w, opts)
```

## License

GPLv3