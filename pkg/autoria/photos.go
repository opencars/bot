// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package autoria

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Photo struct {
	PhotoID  int      `json:"photo_id"`
	AutoID   int      `json:"auto_id"`
	Status   int      `json:"status"`
	Checked  int      `json:"checked"`
	Standard int      `json:"standard"`
	DateAdd  string   `json:"date_add"`
	URLPath  string   `json:"url"`
	Formats  []string `json:"formats"`
}

type CarPhotosResponse struct {
	Status int                    `json:"status"`
	Data   map[string]interface{} `json:"data"`
	Photos []Photo
}

func (photo Photo) URL() string {
	return photo.Formats[len(photo.Formats) - 1]
}

func (api *API) CarPhotos(ID string) (res *CarPhotosResponse, err error) {
	resp, err := http.Get(api.BuildURL("auto/fotos/" + ID))
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	for k, v := range res.Data[ID].(map[string]interface{}) {
		fmt.Printf("%v\n", k)
		buff, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		img := new(Photo)
		json.Unmarshal(buff, &img)

		res.Photos = append(res.Photos, *img)
	}

	return res, nil
}
