package coordinator

import (
	"time"
)

type EventRaiser interface {
	AddListener(eventName string, f func(interface{}))
}

type EventAggregator struct {
	listeners map[string][]func(interface{})
}

func NewEventAggregator() *EventAggregator {
	ea := EventAggregator{
		listeners: make(map[string][]func(interface{})),
	}
	return &ea
}

func (ea *EventAggregator) AddListener(name string, f func(interface{})) {
	ea.listeners[name] = append(ea.listeners[name], f)
}

func (ea *EventAggregator) PublishEvent(name string, eventData interface{}) {
	if ea.listeners[name] != nil {
		for _, r := range ea.listeners[name] {
			r(eventData)
		}
	}
}

type EventData struct {
	Name      string
	Value     float64
	Timestamp time.Time
}
