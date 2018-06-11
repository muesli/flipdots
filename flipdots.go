package flipdots

import (
	"image"
	"image/draw"
	"io/ioutil"
	"math"
	"net"
	"time"

	"github.com/golang/freetype"
	colorful "github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/font"
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
	b := byte(0)
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
			cf, _ := colorful.MakeColor(color)
			l, _, _ := cf.Lab()
			if l >= 0.5 {
				imgmap = imgmap + "1"
			} else {
				imgmap = imgmap + "0"
			}
		}
	}

	return imgmap
}

func (fd *FlipDots) ScrollImage(img image.Image) error {
	b := img.Bounds()
	for x := 0; x < b.Max.X; x++ {
		subimg := image.NewRGBA(image.Rect(0, 0, fd.Width*16, fd.Height))
		draw.Draw(subimg, b, fd.Background, image.ZP, draw.Src)
		draw.Draw(subimg, b, img, b.Min.Add(image.Pt(x, 0)), draw.Src)

		err := fd.SendImage(subimg)
		if err != nil {
			return err
		}

		time.Sleep(55 * time.Millisecond)
	}

	return nil
}

func (fd *FlipDots) TextToImage(text, ttfPath string) (image.Image, error) {
	rgba := image.NewRGBA(image.Rect(0, 0, fd.Width*16, fd.Height))

	fontBytes, err := ioutil.ReadFile(ttfPath)
	if err != nil {
		return rgba, err
	}
	ttf, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return rgba, err
	}

	draw.Draw(rgba, rgba.Bounds(), fd.Background, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(fd.Dpi)
	c.SetFont(ttf)
	c.SetFontSize(fd.FontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fd.Foreground)

	// c.SetHinting(font.HintingNone)
	c.SetHinting(font.HintingFull)

	fh := c.PointToFixed(fd.FontSize/2.0) / 64
	yc := 1 + (float64(fd.Height) / 2.0) + (float64(fh) / 2.0)

	pt := freetype.Pt(1, int(yc))
	s, err := c.DrawString(text, pt)
	if err != nil {
		return rgba, err
	}
	// pt.Y += c.PointToFix32(fd.FontSize * fd.FontSpacing)

	return rgba.SubImage(image.Rect(0, 0, int(float64(s.X)/64), fd.Height)), nil
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
