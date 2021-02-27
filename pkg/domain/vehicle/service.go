package vehicle

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/opencars/bot/pkg/domain"
	"github.com/opencars/bot/pkg/logger"
)

type Service struct {
	operation    domain.OperationService
	registration domain.RegistrationService
}

func NewService(o domain.OperationService, r domain.RegistrationService) *Service {
	return &Service{
		operation:    o,
		registration: r,
	}
}

func (s *Service) FindByNumber(ctx context.Context, number string) (string, error) {
	g, gCtx := errgroup.WithContext(ctx)

	var operations []domain.Operation
	var registrations []domain.Registration

	g.Go(func() error {
		var err error
		operations, err = s.operation.FindByNumber(gCtx, number)
		if err != nil {
			return err
		}

		return nil
	})

	// TODO: Add field "is_active" to the registration,
	//       then if is_active == true, we can find all registrations for given vin_code.
	g.Go(func() error {
		var err error
		registrations, err = s.registration.FindByNumber(gCtx, number)
		if err != nil {
			return err
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	logger.Debugf("registrations: %#v", registrations)
	logger.Debugf("operations: %#v", operations)

	// TODO: Build the HTML page with following data:
	// - operations.
	// - registrations.
	return fmt.Sprintf("%#v %#v", registrations, operations), nil
}

func (s *Service) FindByVIN(ctx context.Context, vin string) (string, error) {
	registrations, err := s.registration.FindByVIN(ctx, vin)
	if err != nil {
		return "", err // TODO: Wrap error into ErrRemoteFailed!
	}

	var operations []domain.Operation
	for _, r := range registrations {
		var err error
		tmp, err := s.operation.FindByNumber(ctx, r.Number)
		if err != nil {
			return "", nil
		}

		operations = append(operations, tmp...)
	}

	logger.Debugf("registrations: %#v", registrations)
	logger.Debugf("operations: %#v", operations)

	// TODO: Build the HTML page with following data:
	// - operations.
	// - registrations.
	return fmt.Sprintf("%#v %#v", registrations, operations), nil
}
