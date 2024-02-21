# Dumbo

Dumbo is a Go library for building dynamic HTML templates, with heavy support of component based rendering.

## Installation

To install Dumbo, you need to install Go and set your Go workspace first.

1. Download Dumbo:

```sh
go get -u github.com/JHelar/dumbo
```

# Quick Start
Example of how to use Dumbo:

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/JHelar/dumbo/dumbo"
)

var D = dumbo.New()

func Layout(children dumbo.ComponentChildren, props dumbo.ComponentProps) (dumbo.DumboElement, error) {
    template := dumbo.Must(D.ParseFile("layout.html"))

    return template.Render(map[string]interface{}{
        "title":    props["title"],
        "children": children,
    })
}

func Heading(children dumbo.ComponentChildren, props dumbo.ComponentProps) (dumbo.DumboElement, error) {
    template := dumbo.Must(D.Parse(`<h1>{heading}</h1>`))

    return template.Render(map[string]interface{}{
        "heading": fmt.Sprintf("Hello %s", props["heading"]),
    })
}

func main() {
    D.AddComponent("Layout", Layout)
    D.AddComponent("Heading", Heading)

    pageTemplate := dumbo.Must(D.Parse(`
    <Layout title={pageTitle}>
        <Heading heading={heading} />
        <p>From Dumbo</p>
    </Layout>
    `))

    outFile, _ := os.Create("output.html")
    defer outFile.Close()

    pageTemplate.Execute(outFile, map[string]string{
        "pageTitle": "Simple example",
        "heading":   "World",
    })
}
```

# Testing

To run the tests:

```
go test ./...
```