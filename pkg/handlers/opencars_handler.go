package handlers

import (
	"bytes"
	"html/template"
	"log"

	"github.com/opencars/bot/pkg/openalpr"
	"github.com/opencars/toolkit"
)

type OpenCarsHandler struct {
	client     *toolkit.Client
	recognizer *openalpr.API
}

func NewOpenCarsHandler(client *toolkit.Client, recognizer *openalpr.API) *OpenCarsHandler {
	return &OpenCarsHandler{
		client:     client,
		recognizer: recognizer,
	}
}

func (h OpenCarsHandler) getInfoByPlates(number string) (string, error) {
	operations, err := h.client.Operation().FindByNumber(number)
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

func (h OpenCarsHandler) getRegistrations(number string) (string, error) {
	registration, err := h.client.Registration().FindByNumber(number)
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
		Number        string
	}{
		registration,
		number,
	}); err != nil {
		return "", err
	}

	return buff.String(), nil
}
