/*
Copyright 2017 Matt Silverlock
Copyright 2017 Philip Marc Schwartz

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package status

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/oxtoacart/bpool"
	jww "github.com/spf13/jwalterweatherman"
)

var tempPath = "web/endpoints/status/templates/"

var templates map[string]*template.Template
var bufpool *bpool.BufferPool

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	bufpool = bpool.NewBufferPool(64)

	layouts, err := filepath.Glob(tempPath + "layouts/*.tmpl")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Layouts: %+v\n", layouts)

	includes, err := filepath.Glob(tempPath + "includes/*.tmpl")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Includes: %+v\n", includes)

	for _, layout := range layouts {
		files := append(includes, layout)
		templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
	}

	fmt.Printf("Map: %+v\n", templates)
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		efmt := "The template %s does not exist."
		jww.ERROR.Printf(efmt, name)
		return fmt.Errorf(efmt, name)
	}

	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
	return nil
}
