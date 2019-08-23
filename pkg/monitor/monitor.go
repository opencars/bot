package monitor

import (
	"time"

	"github.com/opencars/bot/pkg/gov"
)

type Source interface {
	Package(id string) (*gov.Package, error)
}

// Monitor is responsible from monitoring government data registry.
type Monitor struct {
	pkg       string
	user      int64
	timestamp gov.Time
	Events    chan *gov.Package

	source Source
}

// New creates new instance of Monitor.
func New(
	pkg string,
	user int64,
	source Source,
) (*Monitor, error) {
	return &Monitor{
		pkg:  pkg,
		user: user,
		// Month ago.
		timestamp: gov.Time{Time: time.Now().Add(-730 * time.Hour)},
		source:    source,
		Events:    make(chan *gov.Package),
	}, nil
}

func (monitor *Monitor) Listen() error {
	for {
		res, err := monitor.source.Package(monitor.pkg)
		if err != nil {
			return err
		}

		resource := res.Resources[len(res.Resources)-1]
		if resource.LastModified.After(monitor.timestamp.Time) {
			// Update latest timestamp.
			monitor.timestamp = resource.LastModified
			monitor.Events <- res
		}
	}
}
