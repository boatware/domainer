# domainer

[![Go Reference](https://pkg.go.dev/badge/github.com/boatware/domainer.svg)](https://pkg.go.dev/github.com/boatware/domainer)

Simple Go library to split URLs into their domain parts.

## Installation

```bash
go get github.com/boatware/domainer
```

## Usage

```go
package main

import (
    "fmt"

    "github.com/boatware/domainer"
)

func main() {
    url := "http://www.example.com:8080/path/to/file.html?query=string#fragment"

    d, _ := domainer.FromString(url)

    fmt.Println(d.Protocol) // http
    fmt.Println(d.Subdomain) // www
    fmt.Println(d.Domain) // example
    fmt.Println(d.TLD) // com
    fmt.Println(d.Port) // 8080
    fmt.Println(d.Path) // path/to/file.html
    fmt.Println(d.Query) // []Query{ Query{ Key: "query", Value: "string" } }
    fmt.Println(d.Fragment) // fragment
}
```