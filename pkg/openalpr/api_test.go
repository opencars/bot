package openalpr

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	imageFixture = "../../test/image-fixture.json"
)

func TestAPI_Recognize(t *testing.T) {
	fake := Image{}

	okServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something went wrong", http.StatusOK)
	}))
	defer okServer.Close()

	jsonServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(fake); err != nil {
			assert.NoError(t, err)
		}
	}))
	defer jsonServer.Close()

	t.Run("server is not running", func(t *testing.T) {
		api := New("http://invalid")
		_, err := api.Recognize("http://localhost:1234")
		assert.Contains(t, err.Error(), "lookup invalid")
	})

	t.Run("empty url", func(t *testing.T) {
		api := New(okServer.URL)
		_, err := api.Recognize("")
		assert.Contains(t, err.Error(), "empty url")
	})

	t.Run("invalid response body", func(t *testing.T) {
		_, err := New(okServer.URL).Recognize("http://localhost:1234")
		assert.Contains(t, err.Error(), "invalid response body")
	})

	t.Run("everything is valid", func(t *testing.T) {
		api := New(jsonServer.URL)
		img, err := api.Recognize("https://golang.org/doc/gopher/project.png")

		assert.NoError(t, err)
		assert.Equal(t, fake, *img)
	})
}

func TestResponse_Plate(t *testing.T) {
	content, err := ioutil.ReadFile(imageFixture)
	assert.NoError(t, err)

	t.Run("nothing found", func(t *testing.T) {
		plates, err := new(Image).Plates()

		assert.Empty(t, plates)
		assert.EqualError(t, err, "no plates found")
	})

	t.Run("few plates found", func(t *testing.T) {
		img := Image{}

		assert.NoError(t, json.Unmarshal(content, &img))
		img.Recognized = append(img.Recognized, img.Recognized[0])
		plates, err := img.Plates()

		assert.Equal(t, 2, len(plates))
		assert.NoError(t, err)
	})

	t.Run("not valid plates found", func(t *testing.T) {
		img := Image{}

		assert.NoError(t, json.Unmarshal(content, &img))
		img.Recognized[0].Candidates[0].Plate = "INVALID"
		plates, err := img.Plates()

		assert.Equal(t, "INVALID", plates[0])
		assert.NoError(t, err)
	})

	t.Run("valid plates found", func(t *testing.T) {
		img := Image{}

		assert.NoError(t, json.Unmarshal(content, &img))
		plates, err := img.Plates()

		assert.Equal(t, 1, len(plates))
		assert.Equal(t, "BH4316IH", plates[0])
		assert.NoError(t, err)
	})
}
