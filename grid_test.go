package femtioelva_test

import (
	"image/color"
	"testing"

	"github.com/vikblom/femtioelva"
)

func sameColor(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

func TestBlank(t *testing.T) {
	box := femtioelva.GeoBox(0, 0, 100)
	grid := femtioelva.NewGrid(box, 4)
	img := grid.Draw(1, 1)

	rect := img.Bounds()
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			c := img.At(x, y)
			if !sameColor(c, color.White) {
				t.Fatalf("unexpected non-white pixel at %d %d", x, y)
			}
		}
	}
}

func TestSinglePos(t *testing.T) {
	box := femtioelva.GeoBox(femtioelva.GBG_LAT, femtioelva.GBG_LON, 100)
	grid := femtioelva.NewGrid(box, 10)

	east, north := femtioelva.LatLong2UTM(femtioelva.GBG_LAT, femtioelva.GBG_LON)
	grid.IncrUTM(east, north)

	got := 0      // Count how many black pixels
	cellsize := 2 // with a 2x2 cell there should be
	expected := 4

	img := grid.Draw(cellsize, 1)
	rect := img.Bounds()
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			c := img.At(x, y)
			if sameColor(c, color.Black) {
				got++
			} else if !sameColor(c, color.White) {
				t.Fatalf("color which is not black or white in image %v", c)
			}
		}
	}

	if got != expected {
		t.Errorf("expected %d black pixels but found %d", expected, got)
	}
}
