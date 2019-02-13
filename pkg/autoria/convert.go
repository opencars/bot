package autoria

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (api *API) convert(endpoint string, params map[string]string) (map[string]string, error) {
	strParams := make([]string, len(params))

	for k, v := range params {
		strParams = append(strParams, fmt.Sprintf("%s=%s", k, v))
	}

	resp, err := http.Get(api.BuildURL(endpoint, strParams...))

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

		res := make(map[string]string)
		for k, v := range converted {
			res[k] = fmt.Sprint(v)
		}

		return res, nil
	}

	err = NewErr(fmt.Sprintf("invalid response code: %d", resp.StatusCode))
	return nil, err
}

func (api *API) ConvertNewToOld(params map[string]string) (map[string]string, error) {
	return api.convert("/new_to_old", params)
}

func (api *API) ConvertOldToNew(params map[string]string) (map[string]string, error) {
	return api.convert("/old_to_new", params)
}
