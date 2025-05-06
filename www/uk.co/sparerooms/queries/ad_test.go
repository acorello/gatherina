package queries

import (
	"net/url"
	"os"
	"testing"

	"dev.acorello.it/go/gatherina/must"
	"github.com/stretchr/testify/assert"
)

func TestAd(t *testing.T) {
	var input = must.Must(os.Open("testdata/samples/ad_17151018.html"))
	defer input.Close()

	var expectedError error = nil
	var expectedAdDetails = AdDetails{
		ListingId: "17151018",
		Location: Location{
			Latitude:  51.4517223,
			Longitude: -0.1969331,
		},
	}
	gotAd, err := Ad(input)
	if err != expectedError {
		t.Errorf("Ad() error = %v, expectedError %v", err, expectedError)
		return
	}
	assert.Equal(t, expectedAdDetails, gotAd)
}

func TestGetAd(t *testing.T) {
	const sampleAd = "https://www.spareroom.co.uk/flatshare/flatshare_detail.pl?flatshare_id=9523438"
	adURL := must.Must(url.Parse(sampleAd))

	rc, err := GetAd(*adURL)
	if err != nil {
		t.Errorf("GetAd failed: %v", err)
	}
	defer rc.Close()

	gotAd, err := Ad(rc)
	if err != nil {
		t.Errorf("query AdDetails failed: %v", err)
	}
	expected := AdDetails{
		ListingId: "9523438",
		Location: Location{
			Latitude:  51.495941365157,
			Longitude: -0.008709348379069990,
		},
	}
	assert.Equal(t, expected, gotAd)
}
