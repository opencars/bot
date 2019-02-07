package autoria

import (
	"fmt"
	"strings"
)

type API struct {
	key  string
	base string
}

func NewAPI(key string) *API {
	return &API{
		key:  key,
		base: "https://developers.ria.com",
	}
}

func (api *API) BuildURL(path string, params ...string) string {
	options := strings.Join(params, "&")

	if len(options) > 0 {
		return fmt.Sprintf("%s/%s?api_key=%s&%s", api.base, path, api.key, options)
	} else {
		return fmt.Sprintf("%s/%s?api_key=%s", api.base, path, api.key)
	}
}
