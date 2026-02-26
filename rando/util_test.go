package rando

import (
	"os"
	"text/template"
)

func ExampleTemplateFuncs() {
	t, err := template.
		New("example").
		Funcs(TemplateFuncs()).
		Parse(`{{randoAscii 5 true ""}}
		{{randoAnyString 10 "X"}}
		{{randoHostname "common"}}
		{{randoUniqueDNSName "proper" "watttttt.net"}}
		{{randoIsProfane "shit"}}
		`)
	if err != nil {
		panic(err)
	}
	if err = t.Execute(os.Stderr, nil); err != nil {
		panic(err)
	}
	// Output:
}
