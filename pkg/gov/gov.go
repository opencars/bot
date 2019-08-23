package gov

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	// DefaultHost is an immutable default host of government data registry.
	DefaultHost = "https://data.gov.ua"
	// DefaultBasePath is an immutable base path of government data registry.
	DefaultBasePath = "/api/3/action"
)

var (
	// BaseHost equals to DefaultHost by default, it is mutable.
	BaseHost = DefaultHost
	// BasePath equals to DefaultBasePath by default, it is mutable.
	BasePath = DefaultBasePath
)

// Client is core http client struct.
type Client struct {
	http *http.Client
}

// NewClient creates new http client with timeout equal to 5 seconds.
func NewClient() *Client {
	api := new(Client)

	api.http = &http.Client{
		Timeout: 5 * time.Second,
	}

	return api
}

func (c *Client) get(path string) (*Response, error) {
	// TODO: Replace this with url builder, aka url.Url{}, url.String().
	url := fmt.Sprintf("%s/%s", BaseHost, path)

	raw, err := c.http.Get(url)
	if err != nil {
		return nil, err
	}

	response := new(Response)
	if err := json.NewDecoder(raw.Body).Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}

// Package returns detailed information about package from government registry.
// Makes a simple HTTP Get request under the hood.
func (c *Client) Package(id string) (*Package, error) {
	// TODO: Replace this with path builder, aka path.Join(...)
	path := fmt.Sprintf("%s/package_show?id=%s", DefaultBasePath, id)
	res, err := c.get(path)
	if err != nil {
		return nil, err
	}

	pkg := new(Package)

	if err := json.Unmarshal(res.Result, pkg); err != nil {
		return nil, err
	}

	return pkg, nil
}
