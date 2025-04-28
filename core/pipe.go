package core

import (
	"fmt"
	"reflect"
	"strings"
)

// Pipe is a function that can transform or validate data
type Pipe func(interface{}) (interface{}, error)

// PipeOptions defines pipe configuration
type PipeOptions struct {
	Type reflect.Type
}

// PipeMeta stores pipe metadata
type PipeMeta struct {
	Type reflect.Type
	Pipe Pipe
}

var pipes []PipeMeta

// UsePipe decorator for method parameter
func UsePipe(pipe Pipe, options PipeOptions) func(interface{}, string, int) {
	return func(target interface{}, propertyKey string, paramIndex int) {
		pipes = append(pipes, PipeMeta{
			Type: options.Type,
			Pipe: pipe,
		})
	}
}

// GetPipes returns all registered pipes
func GetPipes() []PipeMeta {
	return pipes
}

// ValidationPipe example pipe
func ValidationPipe() Pipe {
	return func(value interface{}) (interface{}, error) {
		// Validate value
		if value == nil {
			return nil, fmt.Errorf("value cannot be nil")
		}
		return value, nil
	}
}

// ParseIntPipe example pipe
func ParseIntPipe() Pipe {
	return func(value interface{}) (interface{}, error) {
		// Try to adapt the value to int
		result, err := Adapt(value, reflect.TypeOf(0))
		if err != nil {
			return nil, fmt.Errorf("failed to parse int: %v", err)
		}
		return result, nil
	}
}

// ParseBoolPipe example pipe
func ParseBoolPipe() Pipe {
	return func(value interface{}) (interface{}, error) {
		// Try to adapt the value to bool
		result, err := Adapt(value, reflect.TypeOf(false))
		if err != nil {
			return nil, fmt.Errorf("failed to parse bool: %v", err)
		}
		return result, nil
	}
}

// TrimPipe example pipe
func TrimPipe() Pipe {
	return func(value interface{}) (interface{}, error) {
		switch v := value.(type) {
		case string:
			return strings.TrimSpace(v), nil
		default:
			return value, nil
		}
	}
}

// DefaultPipe provides a default value if the input is nil
func DefaultPipe(defaultValue interface{}) Pipe {
	return func(value interface{}) (interface{}, error) {
		if value == nil {
			return defaultValue, nil
		}
		return value, nil
	}
}

// RangePipe validates that a number is within a range
func RangePipe(min, max float64) Pipe {
	return func(value interface{}) (interface{}, error) {
		// Try to adapt the value to float64
		num, err := Adapt(value, reflect.TypeOf(float64(0)))
		if err != nil {
			return nil, fmt.Errorf("value must be a number: %v", err)
		}
		
		f := num.(float64)
		if f < min || f > max {
			return nil, fmt.Errorf("value must be between %v and %v", min, max)
		}
		
		return value, nil
	}
}

// LengthPipe validates string or slice length
func LengthPipe(min, max int) Pipe {
	return func(value interface{}) (interface{}, error) {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.String, reflect.Slice, reflect.Array:
			length := v.Len()
			if length < min || length > max {
				return nil, fmt.Errorf("length must be between %d and %d", min, max)
			}
		default:
			return nil, fmt.Errorf("value must be a string, slice, or array")
		}
		return value, nil
	}
}

// RegexPipe validates a string against a regular expression
func RegexPipe(pattern string) Pipe {
	return func(value interface{}) (interface{}, error) {
		str, err := Adapt(value, reflect.TypeOf(""))
		if err != nil {
			return nil, fmt.Errorf("value must be a string: %v", err)
		}
		
		// TODO: Implement regex validation
		// This is a placeholder - you would need to implement actual regex validation
		return str, nil
	}
}

// CustomPipe creates a custom pipe with a validation function
func CustomPipe(validate func(interface{}) error) Pipe {
	return func(value interface{}) (interface{}, error) {
		if err := validate(value); err != nil {
			return nil, err
		}
		return value, nil
	}
} 