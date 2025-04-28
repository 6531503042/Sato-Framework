package core

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gofiber/fiber/v2"
)

// Container manages dependencies
type Container struct {
	services sync.Map
	mu       sync.RWMutex
}

// NewContainer creates a new dependency container
func NewContainer() *Container {
	return &Container{}
}

// Register registers a service with the container
func (c *Container) Register(name string, service interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.services.Load(name); exists {
		return fmt.Errorf("service %s already registered", name)
	}

	c.services.Store(name, service)
	return nil
}

// Get retrieves a service from the container
func (c *Container) Get(name string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	service, exists := c.services.Load(name)
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}

	return service, nil
}

// Inject injects dependencies into a struct
func (c *Container) Inject(target interface{}) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return ErrInvalidTarget
	}

	val = val.Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if tag, ok := fieldType.Tag.Lookup("inject"); ok {
			service, err := c.Get(tag)
			if err != nil {
				return fmt.Errorf("failed to inject %s: %v", tag, err)
			}

			serviceVal := reflect.ValueOf(service)
			if !serviceVal.Type().AssignableTo(field.Type()) {
				return fmt.Errorf("service %s is not assignable to field %s", tag, fieldType.Name)
			}

			field.Set(serviceVal)
		}
	}

	return nil
}

// Inject errors
var (
	ErrInvalidTarget = fiber.NewError(fiber.StatusInternalServerError, "invalid injection target")
) 