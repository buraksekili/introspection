package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/jensneuse/graphql-go-tools/pkg/astprinter"
	"github.com/jensneuse/graphql-go-tools/pkg/introspection"
)

const (
	URL = "http://localhost:8080/posts-subgraph/"
)

// Reference: https://github.com/TykTechnologies/graphql-go-tools/blob/master/pkg/introspection/converter_test.go#L18
func main() {
	hc := http.Client{}

	req, err := http.NewRequest(http.MethodPost, URL, strings.NewReader(introspectionQuery))
	if err != nil {
		log.Fatalf("failed to create a new http request, err: %v", err)
	}

	res, err := hc.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		resp, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resp))
		log.Fatal("unexpected status code: ", res.StatusCode)
	}
	defer res.Body.Close()

	type Response struct {
		Data json.RawMessage `json:"data"`
	}

	introspectionResponse := &Response{}
	if err := json.NewDecoder(res.Body).Decode(&introspectionResponse); err != nil {
		log.Fatalf("failed to decode, err: %v", err)
	}

	converter := introspection.JsonConverter{}
	buf := bytes.NewBuffer(introspectionResponse.Data)
	doc, err := converter.GraphQLDocument(buf)
	if err != nil {
		log.Fatal(err)
	}

	outWriter := &bytes.Buffer{}
	err = astprinter.PrintIndent(doc, nil, []byte("  "), outWriter)
	if err != nil {
		log.Fatal(err)
	}

	schemaOutputPretty := outWriter.String()
	fmt.Println(schemaOutputPretty)
}
