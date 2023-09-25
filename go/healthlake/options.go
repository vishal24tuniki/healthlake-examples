package healthlake

import (
	"github.com/eka-care/sdk-go-eka/pkg/internal/httpclient"
)

type Options struct {
	httpClient httpclient.HTTPClient
}

func (t *Options) Copy() Options {
	to := *t
	return to
}

func (t *Options) GetHTTPClient() httpclient.HTTPClient {
	return t.httpClient
}

func (t *Options) SetHTTPClient(httpClient httpclient.HTTPClient) {
	t.httpClient = httpClient
}
