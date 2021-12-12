package femtioelva_test

import (
	"math"
	"testing"

	"github.com/vikblom/femtioelva"
)

func TestWhereIsGothenburg(t *testing.T) {
	// Sanity check against https://www.latlong.net/lat-long-utm.html
	east, north := femtioelva.LatLong2UTM(femtioelva.GBG_LAT, femtioelva.GBG_LON)
	if int(east) != 677214 {
		t.Error("UTM easting of Gothenburg in unexpected place.")
	}
	if int(north) != 6400188 {
		t.Error("UTM northing of Gothenburg in unexpected place.")
	}
}

func TestBoxIsSquare(t *testing.T) {
	side := 1000.0
	b := femtioelva.GeoBox(femtioelva.GBG_LAT, femtioelva.GBG_LON, side)

	e1, n1 := femtioelva.LatLong2UTM(b.LowLat, b.LowLong)
	e2, n2 := femtioelva.LatLong2UTM(b.HighLat, b.HighLong)

	width := e2 - e1
	height := n2 - n1

	diagonal := math.Sqrt(width*width + height*height)
	diff := diagonal - side*math.Sqrt(2)
	if diff < -1 || 1 < diff {
		t.Errorf("Expected diagonal length close to %f but got %f", side*math.Sqrt(2), diagonal)
	}

}
