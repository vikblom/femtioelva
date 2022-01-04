package main

import (
	"flag"
	"image/png"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/vikblom/femtioelva"
)

// Server state
// TODO: Mutex this?
var (
	box  = femtioelva.GeoBox(femtioelva.GBG_LAT, femtioelva.GBG_LON, 10_000)
	grid = femtioelva.NewGrid(box, 96)
)

func serveGrid(w http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		log.Errorf("serveGrid got a %s request", request.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	img := grid.Draw(8, 2) // TODO: Move graphic options to Grid
	err := png.Encode(w, img)
	if err != nil {
		log.Errorf("encoding png failed: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func serveAssets(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets")
}

func main() {
	verboseFlag := flag.Bool("v", false, "verbose logging")
	flag.Parse()

	if *verboseFlag {
		log.SetLevel(log.DebugLevel)
		log.Debug("Verbose prints enabled")
	}

	apikey := os.Getenv("VASTTRAFIKAPI")
	if apikey == "" {
		log.Fatal("Could not read API key from env: VASTTRAFIKAPI")
		os.Exit(1)
	}

	port := os.Getenv("PORT") // Heroku requirement
	if port == "" {
		port = "8080"
	}
	log.Debug("port:", port)

	token, err := femtioelva.GetAccessToken(apikey)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Retrived token:", token)

	// Work loop: query update draw dump sleep
	// seen := []femtioelva.Vehicle{}
	go func() {
		for {
			vs, err := femtioelva.GetVehicleLocations(token, box)
			if err != nil {
				log.Fatal(err)
			}
			for _, p := range vs {
				grid.IncrUTM(femtioelva.LatLong2UTM(p.Lat, p.Long))
			}

			<-time.After(10 * time.Second)
		}
	}()

	// HTTP server
	http.Handle("/", http.HandlerFunc(serveAssets))
	http.Handle("/vasttrafik.png", http.HandlerFunc(serveGrid))
	http.ListenAndServe(":"+port, nil)

	select {}
}
