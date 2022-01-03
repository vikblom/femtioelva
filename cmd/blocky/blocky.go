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
	box := femtioelva.GeoBox(femtioelva.GBG_LAT, femtioelva.GBG_LON, 10_000)
	minEast, minNorth := femtioelva.LatLong2UTM(box.LowLat, box.LowLong)
	maxEast, maxNorth := femtioelva.LatLong2UTM(box.HighLat, box.HighLong)

	new := []pos{}
	for _, v := range old {
		if minEast <= v.east && v.east <= maxEast && minNorth <= v.north && v.north <= maxNorth {
			new = append(new, pos{v.east - minEast, v.north - minNorth})
		}
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

func DrawMatrix(m femtioelva.Matrix, bs, margin int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0,
		m.Width()*(bs+margin)+margin,
		m.Height()*(bs+margin)+margin))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)

	// Make the most popular spots completely dark
	scale := float64(m.Max()) / 255.0

	for row := 0; row < m.Height(); row++ {
		for col := 0; col < m.Width(); col++ {
			if m.At(row, col) == 0 {
				continue
			}
			x0 := margin + row*(margin+bs)
			y0 := margin + col*(margin+bs)
			r := image.Rect(x0, y0, x0+bs, y0+bs)
			val := 255 - uint8(float64(m.At(row, col))/scale) // Higher count -> darker
			draw.Draw(img, r, &image.Uniform{color.Gray{val}}, image.ZP, draw.Src)
		}
	}

	return img
}

func WriteImage(img *image.RGBA, file string) {
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

	ps := RetrievePositions()
	ps, width, height := Center(ps)
	lim := 0.0
	if width < height {
		lim = width
	} else {
		lim = height
	}
	m := PosMatrix(ps, 64+32, lim)
	img := DrawMatrix(m, 8, 2)

	WriteImage(img, "img.png")
}
