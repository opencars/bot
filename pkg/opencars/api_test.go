package opencars

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_Search(t *testing.T) {
	fake := make([]Transport, 0)

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
		_, err := api.Search("AX1234BT")
		assert.Error(t, err)
	})

	t.Run("number is empty", func(t *testing.T) {
		_, err := New(okServer.URL).Search(" ")
		assert.EqualError(t, err, "number is empty")
	})

	t.Run("invalid response body", func(t *testing.T) {
		_, err := New(okServer.URL).Search("AX1234BT")
		assert.EqualError(t, err, "invalid response body")
	})

	t.Run("everything is valid", func(t *testing.T) {
		resp, err := New(jsonServer.URL).Search("AX1234BT")
		assert.NoError(t, err)
		assert.Equal(t, fake, resp)
	})
}
