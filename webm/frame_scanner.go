package webm

import (
	"bytes"
	"errors"
	"io"

	converter "github.com/yazmeyaa/telegram_sticker_converter"
)

type frameScanner struct {
	r    io.Reader
	opts converter.WEBMTransformOptions

	buf   []byte
	frame []byte
	eof   bool
}

func newScanner(r io.Reader, opts converter.WEBMTransformOptions) *frameScanner {
	return &frameScanner{
		r:    r,
		opts: opts,
		buf:  make([]byte, 0, 64*1024),
	}
}

var (
	pngStart = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	pngEnd   = []byte{0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}

	jpegStart = []byte{0xFF, 0xD8}
	jpegEnd   = []byte{0xFF, 0xD9}

	webpStart = []byte{'R', 'I', 'F', 'F'}
	webpWebP  = []byte{'W', 'E', 'B', 'P'}
)

func (f *frameScanner) signatures() ([]byte, []byte) {
	switch f.opts.Format {
	case converter.FormatPNG:
		return pngStart, pngEnd
	case converter.FormatJPEG:
		return jpegStart, jpegEnd
	case converter.FormatWEBP:
		return webpStart, nil
	default:
		return nil, nil
	}
}

func (f *frameScanner) Next() ([]byte, error) {
	if f.eof {
		return nil, io.EOF
	}

	startSig, endSig := f.signatures()
	if startSig == nil {
		return nil, converter.ErrUnknownFormat
	}

	for {
		startIdx := bytes.Index(f.buf, startSig)
		if startIdx >= 0 {
			switch f.opts.Format {
			case converter.FormatWEBP:
				if len(f.buf[startIdx:]) < 12 {
					break
				}
				if !bytes.Equal(f.buf[startIdx+8:startIdx+12], webpWebP) {
					f.buf = f.buf[startIdx+4:]
					continue
				}
				size := int(uint32(f.buf[startIdx+4]) |
					uint32(f.buf[startIdx+5])<<8 |
					uint32(f.buf[startIdx+6])<<16 |
					uint32(f.buf[startIdx+7])<<24)
				total := 8 + size
				if len(f.buf[startIdx:]) < total {
					break
				}
				frame := f.buf[startIdx : startIdx+total]
				f.buf = append([]byte{}, f.buf[startIdx+total:]...)
				return frame, nil

			default:
				endIdx := bytes.Index(f.buf[startIdx+len(startSig):], endSig)
				if endIdx >= 0 {
					endIdx += startIdx + len(startSig)
					frame := f.buf[startIdx : endIdx+len(endSig)]
					f.buf = append([]byte{}, f.buf[endIdx+len(endSig):]...)
					return frame, nil
				}
			}
		}

		tmp := make([]byte, 8192)
		n, err := f.r.Read(tmp)
		if n > 0 {
			f.buf = append(f.buf, tmp[:n]...)
			continue
		}
		if errors.Is(err, io.EOF) {
			f.eof = true
			return nil, io.EOF
		}
		if err != nil {
			return nil, err
		}
	}
}
