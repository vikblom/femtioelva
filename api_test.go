package femtioelva_test

import (
	"os"
	"testing"

	"github.com/vikblom/femtioelva"
)

var apikey string

func TestApi(t *testing.T) {
	apikey = os.Getenv("VASTTRAFIKAPI")
	if apikey == "" {
		t.Skip("No apikey in env var VASTTRAFIKAPI, skipping.")
	}

	t.Run("GetAccessToken", testGetAccessToken)
}

func testGetAccessToken(t *testing.T) {
	_, err := femtioelva.GetAccessToken(apikey)
	if err != nil {
		t.Error("Retrieving an access token did not work.")
	}
}
