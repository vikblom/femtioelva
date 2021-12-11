package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type oAuth2Response struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	// other fields are not relevant yet.
}

func getAccessToken(apikey string) (string, error) {

	url := "https://api.vasttrafik.se/token"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	// add headers
	secret := base64.URLEncoding.EncodeToString([]byte(apikey))
	req.Header.Set("Authorization", "Basic "+secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// specify params
	q := req.URL.Query()
	q.Add("format", "json")
	q.Add("grant_type", "client_credentials")
	req.URL.RawQuery = q.Encode()

	log.Debug(req.URL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var authResp oAuth2Response
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return "", err
	}
	if authResp.ExpiresIn < 60 {
		return "", errors.New("token will expire in less than a minute")
	}
	return authResp.AccessToken, nil
}

func getStopId(stop, token string) (int, error) {
	url := "https://api.vasttrafik.se/bin/rest.exe/v2/location.name"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	q := req.URL.Query()
	q.Add("format", "json")
	q.Add("input", stop)
	req.URL.RawQuery = q.Encode()

	log.Debug(req.URL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	log.Debug(resp.Status)

	return 0, nil
}

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

	token, err := getAccessToken(strings.TrimSpace(apikey))
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Retrived token:", token)

	id, err := getStopId("svingeln", token)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Svingeln has ID:", id)
}
