package message

//go:generate buf generate

//go:generate mockery --quiet --dir ./messagespb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c messagesclient -m messagesclient/models --with-flatten=remove-unused
