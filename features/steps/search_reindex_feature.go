package steps

import (
	"github.com/maxcnunes/httpfake"
)

// SearchReindexFeature contains all the information for a fake component API
type SearchReindexFeature struct {
	ErrorFeature
	FakeAPI *httpfake.HTTPFake
}

func NewSearchReindexFeature() *SearchReindexFeature {
	f := &SearchReindexFeature{
		FakeAPI: httpfake.New(),
	}

	return f
}

func (f *SearchReindexFeature) setJSONResponseForGetHealth(url string, statusCode int) {
	f.FakeAPI.NewHandler().Get(url).Reply(statusCode)
}

// Close closes the fake API
func (f *SearchReindexFeature) Close() {
	f.FakeAPI.Close()
}

// Reset resets the fake API
func (f *SearchReindexFeature) Reset() {
	f.FakeAPI.Reset()
	// f.setJSONResponseForGetHealth("/health", 200)
}
