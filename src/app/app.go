package app

import (
	"sync"

	"magic-mouse-engine/client/httpclient"
	"magic-mouse-engine/data"
)

var (
	app  *Receipt
	once sync.Once
)

// Receipt structure
type Receipt struct {
	SVCPrefix  string
	HTTPClient *httpclient.HTTPClient
}

// Create an object of app
func Create(prefix string,
	c *httpclient.HTTPClient,
	cfg *data.Config) *Receipt {

	f := func() {
		app = &Receipt{
			SVCPrefix:  prefix,
			HTTPClient: c,
		}
	}
	once.Do(f)
	return app
}

// Get return object.Object would have been created using Create().
func Get() *Receipt {
	return app
}

