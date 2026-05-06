 
TEST AGAINST REST API

1. From project root run `make build`
2. Run the command `docker compose  --profile testing up`
3. Go to the `./baskets/ui` and run `npm run pact:publish'
4. Run test for gateway `go test ./internal/rest/`