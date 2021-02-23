package autoria

import (
	"encoding/json"
	"net/http"

	"github.com/opencars/bot/pkg/logger"
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

	logger.Debugf("%v", resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	vv, ok := res.Data[ID].(map[string]interface{})
	if ok {
		for _, v := range vv {
			buff, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}

			var photo Photo
			if err := json.Unmarshal(buff, &photo); err != nil {
				return nil, err
			}

			res.Photos = append(res.Photos, photo)
		}
	}

	return res, nil
}
