package domain

import (
	"context"

	"github.com/opencars/grpc/pkg/core"
)

type VehicleService interface {
	FindByNumber(ctx context.Context, number string) (*core.Result, error)
	FindByVIN(ctx context.Context, vin string) (*core.Result, error)
}

type MessageService interface {
	Create(ctx context.Context, msg *Message) error
}
