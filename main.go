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
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !g.Pressed {
		if g.idx > 0 {
			g.idx--
		}
		fmt.Println(g.trans)
		g.img, g.trans, g.sents = g.loadFiles("image")
		g.Pressed = true
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && !g.Pressed {
		if g.idx < len(g.images)-1 {
			g.idx++
		}
		g.img, g.trans, g.sents = g.loadFiles("image")
		fmt.Println(g.trans)
		g.Pressed = true
	}
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.Pressed = false
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	scale := g.imgScale(screen, g.img)
	op.GeoM.Scale(scale, scale)

	screen.DrawImage(g.img, op)

	g.drawTextWithOrigin(screen)
}

func (g *Game) drawTextWithOrigin(screen *ebiten.Image) {
	mouseX, mouseY := ebiten.CursorPosition()
	gray := color.RGBA{0x80, 0x80, 0x80, 0xff}

	f := &text.GoTextFace{
		Source:    koreanFaceSource,
		Direction: text.DirectionLeftToRight,
		Size:      18,
		Language:  language.Korean,
	}

	for idx, val := range g.trans {
		x1, x2, y1, y2 := val[0], val[1], val[2], val[3]
		if mouseX >= x1 && mouseX <= x2 && mouseY >= y1 && mouseY <= y2 {
			x, y := mouseX, mouseY
			w, h := text.Measure(g.sents[idx], f, 0)

			vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), gray, false)

			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			text.Draw(screen, g.sents[idx], f, op)
		}
	}

}

func (g *Game) loadFiles(dir string) (img *ebiten.Image, trans [][4]int, sents []string) {
	fn := g.images[g.idx]
	img, _, err := ebitenutil.NewImageFromFile(dir + "/" + fn)
	if err != nil {
		log.Fatal(err)
	}

	trans, sents = txtScanner(dir + "/" + fn[:len(fn)-4] + ".txt")
	return
}

func txtScanner(fn string) (result [][4]int, sents []string) {
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
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

func (g *Game) imgScale(screen *ebiten.Image, img *ebiten.Image) float64 {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()

	scaleX := float64(w) / float64(imgW)
	scaleY := float64(h) / float64(imgH)

	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}
	return scale
}

func isImage(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg"
}

func main() {
	g := &Game{}
	g.loadFileNames("image")
	if len(g.images) == 0 {
		log.Fatal("0")
	}
	g.img, g.trans, g.sents = g.loadFiles("image")
	fmt.Println(g.trans)
	ebiten.SetWindowSize(595, 841)
	ebiten.SetWindowTitle("ImgViewer")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
