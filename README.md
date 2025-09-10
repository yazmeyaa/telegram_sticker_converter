# Telegram sticker converter

## Required system binaries
- ffmpeg 7.1.1

## Build requirements
- Go 1.24

## Supported formats
- TGS
- WEBP
- WEBM

## Supported transformations

| Input    | Output formats                                                                                                                                                                                                       |
| -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **TGS**  | PNG (first frame, all frames, N frame, frame range) <br> JPEG (first frame, all frames, N frame, frame range) <br> WEBP (first frame, all frames, N frame, frame range) <br> Lottie JSON <br> GIF <br> WEBM <br> MP4 |
| **WEBP** | PNG <br> JPEG                                                                                                                                                                                                        |
| **WEBM** | MP4 <br> GIF <br> PNG (first frame, all frames, N frame, frame range) <br> JPEG (first frame, all frames, N frame, frame range) <br> WEBP (first frame, all frames, N frame, frame range)                            |


## Examples
### Extract frame from TGS
```go
package main
import (
	"context"
	"os"

	converter "github.com/yazmeyaa/telegram_sticker_converter"
	"github.com/yazmeyaa/telegram_sticker_converter/tgs"
)


func main() {
	conv := tgs.NewConverter()

	r, err := os.Open("./sticker.tgs")
	defer r.Close()
	w, err := os.Create("./sticker.png")
	defer w.Close()
	opts := converter.TGSTransformOptions{
		Format:       converter.FormatPNG,
		Frame:        converter.FrameN,
		FrameIndex:   10,
		ResizeWidth:  1024,
		ResizeHeight: 1024,
	}
	if err := conv.Transform(context.Background(), r, w, opts); err != nil {
		panic(err)
	}
}
```

### Convert TGS to video
```go
package main
import (
	"context"
	"os"

	converter "github.com/yazmeyaa/telegram_sticker_converter"
	"github.com/yazmeyaa/telegram_sticker_converter/tgs"
)


func main() {
	conv := tgs.NewConverter()

	r, err := os.Open("./sticker.tgs")
	defer r.Close()
	w, err := os.Create("./sticker.mp4")
	defer w.Close()
	opts := converter.TGSTransformOptions{
		Format:       converter.FormatMP4,
		ResizeWidth:  1024,
		ResizeHeight: 1024,
	}
	if err := conv.Transform(context.Background(), r, w, opts); err != nil {
		panic(err)
	}
}
```

### Convert TGS to frames array (ZIP)
```go
package main
import (
	"context"
	"os"

	converter "github.com/yazmeyaa/telegram_sticker_converter"
	"github.com/yazmeyaa/telegram_sticker_converter/tgs"
)


func main() {
	conv := tgs.NewConverter()

	r, err := os.Open("./sticker.tgs")
	defer r.Close()
	w, err := os.Create("./sticker.zip")
	defer w.Close()
	opts := converter.TGSTransformOptions{
		Format:       converter.FormatPNG,
        Frame:        converter.FrameAll,
		ResizeWidth:  1024,
		ResizeHeight: 1024,
	}
	if err := conv.Transform(context.Background(), r, w, opts); err != nil {
		panic(err)
	}
}
```