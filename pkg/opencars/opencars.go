package opencars

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type API struct {
	URI string
}

type Transport struct {
	ID                  int    `json:"id"`
	RegistrationAddress string `json:"registration_address"`
	RegistrationCode    int    `json:"registration_code"`
	Registration        string `json:"registration"`
	Date                string `json:"date"`
	Model               string `json:"model"`
	Year                int    `json:"year"`
	Color               string `json:"color"`
	Kind                string `json:"kind"`
	Body                string `json:"body"`
	Fuel                string `json:"fuel"`
	Capacity            int    `json:"capacity"`
	Weight              int    `json:"own_weight"`
	Number              string `json:"number"`
}

func (client *API) Search(number string) ([]Transport, error) {
	query := fmt.Sprintf("%s/transport?number=%s", client.URI, number)
	response, err := http.Get(query)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	transport := make([]Transport, 0)
	if err = json.Unmarshal(body, &transport); err != nil {
		return nil, err
	}

	return transport, nil
}
