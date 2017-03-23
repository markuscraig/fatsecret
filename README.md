## Golang FatSecret API Client

Will eventually implement all of the FatSecret API.

You must set the following environment variables in order to access the FatSecret API...

1. `"FATSECRET_CONSUMER_KEY"`
2. `"FATSECRET_SHARED_SECRET"`

## FatSecret Client Example

Use the 'fatsecret_client' cli application as a working example...

```bash
$ go run cmd/fatsecret_client/main.go
```

## FS2JSON Example

Use the 'fs2json' cli tool to invoke FatSecret API calls and dump the JSON response...

```bash
# search for coffee
$ go run cmd/fs2json/main.go -method foods.search -params search_expression=coffee | jq .

# list the food categories
go run cmd/fs2json/main.go -method food_categories.get
```

## References

* https://platform.fatsecret.com/api/
* https://platform.fatsecret.com/api/Default.aspx?screen=rapih

Made with :green_heart: in Campbell, CA
