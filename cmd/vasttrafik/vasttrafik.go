package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/vikblom/femtioelva"
)

func main() {
	apikey := os.Getenv("VASTTRAFIKAPI")
	if apikey == "" {
		log.Fatal("Could not read API key from env: VASTTRAFIKAPI")
		os.Exit(1)
	}

	verboseFlag := flag.Bool("v", false, "verbose logging")
	flag.Parse()

	if *verboseFlag {
		log.SetLevel(log.DebugLevel)
		log.Debug("Verbose prints enabled")
	}

	if flag.NArg() == 0 {
		fmt.Println("usage: vasttrafik [path to .gob]")
		os.Exit(1)
	}
	gobfile := flag.Args()[0]
	log.Debug("Will write data to ", gobfile)

	token, err := femtioelva.GetAccessToken(apikey)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Retrived token:", token)

	// Ctrl-c should break the loop immediately.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	box := femtioelva.GeoBox(femtioelva.GBG_LAT, femtioelva.GBG_LON, 10_000)

	alive := true
	seen := []femtioelva.Vehicle{}
	for alive {
		vs, err := femtioelva.GetVehicleLocations(token, box)
		if err != nil {
			log.Fatal(err)
		}
		seen = append(seen, vs...)
		log.Infof("accumulated %d samples from livemap", len(seen))
		if len(seen) > 1_000_000 {
			break
		}

		select {
		case <-time.After(10 * time.Second):
		case <-sigs:
			alive = false
		}
	}

	// Dump history to gob before exiting.
	err = os.MkdirAll(path.Dir(gobfile), 0755)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(gobfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	err = gob.NewEncoder(file).Encode(seen)
	if err != nil {
		log.Fatal(err)
	}
}
