package gif

import (
	"fmt"
	"versions/logs"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func DrawAll(path string, logs []logs.Log) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	pw := imagick.NewPixelWand()
	defer pw.Destroy()

	for i, l := range logs {
		img := fmt.Sprintf("%s/%d_%s", path, i, l.SHA1)
		in := img + ".png"
		check(mw.ReadImage(in))

		dw := imagick.NewDrawingWand()
		defer dw.Destroy()

		dw.SetFont("monospace")
		dw.SetFontSize(24)

		pw.SetColor("black")
		dw.SetStrokeColor(pw)
		dw.SetStrokeWidth(5)

		text := fmt.Sprintf("[%s]: %s", l.SHA1[:5], l.Title)
		dw.Annotation(25, 25+24/2, text)

		dw.SetStrokeWidth(0)

		pw.SetColor("white")
		dw.SetFillColor(pw)
		dw.SetTextAntialias(true)
		dw.Annotation(25, 25+24/2, text)

		check(mw.DrawImage(dw))

		out := img + ".gif"
		check(mw.WriteImage(out))
	}
}
