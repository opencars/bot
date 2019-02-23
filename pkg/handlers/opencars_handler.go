package handlers

import (
	"bytes"
	"html/template"
	"log"

	"github.com/shal/opencars-bot/pkg/openalpr"
	"github.com/shal/opencars-bot/pkg/opencars"
)

type OpenCarsHandler struct {
	OpenCars   *opencars.API
	Recognizer *openalpr.API
}

func (h OpenCarsHandler) getInfoByPlates(plate string) (string, error) {
	transport, err := h.OpenCars.Search(plate)
	if err != nil {
		return "", err
	}

	tpl, err := template.ParseFiles("templates/opencars_info.tpl")
	if err != nil {
		return "", err
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, struct {
		Cars   []opencars.Transport
		Number string
	}{
		transport,
		plate,
	}); err != nil {
		log.Println(err)
	}

	return buff.String(), nil
}
