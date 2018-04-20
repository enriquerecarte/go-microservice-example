package eventbus

import (
	"fmt"
	"reflect"
	"sync"
)

//BusSubscriber defines subscription-related bus behavior
type BusSubscriber interface {
	Subscribe(fn interface{}) error
	SubscribeAsync(fn interface{}, transactional bool) error
	SubscribeOnce(fn interface{}) error
	SubscribeOnceAsync(fn interface{}) error
	Unsubscribe(handler interface{}) error
}

//BusPublisher defines publishing-related bus behavior
type BusPublisher interface {
	Publish(args interface{})
}

//BusController defines bus control behavior (checking handler's presence, synchronization)
type BusController interface {
	WaitAsync()
}

//Bus englobes global (subscribe, publish, control) bus behavior
type Bus interface {
	BusController
	BusSubscriber
	BusPublisher
}

// EventBus - box for handlers and callbacks.
type EventBus struct {
	handlers map[string][]*eventHandler
	lock     sync.Mutex // a lock for the map
	wg       sync.WaitGroup
}

type eventHandler struct {
	callBack      reflect.Value
	flagOnce      bool
	async         bool
	transactional bool
	sync.Mutex // lock for an event handler - useful for running async callbacks serially
}

func New() *EventBus {
	return &EventBus{
		make(map[string][]*eventHandler),
		sync.Mutex{},
		sync.WaitGroup{},
	}
}

// doSubscribe handles the subscription logic and is utilized by the public Subscribe functions
func (bus *EventBus) doSubscribe(fn interface{}, handler *eventHandler) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	subscriberType := reflect.TypeOf(fn)
	if !(subscriberType.Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", subscriberType.Kind())
	}
	if !(subscriberType.NumIn() == 1) {
		return fmt.Errorf("%s must only have one argument", subscriberType.Kind())
	}
	argumentType := subscriberType.In(0).String()
	bus.handlers[argumentType] = append(bus.handlers[argumentType], handler)
	return nil
}

// Subscribe subscribes to a topic.
// Returns error if `fn` is not a function.
func (bus *EventBus) Subscribe(fn interface{}) error {
	return bus.doSubscribe(fn, &eventHandler{
		reflect.ValueOf(fn), false, false, false, sync.Mutex{},
	})
}

// SubscribeAsync subscribes to a topic with an asynchronous callback
// Transactional determines whether subsequent callbacks for a topic are
// run serially (true) or concurrently (false)
// Returns error if `fn` is not a function.
func (bus *EventBus) SubscribeAsync(fn interface{}, transactional bool) error {
	return bus.doSubscribe(fn, &eventHandler{
		reflect.ValueOf(fn), false, true, transactional, sync.Mutex{},
	})
}

// SubscribeOnce subscribes to a topic once. Handler will be removed after executing.
// Returns error if `fn` is not a function.
func (bus *EventBus) SubscribeOnce(fn interface{}) error {
	return bus.doSubscribe(fn, &eventHandler{
		reflect.ValueOf(fn), true, false, false, sync.Mutex{},
	})
}

// SubscribeOnceAsync subscribes to a topic once with an asynchronous callback
// Handler will be removed after executing.
// Returns error if `fn` is not a function.
func (bus *EventBus) SubscribeOnceAsync(fn interface{}) error {
	return bus.doSubscribe(fn, &eventHandler{
		reflect.ValueOf(fn), true, true, false, sync.Mutex{},
	})
}

// Unsubscribe removes callback defined for a topic.
// Returns error if there are no callbacks subscribed to the topic.
func (bus *EventBus) Unsubscribe(handler interface{}) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	subscriberType := reflect.TypeOf(handler)
	fmt.Println("Number of arguments:", subscriberType.NumIn())
	if !(subscriberType.Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", subscriberType.Kind())
	}
	if !(subscriberType.NumIn() == 1) {
		return fmt.Errorf("%s must only have one argument", subscriberType.Kind())
	}

	argumentType := subscriberType.In(0).String()
	if _, ok := bus.handlers[argumentType]; ok && len(bus.handlers[argumentType]) > 0 {
		bus.removeHandler(argumentType, bus.findHandlerIdx(argumentType, reflect.ValueOf(handler)))
		return nil
	}
	return fmt.Errorf("there are no handlers for type %s", argumentType)
}

// Publish executes callback defined for a topic. Any additional argument will be transferred to the callback.
func (bus *EventBus) Publish(args interface{}) {
	bus.lock.Lock() // will unlock if handler is not found or always after setUpPublish
	defer bus.lock.Unlock()

	argumentType := reflect.TypeOf(args).String()

	if handlers, ok := bus.handlers[argumentType]; ok && 0 < len(handlers) {
		// Handlers slice may be changed by removeHandler and Unsubscribe during iteration,
		// so make a copy and iterate the copied slice.
		copyHandlers := make([]*eventHandler, 0, len(handlers))
		copyHandlers = append(copyHandlers, handlers...)
		for i, handler := range copyHandlers {
			if handler.flagOnce {
				bus.removeHandler(argumentType, i)
			}
			if !handler.async {
				bus.doPublish(handler, argumentType, args)
			} else {
				bus.wg.Add(1)
				if handler.transactional {
					handler.Lock()
				}
				go bus.doPublishAsync(handler, argumentType, args)
			}
		}
	}
}

func (bus *EventBus) doPublish(handler *eventHandler, topic string, args interface{}) {
	passedArguments := bus.setUpPublish(topic, args)
	handler.callBack.Call(passedArguments)
}

func (bus *EventBus) doPublishAsync(handler *eventHandler, topic string, args interface{}) {
	defer bus.wg.Done()
	if handler.transactional {
		defer handler.Unlock()
	}
	bus.doPublish(handler, topic, args)
}

func (bus *EventBus) removeHandler(topic string, idx int) {
	if _, ok := bus.handlers[topic]; !ok {
		return
	}
	l := len(bus.handlers[topic])

	if !(0 <= idx && idx < l) {
		return
	}

	copy(bus.handlers[topic][idx:], bus.handlers[topic][idx+1:])
	bus.handlers[topic][l-1] = nil // or the zero value of T
	bus.handlers[topic] = bus.handlers[topic][:l-1]
}

func (bus *EventBus) findHandlerIdx(topic string, callback reflect.Value) int {
	if _, ok := bus.handlers[topic]; ok {
		for idx, handler := range bus.handlers[topic] {
			if handler.callBack == callback {
				return idx
			}
		}
	}
	return -1
}

func (bus *EventBus) setUpPublish(topic string, args ...interface{}) []reflect.Value {

	passedArguments := make([]reflect.Value, 0)
	for _, arg := range args {
		passedArguments = append(passedArguments, reflect.ValueOf(arg))
	}
	return passedArguments
}

// WaitAsync waits for all async callbacks to complete
func (bus *EventBus) WaitAsync() {
	bus.wg.Wait()
}
