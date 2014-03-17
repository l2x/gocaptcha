package gocaptcha

import (
	"bufio"
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"code.google.com/p/draw2d/draw2d"
	"code.google.com/p/freetype-go/freetype"
)

var (
	text = []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9",
		"A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "M", "N", "p", "Q", "R", "S", "T", "W", "X", "Y", "Z",
	}
)

type Capthca struct {
	Dpi      float64
	Font     string
	FontSize float64
	Width    int
	Height   int
	Length   int
}

func New() *Capthca {
	return &Capthca{
		Dpi:      72,
		Font:     "./luxisr.ttf",
		FontSize: 20,
		Width:    80,
		Height:   30,
		Length:   4,
	}
}

func (c *Capthca) Config(conf map[string]interface{}) {
	if v, ok := conf["Dpi"]; ok {
		c.Dpi = v.(float64)
	}
	if v, ok := conf["Font"]; ok {
		c.Font = v.(string)
	}
	if v, ok := conf["FontSize"]; ok {
		c.FontSize = v.(float64)
	}
	if v, ok := conf["Width"]; ok {
		c.Width = v.(int)
	}
	if v, ok := conf["Height"]; ok {
		c.Height = v.(int)
	}
	if v, ok := conf["Length"]; ok {
		c.Length = v.(int)
	}
}

func (c *Capthca) Create() (*bytes.Buffer, string, error) {
	f := new(bytes.Buffer)
	txt := c.getText()

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(c.Font)
	if err != nil {
		return f, txt, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return f, txt, err
	}

	// Initialize the context.
	fg, bg := image.Black, image.White

	rgba := image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	ft := freetype.NewContext()
	ft.SetDPI(c.Dpi)
	ft.SetFont(font)
	ft.SetFontSize(c.FontSize)
	ft.SetClip(rgba.Bounds())
	ft.SetDst(rgba)
	ft.SetSrc(fg)
	ft.SetHinting(freetype.NoHinting)

	// Draw the text.
	x := 5
	_y := int(ft.PointToFix32(c.FontSize) >> 8)
	y := (c.Height-_y)/2 + _y

	pt := freetype.Pt(x, y)
	_, err = ft.DrawString(txt, pt)
	if err != nil {
		return f, txt, err
	}

	i := random(1, 4)
	for j := 0; j < i; j++ {

		ts := float64(random(0, c.Width))
		te := float64(random(0, c.Height))
		es := float64(random(0, c.Width))
		ee := float64(random(0, c.Height))
		//draw line
		gc := draw2d.NewGraphicContext(rgba)
		gc.SetLineWidth(2)
		gc.MoveTo(ts, te)
		gc.LineTo(es, ee)
		gc.Stroke()
	}

	b := bufio.NewWriter(f)
	err = png.Encode(b, rgba)
	if err != nil {
		return f, txt, err
	}
	err = b.Flush()
	if err != nil {
		return f, txt, err
	}

	return f, txt, nil
}

func (c *Capthca) getText() string {
	txt := []string{}
	l := len(text)
	for i := 0; i < c.Length; i++ {
		k := random(0, l)
		txt = append(txt, text[k])
	}
	return strings.Join(txt, " ")
}

func random(min, max int) int {
	rand.Seed(int64(time.Now().Nanosecond()))
	return rand.Intn(max-min) + min
}
