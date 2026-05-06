package activity

//go:generate buf generate

//go:generate mockery --quiet --dir ./activitypb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c activityclient -m activityclient/models --with-flatten=remove-unused
