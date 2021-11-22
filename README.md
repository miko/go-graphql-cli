# go-graphql-cli
Command -line graphql client written in golang. Based on https://github.com/machinebox/graphql.

# Installation

`go install github.com/miko/go-graphql-cli`

# Usage

```
Usage of go-graphql-cli:
  -debug
    	Debugging
  -header value
    	HTTP Header (key: value)
  -query string
    	Graphql query (or GRAPHQL_QUERY from env) (default "query ($code:String!=\"US\"){countries(filter:{code:{eq:$code}}){ capital name continent {name}}}")
  -url string
    	Graphql server URL (or GRAPHQL_URL from env) (default "https://countries.trevorblades.com/")
  -var value
    	GraphQL variable (key=value)
```
Example:
```
$ export GRAPHQL_URL=https://countries.trevorblades.com/
$ export GRAPHQL_QUERY='query ($code:String!="US"){countries(filter:{code:{eq:$code}}){ capital name continent {name}}}'
$ go-graphql-cli -var code=PL
{
  "countries": [
    {
      "capital": "Warsaw",
      "continent": {
        "name": "Europe"
      },
      "name": "Poland"
    }
  ]
}

$ go-graphql-cli -var user=john -var password=pass -query 'query Login($login:String!, $password:String!) {login(login:$login,password:$password) {first_name last_name}}'
{
  "login": [
    {
      "first_name": "John",
      "last_name": "Doe"
    }
  ]
}

```
