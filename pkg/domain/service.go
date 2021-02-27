package domain

import "context"

type VehicleService interface {
	FindByNumber(ctx context.Context, number string) (string, error)
	FindByVIN(ctx context.Context, vin string) (string, error)
}

type RegistrationService interface {
	FindByNumber(ctx context.Context, number string) ([]Registration, error)
	FindByVIN(ctx context.Context, vin string) ([]Registration, error)
}

type OperationService interface {
	FindByNumber(ctx context.Context, number string) ([]Operation, error)
}
