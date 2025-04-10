/*
MIT License

Copyright (c) 2025 Publieke Dienstverlening op de Kaart

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pdok/featureinfo-generator/pkg/types"

	"github.com/iancoleman/strcase"
)

func main() {

	inURL := flag.String("input-url", "", "URL pointing to scheme.json file (https://example.com/scheme.json)")
	inPath := flag.String("input-path", "", "Path pointing to the scheme.json file (../data/scheme.json)")
	destFolder := flag.String("dest-folder", "", "Path pointing to the destination folder (../data/html)")
	v2FileName := flag.String("file-name", "all", "Path pointing to the destination folder (../data/html)")
	flag.Parse()
	validateInputParameters(inPath, inURL, destFolder)
	generateHTMLTemplates(inPath, inURL, destFolder, v2FileName)
}

func validateInputParameters(configURL *string, gpkgPathParam *string, destFolder *string) {
	if *configURL == "" && *gpkgPathParam == "" {
		log.Fatal("input-path or input-url is required. Run with -h for help.")
	} else if *configURL != "" && *gpkgPathParam != "" {
		log.Fatal("either input-path or input-url is required. Run with -h for help.")
	}
	if *destFolder == "" {
		log.Fatal("dest-folder is required. Run with -h for help.")
	}
}

func generateHTMLTemplates(inPath *string, inURL *string, destFolder *string, v2FileName *string) {
	var err error
	var data []byte
	var scheme types.Scheme

	if *inPath != "" {
		data, err = os.ReadFile(*inPath)
	} else {
		data, err = downloadFile(*inURL)
	}

	checkError(err)
	err = json.Unmarshal(data, &scheme)
	checkError(err)

	if scheme.Version == 0 || scheme.Version == 1 {
		log.Fatal("Scheme version of less than 2 is no longer supported")
	}

	if scheme.Version == 2 {
		generateHTMLTemplateForAll(scheme, destFolder, v2FileName)
	}
}

func generateHTMLTemplateForAll(scheme types.Scheme, destFolder *string, v2FileName *string) int {
	err := validateSchemaV2(scheme)
	checkError(err)

	err = os.MkdirAll(*destFolder, os.ModePerm)
	checkError(err)

	log.Print("Generate HTML for all layers: " + *v2FileName + ".html")
	t, err := template.ParseFiles("../internal/resources/layer_template.html")
	checkError(err)

	for _, layer := range scheme.Layers {
		for index := range layer.Properties {
			property := layer.Properties[index]
			layer.Properties[index].ColumnName = property.Name
			if property.Alias != "" {
				layer.Properties[index].Name = property.Alias
			} else if scheme.AutomaticCasing {
				layer.Properties[index].Name = strcase.ToLowerCamel(strings.ToLower(property.Name))
			}
		}
	}

	htmlBuffer := new(bytes.Buffer)
	err = t.Execute(htmlBuffer, scheme)
	checkError(err)
	writeHTMLFile(*v2FileName, destFolder, htmlBuffer)
	return 1
}

func validateSchemaV2(schema types.Scheme) error {
	for i, layer := range schema.Layers {
		if layer.Name == "" {
			return errors.New("The Layer cannot have an empty name for layer number: " + strconv.Itoa(i))
		}
	}
	return nil
}

func writeHTMLFile(name string, destFolder *string, htmlBuffer *bytes.Buffer) {
	var builder strings.Builder
	builder.WriteString(*destFolder)

	if !isLastCharSlash(destFolder) {
		builder.WriteString("/")
	}
	builder.WriteString(name)
	builder.WriteString(".html")
	//nolint:gosec
	errFile := os.WriteFile(builder.String(), htmlBuffer.Bytes(), 0777)
	if errFile != nil {
		log.Fatal("Cannot create temporary file", errFile)
	}
}

func isLastCharSlash(folderPath *string) bool {
	s := *folderPath
	last := s[len(s)-1:]
	return strings.Compare(last, "/") == 0
}

func downloadFile(url string) ([]byte, error) {
	log.Printf("Starting download for: %s", url)
	//nolint:gosec
	resp, err := http.Get(url)
	checkError(err)
	defer resp.Body.Close()
	var file *os.File
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	checkError(err)
	var data []byte
	_, err2 := file.Read(data)
	return data, err2
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
