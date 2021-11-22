package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	url           string = "https://countries.trevorblades.com/"
	query         string = `query ($code:String!="US"){countries(filter:{code:{eq:$code}}){ capital name continent {name}}}`
	queryfile     string
	headers       stringListFlags
	variables     stringListFlags
	filevariables stringListFlags
	debug         bool
)

func init() {

	if v := os.Getenv("GRAPHQL_QUERYFILE"); v != "" {
		log.Printf("Setting queryfile from env")
		queryfile = v
	}
	if v := os.Getenv("GRAPHQL_QUERY"); v != "" {
		log.Printf("Setting query from env")
		query = v
	}
	if v := os.Getenv("GRAPHQL_URL"); v != "" {
		log.Printf("Setting URL from env")
		url = v
	}
	flag.StringVar(&url, "url", url, "Graphql server URL (or GRAPHQL_URL from env)")
	flag.StringVar(&query, "query", query, "Graphql query (or GRAPHQL_QUERY from env)")
	flag.StringVar(&queryfile, "queryfile", queryfile, "File containing graphql query (or GRAPHQL_QUERYFILE from env)")
	flag.Var(&headers, "header", "HTTP Header (key: value)")
	flag.Var(&variables, "var", "GraphQL variable (key=value)")
	flag.Var(&filevariables, "filevar", "GraphQL variable read from file (key=filename)")
	flag.BoolVar(&debug, "debug", debug, "Debugging")
	flag.Parse()

	if queryfile != "" {
		if debug {
			log.Printf("Reading query from file %s\n", queryfile)
		}
		if buf, err := ioutil.ReadFile(queryfile); err != nil {
			log.Fatalf("Cannot open file %s: %s\n", queryfile, err.Error())
		} else {
			query = string(buf)
		}
	}
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

	for k, v := range filevariables {
		s := strings.Split(v, "=")
		if len(s) > 1 {
			key := s[0]
			value := strings.Join(s[1:], "=")
			if debug {
				log.Printf("[%d] Reading variable %s from file %s\n", k, key, value)
			}
			if b, err := ioutil.ReadFile(value); err != nil {
				log.Fatalf("Cannot read variable from file %s: %s", value, err.Error())
			} else {
				var x interface{}
				json.Unmarshal(b, &x)
				req.Var(key, x)
				if debug {
					log.Printf("[%d] Setting variable %s to %s\n", k, key, string(b))
				}
				req.Var(key, x)
			}
		} else {
			log.Printf("WARN: bad variable string %s, needs key=filename\n", v)
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

	if resp, err := json.MarshalIndent(respData, "", "  "); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", string(resp))
	}

}
