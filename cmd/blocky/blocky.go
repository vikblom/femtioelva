package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/vikblom/femtioelva"
)

func AddPositionsToGrid(gobfile string, grid femtioelva.Grid) {
	file, err := os.Open(gobfile)
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
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Println("usage: blocky [path to .gob] [path to .png]")
		os.Exit(1)
	}
	gobfile := flag.Args()[0]
	pngfile := flag.Args()[1]

	box := femtioelva.GeoBox(femtioelva.GBG_LAT, femtioelva.GBG_LON, 5_000)
	grid := femtioelva.NewGrid(box, 96)
	AddPositionsToGrid(gobfile, grid)
	img := grid.Draw(8, 2)

	WriteImage(img, pngfile)
}
