package domain

import (
	"context"

	"github.com/opencars/bot/pkg/domain/model"
	"github.com/opencars/toolkit"
)

type VehicleService interface {
	FindByNumber(ctx context.Context, number string) (*model.Result, error)
	FindByVIN(ctx context.Context, vin string) (*model.Result, error)
	FindByImage(ctx context.Context, url string) (*model.Result, error)
}

type Recognizer interface {
	Recognize(ctx context.Context, url string) ([]toolkit.ResultALPR, error)
}

type MessageService interface {
	Create(ctx context.Context, msg *model.Message) error
}
