package main

import (
	"fmt"
	"github.com/muesli/flipdots"
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
	fd.SendImage(img)
}
