package ani

import (
	"fmt"
	"log"

	"github.com/Mobilpadde/versions/logs"

	"github.com/fogleman/gg"
)

func DrawAll(path, uri string, logs []logs.Log) {
	for i, l := range logs {
		path := fmt.Sprintf("%s/%d_%s", path, i, l.SHA1)
		in := path + ".png"
		img, err := gg.LoadPNG(in)
		if err != nil {
			log.Println(err)
			continue
		}

		sz := img.Bounds().Size()
		dc := gg.NewContext(sz.X, sz.Y)
		dc.DrawImage(img, 0, 0)

		if err := dc.LoadFontFace("./fonts/Oswald-Bold.ttf", 32); err != nil {
			panic(err)
		}

		pos := 25.0
		dc.SetRGB(0, 0, 0)

		uriPath := ""
		if uri != "" {
			uriPath = fmt.Sprintf(" (%s)", uri)
		}
		s := fmt.Sprintf("[%s]: %s%s", l.SHA1, l.Title, uriPath)
		n := 6
		for dy := -n; dy <= n; dy++ {
			for dx := -n; dx <= n; dx++ {
				if dx*dx+dy*dy >= n*n {
					// give it rounded corners
					continue
				}
				x := pos + float64(dx)
				y := pos + float64(dy)
				dc.DrawStringAnchored(s, x, y, 0, 1)
			}
		}

		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(s, pos, pos, 0, 1)

		dc.SavePNG(path + "_annotated.png")
	}
}
