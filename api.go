package femtioelva

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// A square grid on top of Gothenburg
	BOX = GeoBox(GBG_LAT, GBG_LON, 10_000)
)

type oAuth2Response struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	// other fields are not relevant yet.
}

// apiCall wraps the vasttrafik API, does basic error checking on the response.
// If the response is non-nil, the caller is responsible for closing its body.
func apiCall(method, path string, headers, params map[string]string) (*http.Response, error) {
	url := &url.URL{
		Scheme: "https",
		Host:   "api.vasttrafik.se",
		Path:   path,
	}
	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// request - let caller check error and defer close
	log.Debug(req.URL.String())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		log.Debug(string(raw))
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP request not ok: %s", resp.Status)
	}
	return resp, nil
}

func GetAccessToken(apikey string) (string, error) {
	secret := base64.URLEncoding.EncodeToString([]byte(apikey))
	headers := map[string]string{
		"Authorization": "Basic " + secret,
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	params := map[string]string{
		"format":     "json",
		"grant_type": "client_credentials",
	}
	resp, err := apiCall(http.MethodPost, "token", headers, params)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// decode
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

type Vehicle struct {
	// Some kind of ID.
	// Occurs across requests but other fields can be different
	// even if this is the same. Does it identify a vehicle?
	Gid string

	// "Y*1e6" in Västtrafiks API, increases north.
	Lat float64

	// "X*1e6" in Västtrafiks API, increases east.
	Long float64

	// Common name of transport.
	Name string

	// When this data was retrieved.
	Time time.Time
}

// VasttrafikCoord matches the API which uses stringed int millionths.
func apiCoord(coord float64) string {
	return strconv.Itoa(int(1e6 * coord))
}

func GetVehicleLocations(token string, box Box) ([]Vehicle, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}
	log.Info(headers)

	params := map[string]string{
		"format":       "json",
		"miny":         apiCoord(box.LowLat),
		"maxy":         apiCoord(box.HighLat),
		"minx":         apiCoord(box.LowLong),
		"maxx":         apiCoord(box.HighLong),
		"onlyRealtime": "yes",
	}
	resp, err := apiCall(http.MethodGet, "bin/rest.exe/v2/livemap", headers, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode JSON with some anon. structs.
	type VehicleJSON struct {
		X    string
		Y    string
		Name string
		Gid  string
	}
	var lm struct {
		Livemap struct {
			Vehicles []VehicleJSON
		}
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&lm)
	if err != nil {
		return nil, err
	}

	// Pack output
	vs := []Vehicle{}
	for _, v := range lm.Livemap.Vehicles {

		// Store as the actual Long/Lat
		x, err := strconv.Atoi(v.X)
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(v.Y)
		if err != nil {
			return nil, err
		}

		vs = append(vs, Vehicle{
			Gid:  v.Gid,
			Lat:  float64(y) * 1e-6,
			Long: float64(x) * 1e-6,
			Name: v.Name,
			Time: time.Now(),
		})
	}
	return vs, nil
}
