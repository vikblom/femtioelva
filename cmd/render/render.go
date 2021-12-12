package main

import (
	"encoding/gob"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/vikblom/femtioelva"
)

var (
	IMGWIDTH  = 1024
	IMGHEIGHT = 1024

	WIDTHSCALE  = (femtioelva.MAX_LONG - femtioelva.MIN_LONG) / IMGWIDTH
	HEIGHTSCALE = (femtioelva.MAX_LAT - femtioelva.MIN_LAT) / IMGHEIGHT
)

type pos struct {
	// x is from the left
	x int
	// y is from the top
	y int
}

type Path []pos

// RetrievePositions of one vehicle
func RetrievePositions() map[string]Path {
	file, err := os.Open("pos.gob")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	seen := []femtioelva.Vehicle{}
	err = gob.NewDecoder(file).Decode(&seen)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(len(seen))

	Paths := make(map[string]Path)
	for _, v := range seen {
		// Skip boats
		if strings.Contains(v.Name, "Älv") {
			continue
		}
		// scale and flip to img coordinates
		north := (v.Lat - femtioelva.MIN_LAT) / HEIGHTSCALE
		east := (v.Long - femtioelva.MIN_LONG) / WIDTHSCALE
		Paths[v.Gid] = append(Paths[v.Gid], pos{x: east, y: IMGHEIGHT - north})
	}

	return Paths
}

func WriteImage(paths map[string]Path) {
	img := image.NewRGBA(image.Rect(0, 0, IMGWIDTH, IMGHEIGHT))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)

	// TODO: Draw positions
	for _, ps := range paths {
		for _, p := range ps {
			r := image.Rect(p.x, p.y, p.x+1, p.y+1)
			draw.Draw(img, r, &image.Uniform{color.Black}, image.ZP, draw.Src)
		}
	}

	file, err := os.Create("img.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	ps := RetrievePositions()
	// TODO: Convert to map[id][]Pos or similar
	WriteImage(ps)

}
