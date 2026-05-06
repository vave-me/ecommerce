package users

//go:generate buf generate

//go:generate mockery --quiet --dir ./userspb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c usersclient -m usersclient/models --with-flatten=remove-unused
