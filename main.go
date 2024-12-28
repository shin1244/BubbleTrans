package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	idx     int
	Pressed bool
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !g.Pressed {
		if g.idx > 0 {
			g.idx--
		}
		g.Pressed = true
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && !g.Pressed {
		if g.idx < len(Images)-1 {
			g.idx++
		}
		g.Pressed = true
	}
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.Pressed = false
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	scale := g.imgScaleing(screen)
	op.GeoM.Scale(scale, scale)

	screen.DrawImage(Images[g.idx], op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func ImageLoadInFolder(dir string) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isImage(info.Name()) {
			Img, _, err = ebitenutil.NewImageFromFile(dir + "/" + info.Name())
			if err != nil {
				log.Fatal(err)
			}
			Images = append(Images, Img)
		}
		return nil
	})
}

func isImage(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg"
}

func (g *Game) imgScaleing(screen *ebiten.Image) float64 {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	imgW, imgH := Images[g.idx].Bounds().Dx(), Images[g.idx].Bounds().Dy()

	scaleX := float64(w) / float64(imgW)
	scaleY := float64(h) / float64(imgH)

	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}
	return scale
}

var Images []*ebiten.Image
var Img *ebiten.Image
var err error

func main() {
	ImageLoadInFolder("image")
	if len(Images) == 0 {
		log.Fatal("0")
	}
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("ImgViewer")

	if err = ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
