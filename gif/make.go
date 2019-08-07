package gif

import (
	"fmt"
	"image"
	"os"
	"versions/logs"

	"image/gif"
)

func Make(path string, logs []logs.Log) {
	out := &gif.GIF{}

	for i, l := range logs {
		img := fmt.Sprintf("%s/%d_%s.gif", path, i, l.SHA1)
		f, _ := os.Open(img)
		defer f.Close()

		in, _ := gif.Decode(f)
		out.Image = append(out.Image, in.(*image.Paletted))
		out.Delay = append(out.Delay, 100)
	}

	f, _ := os.OpenFile(fmt.Sprintf("%s/../out.gif", path), os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, out)

	// imagick.Initialize()
	// defer imagick.Terminate()

	// mw := imagick.NewMagickWand()

	// mw.SetImageFormat("gif")
	// coalesce := mw.CoalesceImages()
	// mw.Destroy()
	// defer coalesce.Destroy()

	// for i, l := range logs {
	// 	img := fmt.Sprintf("%s/%d_%s.png", path, i, l.SHA1)

	// 	mw := imagick.NewMagickWand()
	// 	check(mw.ReadImage(img))
	// 	defer mw.Destroy()

	// 	coalesce.AddImage(mw)
	// 	coalesce.SetImageDelay(100)
	// }

	// check(coalesce.WriteImage(fmt.Sprintf("%s/%s.gif", path, "versions")))
}
