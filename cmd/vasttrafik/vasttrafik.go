package main

import (
	"encoding/gob"
	"flag"
	"os"
	"os/signal"
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

	token, err := femtioelva.GetAccessToken(apikey)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Retrived token:", token)

	// Ctrl-c should break the loop immediately.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	alive := true
	seen := []femtioelva.Vehicle{}
	for alive {
		vs, err := femtioelva.GetVehicleLocations(token)
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("queried %d vehicle from livemap", len(vs))
		seen = append(seen, vs...)

		if len(seen) > 1_000_000 {
			break
		}

		select {
		case <-time.After(15 * time.Second):
		case <-sigs:
			alive = false
		}
	}

	// Dump history to gob before exiting.
	file, err := os.Create("pos.gob")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	err = gob.NewEncoder(file).Encode(seen)
	if err != nil {
		log.Fatal(err)
	}
}
