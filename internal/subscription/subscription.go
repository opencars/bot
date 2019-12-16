// Package subscription implements user subscription.
package subscription

import "time"

// Subscription represents user subscription to car search.
type Subscription struct {
	Cars    []string
	Active  bool
	period  time.Duration
	quitter chan struct{}
}

// New creates new clean Subscription.
func New(period time.Duration) *Subscription {
	return &Subscription{
		Cars:    make([]string, 0),
		quitter: make(chan struct{}),
		Active:  false,
		period:  period,
	}
}

// Stop stops all running goroutines for current subscription.
func (sub *Subscription) Stop() {
	// Avoid closing closed channel.
	if !sub.Active {
		return
	}

	close(sub.quitter)
	sub.quitter = make(chan struct{})
	sub.Active = false
}

// Start initializes a new goroutine for current subscription with specified callback.
func (sub *Subscription) Start(callback func()) {
	if sub.Active {
		close(sub.quitter)
		sub.quitter = make(chan struct{})
	}

	go func() {
		for {
			select {
			case <-sub.quitter:
				return
			default:
				callback()
				<-time.After(sub.period)
			}
		}
	}()

	sub.Active = true
}

// CarExists returns whatever car included in the list of cars for current subscription.
func (sub *Subscription) CarExists(ID string) bool {
	for _, carID := range sub.Cars {
		if carID == ID {
			return true
		}
	}

	return false
}

// NewCars returns list of cars, which aren't in the list yet.
func (sub *Subscription) NewCars(cars []string) []string {
	newCars := make([]string, 0)

	for _, newCar := range cars {
		if !sub.CarExists(newCar) {
			newCars = append(newCars, newCar)
		}
	}

	return newCars
}
