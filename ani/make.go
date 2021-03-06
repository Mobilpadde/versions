package ani

import (
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/Mobilpadde/versions/logs"

	"github.com/sizeofint/webpanimation"
)

func Make(dumps, outPath string, logs []logs.Log) {
	spd := 1000

	ll := len(logs)
	i := ll - 1
	l := logs[i]

	path := fmt.Sprintf("%s/%d_%s.gif", dumps, i, l.SHA1)
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		log.Println(err)
	}

	b := img.Bounds().Max
	webpanim := webpanimation.NewWebpAnimation(b.X, b.Y, 0)
	webpanim.WebPAnimEncoderOptions.SetKmin(9)
	webpanim.WebPAnimEncoderOptions.SetKmax(17)
	defer webpanim.ReleaseMemory()
	webpConfig := webpanimation.NewWebpConfig()
	webpConfig.SetLossless(1)

	err = webpanim.AddFrame(img, spd*(ll-(i+1)), webpConfig)
	if err != nil {
		log.Println(err)
	}

	for i := len(logs) - 2; i >= 0; i-- {
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

		err = webpanim.AddFrame(img, spd*(ll-(i+1)), webpConfig)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	var buf bytes.Buffer
	err = webpanim.Encode(&buf)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(outPath, buf.Bytes(), 0777)
}
