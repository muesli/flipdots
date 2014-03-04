package flipdots

import (
	"code.google.com/p/freetype-go/freetype"
	"image"
	"image/draw"
	"io/ioutil"
	"math"
	"net"
)

type FlipDots struct {
	Addr        *net.UDPAddr
	Connection  *net.UDPConn
	Width       int
	Height      int
	Dpi         float64
	FontSize    float64
	FontSpacing float64
	Foreground  *image.Uniform
	Background  *image.Uniform
}

func listToByte(s string) byte {
	var b byte
	b = 0
	for i := 0; i < 8; i++ {
		if s[i] == '1' {
			b += byte(math.Pow(float64(2), float64(7-i)))
		}
	}

	return b
}

func matrixToPacket(s string) []byte {
	b := []byte{}

	for i := 0; i < len(s)/8; i++ {
		b = append(b, listToByte(s[i*8:i*8+8]))
	}

	return b
}

func (fd *FlipDots) ImageToMatrix(img image.Image) string {
	imgmap := ""
	for row := 0; row < fd.Height; row++ {
		for column := 0; column < fd.Width; column++ {
			color := img.At(column, row)
			pr, pg, pb, _ := color.RGBA()
			if pr > 32767 || pg > 32767 || pb > 32767 {
				imgmap = imgmap + "1"
			} else {
				imgmap = imgmap + "0"
			}
		}
	}

	return imgmap
}

func (fd *FlipDots) TextToImage(text, ttfPath string) (image.Image, error) {
	rgba := image.NewRGBA(image.Rect(0, 0, fd.Width, fd.Height))

	fontBytes, err := ioutil.ReadFile(ttfPath)
	if err != nil {
		return rgba, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return rgba, err
	}

	draw.Draw(rgba, rgba.Bounds(), fd.Background, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(fd.Dpi)
	c.SetFont(font)
	c.SetFontSize(fd.FontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fd.Foreground)

	//c.SetHinting(freetype.NoHinting)
	c.SetHinting(freetype.FullHinting)

	pt := freetype.Pt(1, 1+int(c.PointToFix32(fd.FontSize)>>8))
	_, err = c.DrawString(text, pt)
	if err != nil {
		return rgba, err
	}
	//	pt.Y += c.PointToFix32(fd.FontSize * fd.FontSpacing)

	return rgba, nil
}

func (fd *FlipDots) Clear() error {
	i := image.NewRGBA(image.Rect(0, 0, fd.Width, fd.Height))
	draw.Draw(i, i.Bounds(), fd.Background, image.ZP, draw.Src)
	return fd.SendImage(i)
}

func (fd *FlipDots) SendImage(img image.Image) error {
	imgmap := fd.ImageToMatrix(img)
	_, err := fd.Connection.Write(matrixToPacket(imgmap))
	return err
}

func Init(addr string, width int, height int) (FlipDots, error) {
	fd := FlipDots{
		Width:       width,
		Height:      height,
		Dpi:         72.0,
		FontSize:    12.0,
		FontSpacing: 1.1,
		Foreground:  image.Black,
		Background:  image.White,
	}

	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fd, err
	}
	fd.Addr = serverAddr

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return fd, err
	}
	fd.Connection = conn
	return fd, nil
}
