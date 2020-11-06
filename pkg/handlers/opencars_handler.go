package handlers

import (
	"bytes"
	"html/template"

	"github.com/opencars/toolkit"

	"github.com/opencars/bot/pkg/logger"
)

type OpenCarsHandler struct {
	client *toolkit.Client
}

func NewOpenCarsHandler(client *toolkit.Client) *OpenCarsHandler {
	return &OpenCarsHandler{
		client: client,
	}
}

func (h *OpenCarsHandler) getInfoByNumber(number string) (string, error) {
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
		logger.Errorf(err.Error())
	}

	return buff.String(), nil
}

func (h *OpenCarsHandler) getRegistrationsByNumber(number string) (string, error) {
	registrations, err := h.client.Registration().FindByNumber(number)
	if err != nil {
		return "", err
	}

	tpl, err := template.ParseFiles("templates/registrations.tpl")
	if err != nil {
		return "", err
	}

	type payload struct {
		Registrations []toolkit.Registration
		Code, Number  string
	}

	tmp := payload{
		Registrations: registrations,
		Number:        number,
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, tmp); err != nil {
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
