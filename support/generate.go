package support

//go:generate buf generate

//go:generate mockery --quiet --dir ./supportpb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c supportclient -m supportclient/models --with-flatten=remove-unused
