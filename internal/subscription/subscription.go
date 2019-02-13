// Package subscription implements user subscription.
package subscription

// Subscription represents user subscription to car search.
type Subscription struct {
	Cars    []string          `json:"cars"`
	Active  bool              `json:"active"`
	Params  map[string]string `json:"params"`
	quitter chan struct{}
}

// New creates new clean Subscription.
func New(params map[string]string) *Subscription {
	return &Subscription{
		Cars:    make([]string, 0),
		quitter: make(chan struct{}),
		Active:  false,
		Params:  params,
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
func (sub *Subscription) Start(callback func(chan struct{})) {
	if sub.Active {
		close(sub.quitter)
		sub.quitter = make(chan struct{})
	}

	go func(quit chan struct{}) {
		for {
			select {
			case <-quit:
				return
			default:
				callback(sub.quitter)
			}
		}
	}(sub.quitter)

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
