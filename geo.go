package femtioelva

import "github.com/wroge/wgs84"

const (
	GBG_LAT = 57.708870
	GBG_LON = 11.974560
)

var (
	// UTM for Sweden
	UTM = wgs84.UTM(32, true)
)

type Box struct {
	LowLat     float64
	LowLong    float64
	HighLat    float64
	HighLong   float64
	sideLength float64
}

func LatLong2UTM(lat, long float64) (east float64, north float64) {
	east, north, _ = wgs84.LonLat().To(UTM)(long, lat, 0)
	return
}

func UTM2LatLong(east, north float64) (lat float64, long float64) {
	long, lat, _ = UTM.To(wgs84.LonLat())(east, north, 0)
	return // flipped order
}

// GeoBox returns min/max long/lat for a square box with sides of width meters
// centered around long and lat.
func GeoBox(lat, long, width float64) Box {

	east, north := LatLong2UTM(lat, long)

	// bottom left
	lowLat, lowLong := UTM2LatLong(east-width/2, north-width/2)
	// top right
	highLat, highLong := UTM2LatLong(east+width/2, north+width/2)

	return Box{LowLat: lowLat, LowLong: lowLong, HighLat: highLat, HighLong: highLong, sideLength: width}
}
