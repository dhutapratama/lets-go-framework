package boot

import (
	"lets-go-framework/lets"
	"lets-go-framework/lets/drivers"
	"lets-go-framework/lets/frameworks"
	"lets-go-framework/lets/loader"
)

// List of initializer
var Initializer = []func(){
	loader.Environment,
	loader.Logger,
	loader.Launching,
}

// List of framework that start on lets
var Servers = []func(){
	frameworks.Http,
	frameworks.Grpc,
	frameworks.RabbitMQ,
	frameworks.Tcp,
	drivers.MySQL,
	drivers.Redis,
	drivers.MongoDB,
}

// Add initialization function and run before application starting
func AddInitializer(init func()) {
	Initializer = append(Initializer, init)
}

// Bootstrap vars and configuration
func OnInit() {
	for _, initializer := range Initializer {
		// fmt.Printf("%v. Initializing %s\n", i, runtime.FuncForPC(reflect.ValueOf(initializer).Pointer()).Name())
		initializer()
	}
}

// Bootstrap frameworks
func OnMain() {
	lets.LogI("Booting ...")
	for _, runner := range Servers {
		go runner()
	}

	loader.RunningForever()
}
