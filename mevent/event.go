package mevent

import (
	"github.com/healer1219/martini/global"
	"log"
	"sync"
)

type EventType string

type Event interface {
	// OnEvent onEvent
	OnEvent(ctx *global.Context)
}

type Bus struct {
	sync.Mutex
	eventMap map[EventType]*EventGroup
}

func (b *Bus) AddBlockEvent(name EventType, event Event) {
	defer b.Unlock()
	b.Lock()
	eventGroup, ok := b.eventMap[name]
	if !ok {
		eventGroup = NewEventGroup()
		b.eventMap[name] = eventGroup
	}
	eventGroup.addBlock(event)
}

func (b *Bus) AddEvent(name EventType, event Event) {
	defer b.Unlock()
	b.Lock()
	eventGroup, ok := b.eventMap[name]
	if !ok {
		eventGroup = NewEventGroup()
		b.eventMap[name] = eventGroup
	}
	eventGroup.addSync(event)
}

func (b *Bus) Publish(eventType EventType) {
	defer b.Unlock()
	b.Lock()
	eventGroup, ok := b.eventMap[eventType]
	if !ok {
		log.Printf("event %v not exists", eventType)
		return
	}
	eventGroup.publish()
}

type EventGroup struct {
	syncGroup  []Event
	blockGroup []Event
}

func NewEventGroup() *EventGroup {
	return &EventGroup{
		syncGroup:  make([]Event, 0),
		blockGroup: make([]Event, 0),
	}
}

func (e *EventGroup) addSync(event Event) {
	if e.syncGroup == nil {
		e.syncGroup = make([]Event, 0)
	}
	e.syncGroup = append(e.syncGroup, event)
}

func (e *EventGroup) addBlock(event Event) {
	if e.blockGroup == nil {
		e.blockGroup = make([]Event, 0)
	}
	e.blockGroup = append(e.blockGroup, event)
}

func (e *EventGroup) publish() {
	for idx := range e.syncGroup {
		event := e.syncGroup[idx]
		if event != nil {
			go event.OnEvent(global.Ctx())
		}
	}

	for idx := range e.blockGroup {
		event := e.blockGroup[idx]
		if event != nil {
			event.OnEvent(global.Ctx())
		}
	}
}

var bus = Bus{
	eventMap: make(map[EventType]*EventGroup),
}

func Publish(eventType EventType) {
	bus.Publish(eventType)
}

func Add(name EventType, event Event) {
	bus.AddEvent(name, event)
}

func AddBlock(name EventType, event Event) {
	bus.AddBlockEvent(name, event)
}
