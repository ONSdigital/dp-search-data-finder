package steps

import (
	"github.com/maxcnunes/httpfake"
)

// SearchReindexFeature contains all the information for a fake component API
type SearchReindexFeature struct {
	ErrorFeature
	FakeSearchAPI *httpfake.HTTPFake
}

func NewSearchReindexFeature() *SearchReindexFeature {
	f := &SearchReindexFeature{
		FakeSearchAPI: httpfake.New(),
	}

	return f
}

func (f *SearchReindexFeature) setJSONResponseForGetHealth(url string, statusCode int) {
	f.FakeSearchAPI.NewHandler().Get(url).Reply(statusCode)
}

// Close closes the fake API
func (f *SearchReindexFeature) Close() {
	f.FakeSearchAPI.Close()
}

// Reset resets the fake API
func (f *SearchReindexFeature) Reset() {
	f.FakeSearchAPI.Reset()
}
