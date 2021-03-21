package mocks

//go:generate mockgen -destination=./service.go -package=mocks github.com/opencars/bot/pkg/domain VehicleService,MessageService
