package following

//go:generate buf generate

//go:generate mockery --quiet --dir ./storespb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c followingclient -m followingclient/models --with-flatten=remove-unused
