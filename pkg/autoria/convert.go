package autoria

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (api *API) convert(endpoint string, values url.Values) (url.Values, error) {
	resp, err := http.Get(api.buildURL(endpoint, values))

	if err != nil {
		return nil, err
	}

	var data interface{}

	if resp.StatusCode == http.StatusOK {
		bodyAsByteArray, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(bodyAsByteArray, &data)

		if err != nil {
			return nil, err
		}

		converted := data.(map[string]interface{})["converted"].(map[string]interface{})

		values := url.Values{}
		for k, v := range converted {
			values.Add(k, fmt.Sprint(v))
		}

		return values, nil
	}

	err = fmt.Errorf("invalid response code: %d", resp.StatusCode)
	return nil, err
}

func (api *API) ConvertNewToOld(values url.Values) (url.Values, error) {
	return api.convert("/new_to_old", values)
}

func (api *API) ConvertOldToNew(values url.Values) (url.Values, error) {
	return api.convert("/old_to_new", values)
}
