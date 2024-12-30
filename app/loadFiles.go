package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func ImgLoadInFolder(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isImage(info.Name()) {
			img, _, err := ebitenutil.NewImageFromFile(dir + "/" + info.Name())
			if err != nil {
				log.Fatal(err)
			}
			images = append(images, img)
		}
		return nil
	})
}

func isImage(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg"
}

func (g *Game) ImgScaleing(screen *ebiten.Image) float64 {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	imgW, imgH := images[g.idx].Bounds().Dx(), images[g.idx].Bounds().Dy()

	scaleX := float64(w) / float64(imgW)
	scaleY := float64(h) / float64(imgH)

	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}
	return scale
}
