package mailer

//go:generate buf generate

//go:generate mockery --quiet --dir ./mailerpb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c mailerclient -m mailerclient/models --with-flatten=remove-unused
