package template_test

import (
	"testing"

	"github.com/JHelar/dumbo/template"
)

func TestParse(t *testing.T) {
	content := []byte("    test")

	template.Parse(content)
}
