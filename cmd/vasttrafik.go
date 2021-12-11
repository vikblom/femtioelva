package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type oAuth2Response struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	// other fields are not relevant yet.
}

func dumpResponse(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	println(body)
	return nil
}

func getAccessToken(apikey string) string {

	url := "https://api.vasttrafik.se/token"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		panic(err)
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

	println(req.URL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var authResp oAuth2Response
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		panic(err)
	}
	if authResp.ExpiresIn < 60 {
		panic("Token will expire in less than a minute!")
	}
	return authResp.AccessToken
}

func getStopId(stop, token string) int {
	url := "https://api.vasttrafik.se/bin/rest.exe/v2/location.name"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	q := req.URL.Query()
	q.Add("format", "json")
	q.Add("input", stop)
	req.URL.RawQuery = q.Encode()

	println(req.URL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	println(resp.Status)

	//dumpResponse(resp)

	return 0
}

func main() {
	apikey := os.Getenv("VASTTRAFIKAPI")
	if apikey == "" {
		println("Could not read API key from env: VASTTRAFIKAPI")
		os.Exit(1)
	}

	token := getAccessToken(strings.TrimSpace(apikey))
	println(token)
	println("Bearer " + token)
	println(getStopId("svingeln", token))
}
