package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/text/language"
)

type Game struct {
	idx     int
	Pressed bool
	images  []string
	img     *ebiten.Image
	trans   [][4]int
	sents   []string
}

//go:embed NanumGothic.ttf
var koreanTTF []byte
var koreanFaceSource *text.GoTextFaceSource

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(koreanTTF))
	if err != nil {
		log.Fatal(err)
	}
	koreanFaceSource = s
}

func (g *Game) Update() error {
	g.pageChange()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := imgScale(screen, g.img)
	screen.DrawImage(g.img, op)
	g.drawTextAndBorder(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	g := &Game{}
	g.loadFileNames("image")
	if len(g.images) == 0 {
		log.Fatal("0")
	}
	g.img, g.trans, g.sents = g.loadFiles("image")
	ebiten.SetWindowSize(640, 904)
	ebiten.SetWindowTitle("ImgViewer")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) pageChange() {
	if (ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)) && !g.Pressed {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.idx > 0 {
			g.idx--
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && g.idx < len(g.images)-1 {
			g.idx++
		}
		g.img, g.trans, g.sents = g.loadFiles("image")
		g.Pressed = true
	} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.Pressed = false
	}
}

func (g *Game) drawTextAndBorder(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()

	f := &text.GoTextFace{
		Source:    koreanFaceSource,
		Direction: text.DirectionLeftToRight,
		Size:      18,
		Language:  language.Korean,
	}

	for idx, val := range g.trans {
		y1, y2, x1, x2 := val[0], val[1], val[2], val[3]
		drawBorder(screen, y1, y2, x1, x2)
		if x >= x1 && x <= x2 && y >= y1 && y <= y2 {
			drawText(screen, x, y, g.sents[idx], f)
		}
	}
}

func drawBorder(screen *ebiten.Image, y1, y2, x1, x2 int) {
	lineWidth := 2.0
	color := color.Black
	vector.StrokeLine(screen, float32(x1), float32(y1), float32(x1), float32(y2), float32(lineWidth), color, false)
	vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y1), float32(lineWidth), color, false)
	vector.StrokeLine(screen, float32(x2), float32(y1), float32(x2), float32(y2), float32(lineWidth), color, false)
	vector.StrokeLine(screen, float32(x1), float32(y2), float32(x2), float32(y2), float32(lineWidth), color, false)
}

func drawText(screen *ebiten.Image, x, y int, t string, f *text.GoTextFace) {
	w, h := text.Measure(t, f, 0)
	gray := color.RGBA{0x80, 0x80, 0x80, 0xff}

	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), gray, false)

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	text.Draw(screen, t, f, op)
}

func (g *Game) loadFiles(dir string) (img *ebiten.Image, trans [][4]int, sents []string) {
	fn := g.images[g.idx]
	img, _, err := ebitenutil.NewImageFromFile(dir + "/" + fn)
	if err != nil {
		log.Fatal(err)
	}

	trans, sents = loadTexts(dir + "/" + fn[:len(fn)-4] + ".txt")
	return
}

func loadTexts(fn string) (result [][4]int, sents []string) {
	tf, err := os.Open(fn)
	if err != nil {
		return
	}
	defer tf.Close()

	scanner := bufio.NewScanner(tf)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Fields(line)
		if len(parts) < 5 {
			continue
		}

		var nums [4]int
		for i := 0; i < 4; i++ {
			_, err := fmt.Sscanf(parts[i], "%d", &nums[i])
			if err != nil {
				return
			}
		}

		s := strings.Join(parts[4:], " ")

		result = append(result, nums)
		sents = append(sents, s)
	}
	return
}

func (g *Game) loadFileNames(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isImage(info.Name()) {
			g.images = append(g.images, info.Name())
		}
		return nil
	})
}

func imgScale(screen *ebiten.Image, img *ebiten.Image) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}

	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()

	scaleX := float64(w) / float64(imgW)
	scaleY := float64(h) / float64(imgH)

	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}
	op.GeoM.Scale(scale, scale)

	return op
}

func isImage(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg"
}
