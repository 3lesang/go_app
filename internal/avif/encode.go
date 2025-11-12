package avif

import (
	"bytes"
	"image"

	_ "image/jpeg"
	_ "image/png"

	"github.com/gen2brain/avif"
)

func EncodeImageToAVIF(imgBytes []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	options := avif.Options{
		Quality: 10,
		Speed:   4,
	}
	if err := avif.Encode(buf, img, options); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
