package tickets

//go:generate buf generate

//go:generate mockery --quiet --dir ./ticketspb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c ticketsclient -m ticketsclient/models --with-flatten=remove-unused
