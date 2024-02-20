package template_test

import (
	"log"
	"testing"

	"github.com/JHelar/dumbo/template"
)

func TestTemplateParse(t *testing.T) {
	dumbo := template.New()

	template := dumbo.Parse(`<div hx-test="testing" hx-swap={swap}>test <p>{heading}</p></div>`)

	res := template.Render(map[string]string{
		"heading": "Testing",
		"swap":    "none",
	})

	log.Print(res)
}

type TemplateData struct {
	Heading string
	Swap    string
}

func TestTemplateParseWithStruct(t *testing.T) {
	dumbo := template.New()

	template := dumbo.Parse(`<div hx-test="testing" hx-swap={Swap}>test <p>{Heading}</p></div>`)

	res := template.Render(TemplateData{
		Heading: "testing",
		Swap:    "none",
	})

	log.Print(res)
}

func TestComponentRender(t *testing.T) {
	dumbo := template.New()

	componentTemplate := dumbo.Parse(`
	<div hx-select={select}>
		<h1>{heading}</h1>
		<h2>{subheading}</h2>
		{children}
	</div>`)

	dumbo.AddComponent("Heading", func(children string, props map[string]interface{}) string {
		hxSelect := "Heading select"
		heading := ""
		if props["heading"].(int) > 9000 {
			heading = "Heading is over 9000"
		} else {
			heading = "Heading is sadly not enough"
		}
		return componentTemplate.Render(map[string]string{
			"heading":    heading,
			"subheading": "Heading 2",
			"children":   children,
			"select":     hxSelect,
		})
	})

	layoutTemplate := dumbo.Parse(`
	<div hx-select="hej">
		<Heading heading={320} />
		<Heading heading={heading2}>
			<p>Paragraph child</p>
		</Heading>
	</div>
	`)

	res := layoutTemplate.Render(map[string]interface{}{
		"heading2": 9001,
	})
	log.Print(res)
}

func TestComponentPropDrilling(t *testing.T) {
	dumbo := template.New()

	dumbo.AddComponent("A", func(children string, props map[string]interface{}) string {
		template := dumbo.Parse(`
		<h1>Component A</h1>
			<B prop={propB}/>
		`)

		return template.Render(props)
	})

	dumbo.AddComponent("B", func(children string, props map[string]interface{}) string {
		template := dumbo.Parse(`
		<h2>Component B: {prop}</h2>
		`)

		return template.Render(props)
	})

	template := dumbo.Parse(`<A propB={9000}/>`)

	res := template.Render(nil)

	log.Print(res)
}
