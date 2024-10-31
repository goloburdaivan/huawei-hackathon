package events

type Event interface{}

type Listener interface {
	Handle(event Event)
}

type EventDispatcher struct {
	listeners map[string][]Listener
}

var dispatcher = &EventDispatcher{
	listeners: make(map[string][]Listener),
}

func GetDispatcher() *EventDispatcher {
	return dispatcher
}

func (d *EventDispatcher) Register(eventName string, listener Listener) {
	d.listeners[eventName] = append(d.listeners[eventName], listener)
}

func (d *EventDispatcher) Dispatch(eventName string, event Event) {
	if listeners, found := d.listeners[eventName]; found {
		for _, listener := range listeners {
			listener.Handle(event)
		}
	}
}
