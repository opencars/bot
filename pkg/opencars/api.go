package opencars

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type API struct {
	URI string
}

func New(URI string) *API {
	return &API{URI: URI}
}

func (client *API) Search(number string) ([]Transport, error) {
	if len(strings.TrimSpace(number)) == 0 {
		return nil, errors.New("number is empty")
	}

	query := fmt.Sprintf("%s/transport?number=%s", client.URI, number)
	response, err := http.Get(query)

	if err != nil {
		return nil, err
	}

	transport := make([]Transport, 0)
	if err := json.NewDecoder(response.Body).Decode(&transport); err != nil {
		return nil, errors.New("invalid response body")
	}

	return transport, nil
}
