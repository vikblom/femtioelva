package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type oAuth2Response struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	// other fields are not relevant yet.
}

const KEY = "u0TPd1wPLc4_2P8JIofbIqfSn3Ia"
const SECRET = "VnTZtM_dHUM2kwQN7CaEU0sPXaYa"

func dumpResponse(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	println(body)
	return nil
}

func getAccessToken() string {
	url := "https://api.vasttrafik.se/token"
	secret := base64.URLEncoding.EncodeToString([]byte(KEY + ":" + SECRET))

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		panic(err)
	}
	// add headers
	req.Header.Set("Authorization", "Basic "+secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// specify params
	q := req.URL.Query()
	q.Add("format", "json")
	q.Add("grant_type", "client_credentials")
	req.URL.RawQuery = q.Encode()

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
	token := getAccessToken()
	println(token)
	println("Bearer " + token)
	println(getStopId("svingeln", token))
}
