package main

import (
	"flag"
	"os"
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

	seen := make(map[string]femtioelva.Vehicle)
	for {
		vs, err := femtioelva.GetVehicleLocations(token)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("Livemap queried vehicles:", len(vs))

		updated := 0
		for _, v := range vs {
			old, ok := seen[v.Gid]
			if ok {
				log.Debug("Already seen GID: ", v.Gid)
				log.Debugf("OLD: %#v\n", old)
				log.Debugf("NEW: %#v\n", v)
			} else {
				updated++
				seen[v.Gid] = v
			}
		}
		log.Info("Livemap updated vehicles:", updated)

		time.Sleep(time.Minute)
	}
}
