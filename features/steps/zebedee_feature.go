package steps

import (
	"github.com/maxcnunes/httpfake"
)

// ZebedeeFeature contains all the information for a fake component API
type ZebedeeFeature struct {
	ErrorFeature
	FakeAPI *httpfake.HTTPFake
}

func NewZebedeeFeature() *ZebedeeFeature {
	f := &ZebedeeFeature{
		FakeAPI: httpfake.New(),
	}

	return f
}

func (f *ZebedeeFeature) setJSONResponseForGetHealth(url string, statusCode int) {
	f.FakeAPI.NewHandler().Get(url).Reply(statusCode)
}

func (f *ZebedeeFeature) setJSONResponseForGetPublishIndex(statusCode int) {
	f.FakeAPI.NewHandler().Get("/publishedindex").Reply(statusCode).BodyString(`{}`)
}

// Close closes the fake API
func (f *ZebedeeFeature) Close() {
	f.FakeAPI.Close()
}

// Reset resets the fake API
func (f *ZebedeeFeature) Reset() {
	f.FakeAPI.Reset()
}
