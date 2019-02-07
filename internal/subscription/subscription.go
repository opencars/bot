package subscription

import (
	"log"
)

type Subscription struct {
	Cars    []string          `json:"cars"`
	Active  bool              `json:"active"`
	Params  map[string]string `json:"params"`
	quitter chan struct{}     `json:"-"`
}

func NewSubscription(params map[string]string) *Subscription {
	lastIDs := make([]string, 0)
	quitter := make(chan struct{})

	return &Subscription{
		Cars:    lastIDs,
		quitter: quitter,
		Active:  false,
		Params:  params,
	}
}

func (s *Subscription) Stop() {
	// Avoid closing closed channel.
	if !s.Active {
		return
	}

	close(s.quitter)
	s.quitter = make(chan struct{})
	s.Active = false
}

func (s *Subscription) Start(callback func(chan struct{})) {
	if s.Active {
		log.Println("Subscription is active, recreating...")
		close(s.quitter)
		s.quitter = make(chan struct{})
	}

	go func(quit chan struct{}) {
		for {
			select {
			case <-quit:
				log.Println("Quit was called")
				return
			default:
				callback(s.quitter)
			}
		}
	}(s.quitter)

	s.Active = true
}

func (s *Subscription) carExist(carID string) bool {
	for _, car := range s.Cars {
		if car == carID {
			return true
		}
	}

	return false
}

func (s *Subscription) GetNewCars(cars []string) []string {
	newCars := make([]string, 0)

	for _, newCar := range cars {
		if !s.carExist(newCar) {
			newCars = append(newCars, newCar)
		}
	}

	return newCars
}
