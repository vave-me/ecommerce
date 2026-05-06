package newsletters

//go:generate buf generate

//go:generate mockery --quiet --dir ./newsletterspb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c newslettersclient -m newslettersclient/models --with-flatten=remove-unused
