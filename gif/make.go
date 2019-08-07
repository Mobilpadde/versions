package gif

import (
	"fmt"
	"image"
	"os"
	"versions/logs"

	"image/gif"
)

func Make(dumps, outPath string, logs []logs.Log) {
	out := &gif.GIF{}

	for i, l := range logs {
		img := fmt.Sprintf("%s/%d_%s.gif", dumps, i, l.SHA1)
		f, _ := os.Open(img)
		defer f.Close()

		in, _ := gif.Decode(f)
		out.Image = append(out.Image, in.(*image.Paletted))
		out.Delay = append(out.Delay, 100)
	}

	f, _ := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, out)
}
