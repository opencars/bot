// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package opencars

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type API struct {
	URL string
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
	OwnWeight           int    `json:"own_weight"`
	Number              string `json:"number"`
}

func (client *API) SearchTransport(number string) ([]Transport, error) {
	query := fmt.Sprintf("%s/transport?number=%s", client.URL, number)
	fmt.Println(query)
	resp, err := http.Get(query)

	fmt.Println(resp)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	transport := make([]Transport, 0)

	err = json.Unmarshal(body, &transport)

	if err != nil {
		return nil, err
	}

	return transport, nil
}