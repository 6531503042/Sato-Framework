package core

import (
	"fmt"
	"reflect"
	"sync"
)

// ModuleOptions defines module configuration
type ModuleOptions struct {
	Imports     []interface{}
	Controllers []interface{}
	Providers   []interface{}
	Exports     []interface{}
}

// ModuleMeta stores module metadata
type ModuleMeta struct {
	Instance    interface{}
	Imports     []interface{}
	Controllers []interface{}
	Providers   []interface{}
	Exports     []interface{}
}

var modules []ModuleMeta
var mu sync.RWMutex

// Module decorator for class
func Module(options ModuleOptions) func(m interface{}) interface{} {
	return func(m interface{}) interface{} {
		mu.Lock()
		defer mu.Unlock()

		// Validate module
		if err := validateModule(m, options); err != nil {
			panic(err)
		}

		// Register module
		modules = append(modules, ModuleMeta{
			Instance:    m,
			Imports:     options.Imports,
			Controllers: options.Controllers,
			Providers:   options.Providers,
			Exports:     options.Exports,
		})

		return m
	}
}

// GetModules returns all registered modules
func GetModules() []ModuleMeta {
	mu.RLock()
	defer mu.RUnlock()
	return modules
}

// GetModule returns a module by type
func GetModule(moduleType interface{}) (ModuleMeta, error) {
	mu.RLock()
	defer mu.RUnlock()

	for _, m := range modules {
		if reflect.TypeOf(m.Instance) == reflect.TypeOf(moduleType) {
			return m, nil
		}
	}

	return ModuleMeta{}, fmt.Errorf("module not found")
}

func validateModule(m interface{}, options ModuleOptions) error {
	// Validate imports
	for _, imp := range options.Imports {
		if _, err := GetModule(imp); err != nil {
			return fmt.Errorf("invalid import: %v", err)
		}
	}

	// Validate controllers
	for _, ctrl := range options.Controllers {
		if reflect.TypeOf(ctrl).Kind() != reflect.Ptr {
			return fmt.Errorf("controller must be a pointer")
		}
	}

	// Validate providers
	for _, prov := range options.Providers {
		if reflect.TypeOf(prov).Kind() != reflect.Ptr {
			return fmt.Errorf("provider must be a pointer")
		}
	}

	return nil
} 