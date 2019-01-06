// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

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
	return fmt.Sprintf("%s/%s?api_key=%s&%s", api.base, path, api.key, options)
}
