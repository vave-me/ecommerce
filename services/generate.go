package services

//go:generate buf generate

//go:generate mockery --quiet --dir ./servicespb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c servicesclient -m servicesclient/models --with-flatten=remove-unused
