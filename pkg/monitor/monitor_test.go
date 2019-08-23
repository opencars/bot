package monitor

import (
	"testing"
	"time"

	"github.com/opencars/bot/pkg/gov"
)

type Fake struct {
	ok bool
}

func (f *Fake) Package(id string) (*gov.Package, error) {
	if f.ok {
		return &gov.Package{
			Resources: []gov.Resource{
				{
					LastModified: gov.Time{Time: time.Now()},
				},
				{
					LastModified: gov.Time{Time: time.Now().Add(-1 * time.Second)},
				},
			},
		}, nil
	} else {
		return &gov.Package{
			Resources: []gov.Resource{
				{
					LastModified: gov.Time{Time: time.Now().Add(-731 * time.Hour)},
				},
			},
		}, nil
	}
}

func TestMonitor_ListenOK(t *testing.T) {
	m, err := New("test", 0, &Fake{true})
	if err != nil {
		t.Errorf(err.Error())
	}

	errs := make(chan error)

	go func() {
		err := m.Listen()
		if err != nil {
			errs <- err
		}
	}()

	select {
	case <-errs:
		t.Error(err.Error())
	case pkg := <-m.Events:
		if pkg == nil {
			t.Fail()
		}
	case <-time.After(1 * time.Second):
		t.Fail()
	}
}

func TestMonitor_ListenNothing(t *testing.T) {
	m, err := New("test", 0, &Fake{false})
	if err != nil {
		t.Errorf(err.Error())
	}

	errs := make(chan error)

	go func() {
		err := m.Listen()
		if err != nil {
			errs <- err
		}
	}()

	select {
	case <-errs:
		t.Error(err.Error())
	case <-m.Events:
		t.Fail()
	case <-time.After(1 * time.Second):
		return
	}
}
