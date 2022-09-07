# hop-go

[View Source Documentation](https://pkg.go.dev/github.com/hopinc/hop-go)

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
	c, err := hopgo.NewClient(myToken)
	if err != nil {
		// Handle errors how you wish here.
		panic(err)
	}

	s, err := c.Projects.Secrets.Create(
		context.Background(),
		"PROJECT_ID",
		"SECRET_NAME",
		"SECRET_VALUE",
	)
	if err != nil {
		// Handle errors how you wish here.
		panic(err)
	}
	fmt.Println(s)
}
```
