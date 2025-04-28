package core

type (
	ControllerOptions struct {
		RoutePrefix string
		Guards []Guard
	}

	ControllerMeta struct {
		Instance interface{}
		RoutePrefix string
		Guards []Guard
	}
)

var controllers []ControllerMeta

func Controller(options ControllerOptions) func(c interface{}) interface{} {
	return func(c interface{}) interface{} {
		controllers = append(controllers, ControllerMeta{
			Instance: c,
			RoutePrefix: options.RoutePrefix,
			Guards: options.Guards,
		})

		return c
	}
}

func GetControllers() []ControllerMeta {
	return controllers
}

