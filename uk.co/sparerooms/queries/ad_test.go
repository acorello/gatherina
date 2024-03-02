package queries

import (
	"dev.acorello.it/go/gatherina/must"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAd(t *testing.T) {
	var input = must.Must(os.Open("/Users/am/Projects/gatherina/uk.co/sparerooms/queries/samples/ad_17151018.html"))
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
