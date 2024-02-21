package dumbo

import "io"

func (element *DumboElement) Execute(w io.Writer) error {
	_, err := w.Write([]byte(element.content))

	return err
}

func (template *DumboTemplate) Execute(w io.Writer, data any) error {
	element, err := template.Render(data)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(element.content))

	return err
}
