package main

import (
	"encoding/gob"
	"image"
	"image/png"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/vikblom/femtioelva"
)

func AddPositionsToGrid(grid femtioelva.Grid) {
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

	for _, v := range seen {
		// Skip boats
		if strings.Contains(v.Name, "Älv") {
			continue
		}
		east, north := femtioelva.LatLong2UTM(v.Lat, v.Long)
		grid.IncrUTM(east, north)
	}
}

func WriteImage(img image.Image, file string) {
	fh, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	err = png.Encode(fh, img)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	box := femtioelva.GeoBox(femtioelva.GBG_LAT, femtioelva.GBG_LON, 5_000)
	grid := femtioelva.NewGrid(box, 96)
	AddPositionsToGrid(grid)
	img := grid.Draw(8, 2)

	WriteImage(img, "img.png")
}
