# Example Usage

```
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bxdn/flite"
)

type credentials struct {
	Password string `json:"password"`
}

func login(f *flite.Flite) error {
	creds, e := flite.GetTypedBody[credentials](f.Req.Context())
	if e != nil {
		f.ReturnError("Internal Server Error", http.StatusInternalServerError)
		return e
	}
	fmt.Println(creds.Password)
	if creds.Password != "pass" {
		return f.ReturnText("YOU FAIL")
	} else {
		return f.ReturnText("YOU MAY ENTER")
	}
}

func simpleGreet(f *flite.Flite) error {
	return f.ReturnText("Hello Universe!")
}

func getCredentials(f *flite.Flite) error {
	return f.ReturnJSON(credentials{"pass"})
}

func main() {

	// Creates an endpoint at the root path ("/") and adds a GET handler to it
	endpoint := flite.CreateEndpoint("/")
	endpoint.GET(simpleGreet)

	// Creates an endpoint at "/credentials" and adds a GET handler to it
	endpoint2 := flite.CreateEndpoint("/credentials")
	endpoint2.GET(getCredentials)

	// Creates an endpoint at "/login" and adds a POST handler to it
	// The flite.Json middleware will convert the incoming json body to the credentials type before reaching the login() handler
	endpoint3 := flite.CreateEndpoint("/login")
	endpoint3.POST(flite.Json[credentials], login)

	// Creates the server, registers the endpoints, and serves at port 8080
	s := flite.NewFliteServer()
	s.Register(endpoint, endpoint2, endpoint3)
	if e := s.Serve(8080); e != nil {
		log.Fatal(e)
	}
}
```