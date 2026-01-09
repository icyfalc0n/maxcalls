package calls

import (
	"io"
	"net/http"
	"net/url"
)

const (
	ENDPOINT        = "https://calls.okcdn.ru/fb.do"
	APPLICATION_KEY = "CNHIJPLGDIHBABABA"
)

type RawApiClient struct{}

func (c *RawApiClient) CallMethod(method string, additionalParams map[string]string) ([]byte, error) {
	form := url.Values{}
	form.Set("method", method)
	form.Set("format", "JSON")
	form.Set("application_key", APPLICATION_KEY)

	for k, v := range additionalParams {
		form.Set(k, v)
	}

	resp, err := http.PostForm(ENDPOINT, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
