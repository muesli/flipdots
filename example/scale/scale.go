package main

import (
	"fmt"
	"github.com/muesli/flipdots"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"time"
)

var (
	ttfPath = "/Library/Fonts/Verdana.ttf"
	text    = "flipdots.go!"
)

func main() {
	fmt.Println("FLIPDOTS GO!")

	fd, err := flipdots.Init("flipdot.ffa:2323", 80, 16)
	if err != nil {
		panic(err)
	}

	img, _ := fd.TextToImage(text, ttfPath)
	for i := 16; i > 0; i-- {
		m := image.NewRGBA(image.Rect(0, 0, fd.Width, fd.Height))
		draw.Draw(m, m.Bounds(), fd.Background, image.ZP, draw.Src)

		scaleimg := resize.Resize(uint(fd.Width), uint(i), img, resize.Lanczos3)
		sr := image.Rectangle{image.Pt(0, 0), image.Pt(fd.Width, i)}

		mid := fd.Height / 2
		nmid := i / 2
		dp := image.Pt(0, mid-nmid)
		r := image.Rectangle{dp, dp.Add(sr.Size())}

		draw.Draw(m, r, scaleimg, sr.Min, draw.Src)
		fd.SendImage(m)
		time.Sleep(100)
	}

	fd.Clear()
}
