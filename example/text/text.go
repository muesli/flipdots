package main

import (
	"github.com/muesli/flipdots"
)

var (
	ttfPath = "/usr/share/fonts/TTF/Roboto-Medium.ttf"
	text    = "Fuddelwuddelduddelmoo!"
)

func main() {
	fd, err := flipdots.New("flipdot.lab:2323", 80, 16)
	if err != nil {
		panic(err)
	}

	img, _ := fd.TextToImage(text, ttfPath)
	fd.ScrollImage(img)
}
