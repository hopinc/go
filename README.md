# hop-go

[View Hop Documentation](https://docs.hop.io/sdks/server/go) | [View Source Documentation](https://pkg.go.dev/github.com/hopinc/hop-go)

Hop's Go library. Requires Go 1.18+.

```go
package main

import (
	"context"
	"fmt"

	"github.com/hopinc/hop-go"
)

func main() {
	myToken := "ptk_xxx"
	c, err := hop.NewClient(myToken)
	if err != nil {
		// Handle errors how you wish here.
		panic(err)
	}

	s, err := c.Projects.Secrets.Create(
		context.Background(),
		"SECRET_NAME",
		"SECRET_VALUE",
		hop.WithProjectID("PROJECT_ID"), // If not using a project token, you will need to specify the project ID.
	)
	if err != nil {
		// Handle errors how you wish here.
		panic(err)
	}
	fmt.Println(s)
}
```

Client options (such as the project ID) can be set either at a client level like `c.AddClientOptions(hop.WithProjectID("PROJECT_ID"))` or at a functional level like shown above. If options are provided to the function, they override the client level option configuration.