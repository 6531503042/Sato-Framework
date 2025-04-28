package core

import (
	"reflect"
	"strings"
)

// HTTP Methods
const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	PATCH  = "PATCH"
)

// ControllerOptions defines controller configuration
type ControllerOptions struct {
	Path    string
	Guards  []Guard
	Version string
}

// RouteOptions defines route configuration
type RouteOptions struct {
	Path   string
	Method string
}

// ControllerMeta stores controller metadata
type ControllerMeta struct {
	Instance interface{}
	Path     string
	Guards   []Guard
	Version  string
	Routes   []RouteMeta
}

// RouteMeta stores route metadata
type RouteMeta struct {
	Path     string
	Method   string
	Handler  string
	Guards   []Guard
	Pipes    []PipeMeta
}

var controllers []ControllerMeta

// Controller decorator for class
func Controller(options ControllerOptions) func(c interface{}) interface{} {
	return func(c interface{}) interface{} {
		path := options.Path
		if path == "" {
			path = "/" + strings.ToLower(reflect.TypeOf(c).Name())
		}

		controllers = append(controllers, ControllerMeta{
			Instance: c,
			Path:     path,
			Guards:   options.Guards,
			Version:  options.Version,
			Routes:   make([]RouteMeta, 0),
		})

		return c
	}
}

// Get decorator for method
func Get(options RouteOptions) func(interface{}, string) {
	return func(target interface{}, propertyKey string) {
		addRoute(target, RouteMeta{
			Path:    options.Path,
			Method:  GET,
			Handler: propertyKey,
		})
	}
}

// Post decorator for method
func Post(options RouteOptions) func(interface{}, string) {
	return func(target interface{}, propertyKey string) {
		addRoute(target, RouteMeta{
			Path:    options.Path,
			Method:  POST,
			Handler: propertyKey,
		})
	}
}

// Put decorator for method
func Put(options RouteOptions) func(interface{}, string) {
	return func(target interface{}, propertyKey string) {
		addRoute(target, RouteMeta{
			Path:    options.Path,
			Method:  PUT,
			Handler: propertyKey,
		})
	}
}

// Delete decorator for method
func Delete(options RouteOptions) func(interface{}, string) {
	return func(target interface{}, propertyKey string) {
		addRoute(target, RouteMeta{
			Path:    options.Path,
			Method:  DELETE,
			Handler: propertyKey,
		})
	}
}

// Patch decorator for method
func Patch(options RouteOptions) func(interface{}, string) {
	return func(target interface{}, propertyKey string) {
		addRoute(target, RouteMeta{
			Path:    options.Path,
			Method:  PATCH,
			Handler: propertyKey,
		})
	}
}

// UseGuards decorator for method
func UseGuards(guards ...Guard) func(interface{}, string) {
	return func(target interface{}, propertyKey string) {
		addGuards(target, propertyKey, guards)
	}
}

// UsePipes decorator for method
func UsePipes(pipes ...PipeMeta) func(interface{}, string) {
	return func(target interface{}, propertyKey string) {
		addPipes(target, propertyKey, pipes)
	}
}

func addRoute(target interface{}, route RouteMeta) {
	for i, c := range controllers {
		if c.Instance == target {
			controllers[i].Routes = append(controllers[i].Routes, route)
			break
		}
	}
}

func addGuards(target interface{}, handler string, guards []Guard) {
	for i, c := range controllers {
		if c.Instance == target {
			for j, r := range c.Routes {
				if r.Handler == handler {
					controllers[i].Routes[j].Guards = guards
					break
				}
			}
			break
		}
	}
}

func addPipes(target interface{}, handler string, pipes []PipeMeta) {
	for i, c := range controllers {
		if c.Instance == target {
			for j, r := range c.Routes {
				if r.Handler == handler {
					controllers[i].Routes[j].Pipes = pipes
					break
				}
			}
			break
		}
	}
}

// GetControllers returns all registered controllers
func GetControllers() []ControllerMeta {
	return controllers
}

// ApplyPipes applies pipes to a value
func ApplyPipes(value interface{}, pipes []PipeMeta) (interface{}, error) {
	var err error
	result := value
	
	for _, pipe := range pipes {
		result, err = pipe.Pipe(result)
		if err != nil {
			return nil, err
		}
	}
	
	return result, nil
}

