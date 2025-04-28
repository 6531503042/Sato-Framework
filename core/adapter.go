package core

import (
	"fmt"
	"reflect"
	"sync"
)

// AdapterFunc is a function that can adapt one type to another
type AdapterFunc func(interface{}) (interface{}, error)

// AdapterRegistry stores registered adapters
type AdapterRegistry struct {
	adapters map[string]AdapterFunc
	mu       sync.RWMutex
}

// NewAdapterRegistry creates a new adapter registry
func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		adapters: make(map[string]AdapterFunc),
	}
}

// RegisterAdapter registers a new adapter
func (r *AdapterRegistry) RegisterAdapter(from, to reflect.Type, adapter AdapterFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	key := fmt.Sprintf("%s->%s", from.String(), to.String())
	r.adapters[key] = adapter
}

// GetAdapter retrieves an adapter for the given types
func (r *AdapterRegistry) GetAdapter(from, to reflect.Type) (AdapterFunc, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	key := fmt.Sprintf("%s->%s", from.String(), to.String())
	adapter, exists := r.adapters[key]
	return adapter, exists
}

// Adapt attempts to adapt a value to the target type
func (r *AdapterRegistry) Adapt(value interface{}, targetType reflect.Type) (interface{}, error) {
	valueType := reflect.TypeOf(value)
	
	// If types match, return the value directly
	if valueType == targetType {
		return value, nil
	}
	
	// Try to find a direct adapter
	if adapter, exists := r.GetAdapter(valueType, targetType); exists {
		return adapter(value)
	}
	
	// Try to find an indirect adapter through intermediate types
	// This is a simplified version - in a real implementation, you'd want to
	// find the shortest path between types
	for _, adapter := range r.adapters {
		adapted, err := adapter(value)
		if err != nil {
			continue
		}
		
		adaptedType := reflect.TypeOf(adapted)
		if adaptedType == targetType {
			return adapted, nil
		}
		
		// Try to adapt the intermediate result
		if final, err := r.Adapt(adapted, targetType); err == nil {
			return final, nil
		}
	}
	
	return nil, fmt.Errorf("no adapter found from %s to %s", valueType, targetType)
}

// Global adapter registry
var globalRegistry = NewAdapterRegistry()

// RegisterAdapter registers a global adapter
func RegisterAdapter(from, to reflect.Type, adapter AdapterFunc) {
	globalRegistry.RegisterAdapter(from, to, adapter)
}

// Adapt attempts to adapt a value to the target type using the global registry
func Adapt(value interface{}, targetType reflect.Type) (interface{}, error) {
	return globalRegistry.Adapt(value, targetType)
}

// Built-in adapters
func init() {
	// String to int
	RegisterAdapter(
		reflect.TypeOf(""),
		reflect.TypeOf(0),
		func(v interface{}) (interface{}, error) {
			str := v.(string)
			var result int
			_, err := fmt.Sscanf(str, "%d", &result)
			if err != nil {
				return nil, fmt.Errorf("failed to convert string to int: %v", err)
			}
			return result, nil
		},
	)
	
	// Int to string
	RegisterAdapter(
		reflect.TypeOf(0),
		reflect.TypeOf(""),
		func(v interface{}) (interface{}, error) {
			return fmt.Sprintf("%d", v.(int)), nil
		},
	)
	
	// String to bool
	RegisterAdapter(
		reflect.TypeOf(""),
		reflect.TypeOf(false),
		func(v interface{}) (interface{}, error) {
			str := v.(string)
			switch str {
			case "true", "1", "yes":
				return true, nil
			case "false", "0", "no":
				return false, nil
			default:
				return nil, fmt.Errorf("invalid boolean string: %s", str)
			}
		},
	)
	
	// Bool to string
	RegisterAdapter(
		reflect.TypeOf(false),
		reflect.TypeOf(""),
		func(v interface{}) (interface{}, error) {
			return fmt.Sprintf("%v", v.(bool)), nil
		},
	)
} 