package domain

import (
	"context"

	"github.com/opencars/toolkit"
)

type VehicleService interface {
	FindByNumber(ctx context.Context, number string) (*Result, error)
	FindByVIN(ctx context.Context, vin string) (*Result, error)
	FindByImage(ctx context.Context, url string) (*Result, error)
}

type Recognizer interface {
	Recognize(ctx context.Context, url string) ([]toolkit.ResultALPR, error)
}

type MessageService interface {
	Create(ctx context.Context, msg *Message) error
}
