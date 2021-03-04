package gif

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"

	"versions/logs"

	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
)

func Make(dumps, outPath string, logs []logs.Log) {
	out := &gif.GIF{}

	for i := len(logs) - 1; i >= 0; i-- {
		l := logs[i]

		path := fmt.Sprintf("%s/%d_%s.gif", dumps, i, l.SHA1)
		f, err := os.Open(path)
		if err != nil {
			log.Println(err)
			continue
		}

		defer f.Close()

		img, err := png.Decode(f)
		if err != nil {
			log.Println(err)
			continue
		}

		pm := image.NewPaletted(img.Bounds(), palette.WebSafe)
		draw.Draw(pm, img.Bounds(), img, image.ZP, draw.Over)

		var buf bytes.Buffer
		err = gif.Encode(&buf, pm, nil)

		out.Image = append(out.Image, pm)
		out.Delay = append(out.Delay, 100)
	}

	f, _ := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()

	err := gif.EncodeAll(f, out)
	if err != nil {
		log.Println(err)
	}
}
