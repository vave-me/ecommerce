package assistants

//go:generate buf generate

//go:generate mockery --quiet --dir ./assistantspb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c assistantsclient -m assistantsclient/models --with-flatten=remove-unused
