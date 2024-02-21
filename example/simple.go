package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JHelar/dumbo/dumbo"
)

var D = dumbo.New()

func Layout(children dumbo.ComponentChildren, props dumbo.ComponentProps) (dumbo.DumboElement, error) {
	template := dumbo.Must(D.ParseFile("./example/layout.html"))

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

	outFile, err := os.Create("./example/output.html")
	if err != nil {
		log.Panicf("Error creating file %s", err.Error())
	}
	defer outFile.Close()

	err = pageTemplate.Execute(outFile, map[string]string{
		"pageTitle": "Simple example",
		"heading":   "World",
	})

	if err != nil {
		log.Panicf("Error executing template %s", err.Error())
	}

	log.Print("Render completed see: ./example/output.html")
}
