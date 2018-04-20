package events

import (
	"github.com/enriquerecarte/microservices-example/pkg/eventbus"
)

var ev *eventbus.EventBus

func StartEvents() {
	ev = eventbus.New()
}

func Publish(event interface{}) {
	ev.Publish(event)
}

func Subscribe(eventListener interface{}) {
	ev.Subscribe(eventListener)
}
