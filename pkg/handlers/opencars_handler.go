package handlers

import (
	"bytes"
	"html/template"
	"log"

	"github.com/opencars/bot/pkg/openalpr"
	"github.com/opencars/toolkit"
)

type OpenCarsHandler struct {
	OpenCars   *toolkit.Client
	Recognizer *openalpr.API
}

func (h OpenCarsHandler) getInfoByPlates(number string) (string, error) {
	operations, err := h.OpenCars.Operations(number, 5)
	if err != nil {
		return "", err
	}

	tpl, err := template.ParseFiles("templates/opencars/operations.tpl")
	if err != nil {
		return "", err
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, struct {
		Operations []toolkit.Operation
		Number     string
	}{
		operations,
		number,
	}); err != nil {
		log.Println(err)
	}

	return buff.String(), nil
}

func (h OpenCarsHandler) getRegistrations(code string) (string, error) {
	registrations, err := h.OpenCars.Registrations(code)
	if err != nil {
		return "", err
	}

	tpl, err := template.ParseFiles("templates/opencars/registrations.tpl")
	if err != nil {
		return "", err
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, struct {
		Registrations []toolkit.Registration
		Code          string
	}{
		registrations,
		code,
	}); err != nil {
		log.Println(err)
	}

	return buff.String(), nil
}
