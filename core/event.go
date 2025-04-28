package core

import (
	"fmt"
	"sync"
	"time"
)

// Event represents an event that can be published
type Event interface {
	GetName() string
	GetTimestamp() time.Time
}

// EventHandler is a function that handles an event
type EventHandler func(Event) error

// EventBus manages event subscriptions and publishing
type EventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe subscribes to an event
func (b *EventBus) Subscribe(eventName string, handler EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	b.handlers[eventName] = append(b.handlers[eventName], handler)
	return nil
}

// Publish publishes an event to all subscribers
func (b *EventBus) Publish(event Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	b.mu.RLock()
	handlers := b.handlers[event.GetName()]
	b.mu.RUnlock()

	var errs []error
	for _, handler := range handlers {
		if err := handler(event); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to handle event: %v", errs)
	}

	return nil
}

// Unsubscribe removes a handler from an event
func (b *EventBus) Unsubscribe(eventName string, handler EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers, exists := b.handlers[eventName]
	if !exists {
		return fmt.Errorf("event %s not found", eventName)
	}

	for i, h := range handlers {
		if &h == &handler {
			b.handlers[eventName] = append(handlers[:i], handlers[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("handler not found for event %s", eventName)
}

// Clear removes all handlers for an event
func (b *EventBus) Clear(eventName string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.handlers[eventName]; !exists {
		return fmt.Errorf("event %s not found", eventName)
	}

	delete(b.handlers, eventName)
	return nil
}

// ClearAll removes all handlers
func (b *EventBus) ClearAll() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers = make(map[string][]EventHandler)
} 