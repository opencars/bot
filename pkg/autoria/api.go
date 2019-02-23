package autoria

import (
	"net/url"
)

type API struct {
	key string
	url url.URL
}

func New(key string) *API {
	return &API{
		key: key,
		url: url.URL{
			Scheme: "https",
			Host:   "developers.ria.com",
		},
	}
}

func (api *API) buildURL(path string, values url.Values) string {
	api.url.Path = path

	if values != nil {
		values.Add("api_key", api.key)
		api.url.RawQuery = values.Encode()
	} else {
		api.url.RawQuery = "api_key" + api.key
	}

	return api.url.String()
}
