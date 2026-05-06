package offers

//go:generate buf generate

//go:generate mockery --quiet --dir ./offerspb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c offersclient -m offersclient/models --with-flatten=remove-unused
