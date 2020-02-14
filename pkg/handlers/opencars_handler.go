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

func (h OpenCarsHandler) getInfoByNumber(number string) (string, error) {
	operations, err := h.client.Operation().FindByNumber(number)
	if err != nil {
		return "", err
	}

	tpl, err := template.ParseFiles("templates/operations.tpl")
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

func (h OpenCarsHandler) getRegistrationsByNumber(number string) (string, error) {
	registration, err := h.client.Registration().FindByNumber(number)
	if err != nil {
		return "", err
	}

	tpl, err := template.ParseFiles("templates/registrations.tpl")
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

func (h *OpenCarsHandler) GetReportByVIN(vin string) (string, error) {
	registrations, err := h.client.Registration().FindByVIN(vin)
	if err != nil {
		return "", err
	}

	if len(registrations) == 0 {
		return "Нажаль, не знайшли такого VIN", nil
	}

	operations, err := h.client.Operation().FindByNumber(registrations[0].Number)
	if err != nil {
		return "", err
	}

	tpl, err := template.ParseFiles("templates/report_by_vin.tpl")
	if err != nil {
		return "", err
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, struct {
		Registration toolkit.Registration
		Operations   []toolkit.Operation
		VIN          string
	}{
		Registration: registrations[0],
		Operations:   operations,
		VIN:          vin,
	}); err != nil {
		return "", err
	}

	return buff.String(), nil
}
