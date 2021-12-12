package femtioelva

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// A square grid on top of Gothenburg
	box      = GeoBox(GBG_LAT, GBG_LON, 10_000)
	MIN_LAT  = int(box.LowLat * 1_000_000)
	MAX_LAT  = int(box.HighLat * 1_000_000)
	MIN_LONG = int(box.LowLong * 1_000_000)
	MAX_LONG = int(box.HighLong * 1_000_000)
)

type oAuth2Response struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	// other fields are not relevant yet.
}

func GetAccessToken(apikey string) (string, error) {
	url := "https://api.vasttrafik.se/token"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	// headers
	secret := base64.URLEncoding.EncodeToString([]byte(apikey))
	req.Header.Set("Authorization", "Basic "+secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// params
	q := req.URL.Query()
	q.Add("format", "json")
	q.Add("grant_type", "client_credentials")
	req.URL.RawQuery = q.Encode()

	// request
	log.Debug(req.URL.String())
	resp, err := http.DefaultClient.Do(req)
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

func GetStopId(stop, token string) (int, error) {
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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	log.Debug(resp.Status)

	return 0, nil
}

type Vehicle struct {
	// Some kind of ID.
	// Occurs across requests but other fields can be different
	// even if this is the same. Does it identify a vehicle?
	Gid string

	// X in Västtrafiks API, increases east.
	Long int

	// Y in Västtrafiks API, increases north.
	Lat int

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
	url := "https://api.vasttrafik.se/bin/rest.exe/v2/livemap"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	strconv.Itoa(123)
	req.Header.Set("Authorization", "Bearer "+token)
	q := req.URL.Query()
	q.Add("format", "json")
	q.Add("miny", apiCoord(box.LowLat))
	q.Add("maxy", apiCoord(box.HighLat))
	q.Add("minx", apiCoord(box.LowLong))
	q.Add("maxx", apiCoord(box.HighLong))
	q.Add("onlyRealtime", "yes")
	req.URL.RawQuery = q.Encode()

	log.Debug(req.URL.String())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Debug(resp.Status)

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
			Long: x,
			Lat:  y,
			Name: v.Name,
			Time: time.Now(),
		})
	}
	return vs, nil
}
