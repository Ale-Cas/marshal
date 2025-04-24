# marshal

Marshal is a Go library to encode and decode structs through HTTP.

## Features

*   **HTTP Encoding/Decoding:** Seamlessly encode Go structs into HTTP request bodies (e.g., JSON, XML) and decode HTTP response bodies into Go structs. Simplifies API interactions.
*   **Automatic Type Conversion:** Handles common type conversions between Go structs and HTTP data formats.
*   **Customizable:** Offers options to customize encoding/decoding behavior, such as specifying custom field names or data formats.
*   **Error Handling:** Provides robust error handling for common encoding/decoding issues.

## Example

Here's a basic example of how to use the `marshal` package to make a GET request and decode the response:

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/Ale-Cas/marshal"
)

// ExpectedResponse represents the expected HTTP response body.
type ExpectedResponse struct {
    Data string `json:"data"`
}

func main() {

    // Make a GET request
    resp, err := marshal.Get[ExpectedResponse](http.DefaultClient, "https://example.com/api/data", nil)
    if err != nil {
        panic(err)
    }

    // Print the data from the response
    fmt.Println("Data:", resp.Data)
}
```

## API

- `Request[Body, Response any](client Client, method HTTPMethod, url string, body Body, headers Headers) (*Response, error)`
    - Encode the body struct as JSON.
    - Send HTTP request with the specified parameters.
    - Decode the response body from JSON into the Response type.
- `RequestWithContext[Body, Response any](client Client, method HTTPMethod, url string, body Body, headers Headers, ctx context.Context) (*Response, error)`
    - Same as Request, but includes a context.Context for cancellation and timeouts.

Alongside the generic request functions, the library provides convenience methods for common HTTP verbs:
- `Post` sends a POST request, automatically encoding the request body and decoding the response.
- `Get` sends a GET request with automatic JSON decoding.
- `Put` sends a PUT request, automatically encoding the request body and decoding the response.
- `Patch` sends a PATCH request, automatically encoding the request body and decoding the response.
- `Delete` sends a DELETE request.