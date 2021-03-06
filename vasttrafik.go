package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type oAuth2Response struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	// other fields are not relevant yet.
}

const KEY = "u0TPd1wPLc4_2P8JIofbIqfSn3Ia"
const SECRET = "VnTZtM_dHUM2kwQN7CaEU0sPXaYa"

func main() {
	url := "https://api.vasttrafik.se/token"
	secret := base64.URLEncoding.EncodeToString([]byte(KEY + ":" + SECRET))

	//var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("")))
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
	println(authResp.AccessToken)
	println(authResp.ExpiresIn)
}
