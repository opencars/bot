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

type PhotosResponse struct {
	Status int                    `json:"status"`
	Data   map[string]interface{} `json:"data"`
	Photos []Photo
}

func (photo Photo) URL() string {
	return photo.Formats[len(photo.Formats)-1]
}

func (api *API) CarPhotos(ID string) (res *PhotosResponse, err error) {
	resp, err := http.Get(api.buildURL("auto/fotos/"+ID, nil))
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

		if err := json.Unmarshal(buff, &img); err != nil {
			return nil, err
		}

		res.Photos = append(res.Photos, *img)
	}

	return res, nil
}
