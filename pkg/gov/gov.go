package gov

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"sort"
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

func (c *Client) get(endpoint string, values url.Values) (*Response, error) {
	u, err := url.Parse(BaseHost)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(path.Join(u.Path, BasePath), endpoint)
	u.RawQuery = values.Encode()

	raw, err := c.http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer raw.Body.Close()

	response := new(Response)
	if err := json.NewDecoder(raw.Body).Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}

// Package returns detailed information about package from government registry.
// Makes a simple HTTP Get request under the hood.
func (c *Client) Package(id string) (*Package, error) {
	res, err := c.get("package_show", url.Values{
		"id": {id},
	})
	if err != nil {
		return nil, err
	}

	pkg := new(Package)

	if err := json.Unmarshal(res.Result, pkg); err != nil {
		return nil, err
	}

	// Sort by modification time.
	sort.Sort(Sorter(pkg.Resources))

	return pkg, nil
}
