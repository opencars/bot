package handlers

import (
	"bytes"
	"html/template"
	"log"

	"github.com/opencars/bot/pkg/openalpr"
	"github.com/opencars/toolkit/sdk"
)

type OpenCarsHandler struct {
	OpenCars   *sdk.Client
	Recognizer *openalpr.API
}

func (h OpenCarsHandler) getInfoByPlates(plate string) (string, error) {
	transport, err := h.OpenCars.SearchLimit(plate, 5)
	if err != nil {
		return "", err
	}

	tpl, err := template.ParseFiles("templates/opencars_info.tpl")
	if err != nil {
		return "", err
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, struct {
		Cars   []sdk.Operation
		Number string
	}{
		transport,
		plate,
	}); err != nil {
		log.Println(err)
	}

	return buff.String(), nil
}
