# gogogo
A flexible framework for building RESTFul service in Go

This framework promotes plugin-based architecture pattern and layered architecture.

While plugin based architecture makes application extendable, layered architecture makes it clean.
It makes it easy to separate concerns, application logic and business logic can be developer separately. 
These concerned can be integrated later without affecting each other.


**The motivation** <br/>
As application grows more and more code is pushed into it but hardly any code removed.
The plugin based approach allows you to manage code easily. Every new feature is a plugin, so is every
obsolete feature. One could simply be plugged into the system and other can
simply be plugged out from the system.  


### Components

The framework offers many core components as explained below.


#### router
Router is one of the core component. It allows to define api endpoints to be exposed externally

**main.go**

```Go 
import (
	"github.com/udlabs/gogogo/router"
	"github.com/udlabs/gogogo/server"
	"log"
	"os"
)

func main() {
    path, err := os.Getwd() // find the current directory
	if err != nil {
		log.Println(err)
	}
    
    // instantial the router
	r := new(router.HttpRouter)
	
	// load routes
	r.LoadRoute(path, "" /* default file name route_config.json.template */)
	routerHandler := r.BindChiRouter()

	// Setup and start http server
	s := server.HttpServer{Port: ":3005"}
	s.Start(routerHandler)
}

```

#### server
It's a typical http server module. Provide apis to start the server with configure the routing.

**main.go**

```Go 
import (
	"github.com/udlabs/gogogo/router"
	"github.com/udlabs/gogogo/server"
	"log"
	"os"
)

func main() {
    path, err := os.Getwd() // find the current directory
	if err != nil {
		log.Println(err)
	}
    
    // instantial the router
	r := new(router.HttpRouter)
	
	// load routes
	r.LoadRoute(path, "" /* default file name route_config.json.template */)
	routerHandler := r.BindChiRouter()

	// Setup and start http server
	s := server.HttpServer{Port: ":3005"}
	s.Start(routerHandler)
}

```


### Maintainer

Upendra Dev Singh 
