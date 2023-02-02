package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/fr0stylo/ts-generator/pkg/generator"
	"github.com/fr0stylo/ts-generator/pkg/parser"
	"github.com/fr0stylo/ts-generator/pkg/utils"
)

func main() {
	isInteractive := flag.Bool("i", true, "interactive mode")
	flag.Parse()

	if *isInteractive {
		interactive()
		return
	}

	example := `{
		"name": "test",
		"age": 30,
		"price": 30.56,
		"sizes": [ "L", "XL", "XXL", "XXXL" ],
		"startDate": "2019-06-07",
		"option": {"size": "XL", "color": "red"},
		"options": [{"size": "XL", "color": "red"}]
	}`

	req, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/posts", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	id, err := parser.FromBuffer(bytes.NewBufferString(example), "Main")
	// id, err := parser.FromBuffer(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBufferString("//Code generated by go generate; DO NOT EDIT.\n\n")

	if err := id.Provide(buf); err != nil {
		log.Fatal(err)
	}

	writeToFile("types/types.ts", buf)
	api := generator.Api{BaseURL: fmt.Sprintf("%s://%s/", req.URL.Scheme, req.URL.Host)}
	writeToFile("types/client.ts", api.GenerateApi())
}

func writeToFile(filename string, data io.Reader) error {
	err := os.MkdirAll(path.Dir(filename), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	io.Copy(f, data)
	return nil
}

func interactive() {
	sds := parser.NewSchemaProvider()
	for input := utils.StringPrompt("ts-generate >>"); input != "done" && input != "q" && input != "exit"; input = utils.StringPrompt("ts-generate >>") {
		if len(input) == 0 {
			fmt.Println("Please enter a correct action in order to proceed with api generation")
			fmt.Println("Use `help` for more information")
			continue
		}

		slugs := strings.Split(input, " ")
		switch slugs[0] {
		case "call":
			if len(slugs) != 3 {
				fmt.Print("Please specify type name and url to call")
				continue
			}

			buf := bytes.NewBuffer([]byte{})

			req, err := http.NewRequest("GET", slugs[2], nil)
			if err != nil {
				fmt.Print(err)
				continue
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Print(err)
				continue
			}
			defer resp.Body.Close()

			id, err := parser.FromBuffer(resp.Body, slugs[1])
			if err != nil {
				fmt.Print(err)
				continue
			}

			if err := id.Provide(buf); err != nil {
				fmt.Print(err)
				continue
			}

			sds.Merge(id)

			io.Copy(os.Stdout, buf)
		case "list":
			buf := bytes.NewBufferString("")
			if err := sds.Provide(buf); err != nil {
				fmt.Print(fmt.Errorf("failed to generate, please try again: %s", err))
				continue
			}

			io.Copy(os.Stdout, buf)
		case "save":
			if len(slugs) != 2 {
				fmt.Print("Please specify file name")
			}
			buf := bytes.NewBufferString("//Cde generated by go generate; DO NOT EDIT.\n\n")
			if err := sds.Provide(buf); err != nil {
				fmt.Print(fmt.Errorf("failed to generate, please try again: %s", err))
				continue
			}

			writeToFile(slugs[1], buf)
		default:
			fmt.Print(`
Please refer to corresponding documentation for more information
Application is in testing stage. It is not intended for production use.
Intended for creating API clients and types from provided responses.

Usage:
  call <TYPE_NAME> <URL>    - calls the given URL and parses the response, printing the result and saving the result to inner buffer.
  save <FILE>               - saves the buffer to result file.

  done, q, exit - exits the application.

`)
		}
	}
}
