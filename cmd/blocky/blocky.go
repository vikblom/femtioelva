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
	east  float64
	north float64
}

// RetrievePositions of one vehicle
func RetrievePositions() []pos {
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

	positions := []pos{}
	for _, v := range seen {

		// Skip boats
		if strings.Contains(v.Name, "Älv") {
			continue
		}

		east, north := femtioelva.LatLong2UTM(float64(v.Lat)/1_000_000, float64(v.Long)/1_000_000)
		positions = append(positions, pos{east, north})
	}

	return positions
}

func Center(old []pos) ([]pos, float64, float64) {

	minEast, minNorth := femtioelva.LatLong2UTM(femtioelva.BOX.LowLat, femtioelva.BOX.LowLong)
	maxEast, maxNorth := femtioelva.LatLong2UTM(femtioelva.BOX.HighLat, femtioelva.BOX.HighLong)

	new := []pos{}
	for _, v := range old {
		new = append(new, pos{v.east - minEast, v.north - minNorth})
	}
	return new, maxEast - minEast, maxNorth - minNorth
}

func PosMatrix(ps []pos, n int, max float64) femtioelva.Matrix {
	m := femtioelva.NewMatrix(n, n)

	d := max / float64(n) // size of each cell
	for _, p := range ps {
		if p.east > max || p.north > max {
			continue
		}
		row := int((max - p.north) / d)
		col := int(p.east / d)
		m.Incr(row, col)
	}

	return m
}

func WriteMatrix(m femtioelva.Matrix, bs, margin int) {
	img := image.NewRGBA(image.Rect(0, 0,
		m.Width()*(bs+margin)+margin,
		m.Height()*(bs+margin)+margin))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	file, err := os.Create("img.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for row := 0; row < m.Height(); row++ {
		for col := 0; col < m.Width(); col++ {
			if m.At(row, col) == 0 {
				continue
			}
			x0 := margin + row*(margin+bs)
			y0 := margin + col*(margin+bs)
			r := image.Rect(x0, y0, x0+bs, y0+bs)
			draw.Draw(img, r, &image.Uniform{color.Black}, image.ZP, draw.Src)
		}
	}

	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	ps := RetrievePositions()
	ps, width, height := Center(ps)
	lim := 0.0
	if width < height {
		lim = width
	} else {
		lim = height
	}
	m := PosMatrix(ps, 64, lim)
	WriteMatrix(m, 8, 2)
}
