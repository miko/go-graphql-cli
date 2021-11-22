package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/machinebox/graphql"
)

type ResponseStruct interface{}

type stringListFlags []string

func (i *stringListFlags) String() string {
	return "StringList flags"
}
func (i *stringListFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	url       string = "https://countries.trevorblades.com/"
	query     string = `query ($code:String!="US"){countries(filter:{code:{eq:$code}}){ capital name continent {name}}}`
	headers   stringListFlags
	variables stringListFlags
	debug     bool
)

func init() {
	flag.StringVar(&url, "url", url, "Graphql server URL")
	flag.StringVar(&query, "query", query, "Graphql query")
	flag.Var(&headers, "header", "HTTP Header (key: value)")
	flag.Var(&variables, "var", "GraphQL variable (key=value)")
	flag.BoolVar(&debug, "debug", debug, "Debugging")
	flag.Parse()
}
func main() {
	if debug {
		log.Printf("URL: %s\nQuery: %s\n", url, query)
	}
	client := graphql.NewClient(url)

	req := graphql.NewRequest(query)

	for k, v := range variables {
		s := strings.Split(v, "=")
		if len(s) > 1 {
			key := s[0]
			value := strings.Join(s[1:], "=")
			if debug {
				log.Printf("[%d] Setting variable %s to %s\n", k, key, value)
			}
			req.Var(key, value)
		} else {
			log.Printf("WARN: bad variable string %s, needs key=value\n", v)
		}
	}

	for k, v := range headers {
		s := strings.Split(v, ":")
		if len(s) > 1 {
			key := s[0]
			value := strings.Join(s[1:], ":")
			if debug {
				log.Printf("[%d] Setting header %s to %s\n", k, key, value)
			}
			req.Header.Set(key, value)
		} else {
			log.Printf("WARN: bad header string %s, needs key: value\n", v)
		}
	}

	ctx := context.Background()

	var respData ResponseStruct
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	resp, _ := json.MarshalIndent(respData, "", "  ")

	fmt.Printf("%s\n", string(resp))
}
