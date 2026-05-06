package posts

//go:generate buf generate

//go:generate mockery --quiet --dir ./postspb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c postsclient -m postsclient/models --with-flatten=remove-unused
