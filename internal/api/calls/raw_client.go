package calls

import (
	"io"
	"net/http"
	"net/url"
)

const (
	endpoint       = "https://calls.okcdn.ru/fb.do"
	applicationKey = "CNHIJPLGDIHBABABA"
	endpointFormat = "JSON"
)

type RawApiClient struct{}

func (c *RawApiClient) CallMethod(method string, additionalParams map[string]string) ([]byte, error) {
	form := url.Values{}
	form.Set("method", method)
	form.Set("format", endpointFormat)
	form.Set("application_key", applicationKey)

	for k, v := range additionalParams {
		form.Set(k, v)
	}

	resp, err := http.PostForm(endpoint, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
