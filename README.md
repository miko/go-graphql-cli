# go-graphql-cli
Command -line graphql client written in golang. Based on https://github.com/machinebox/graphql.

# Usage

```
Usage of ./go-graphql-cli:
  -debug
    	Debugging
  -header value
    	HTTP Header (key: value)
  -query string
    	Graphql query (default "query ($code:String!=\"US\"){countries(filter:{code:{eq:$code}}){ capital name continent {name}}}")
  -url string
    	Graphql server URL (default "https://countries.trevorblades.com/")
  -var value
    	GraphQL variable (key=value)
```
Example:
```
./go-graphql-cli -var code=PL
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
```
